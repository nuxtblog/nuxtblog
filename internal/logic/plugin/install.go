package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ── Install ───────────────────────────────────────────────────────────────

// Install fetches the latest GitHub Release zip for the given repo and installs the plugin.
func (s *sPlugin) Install(ctx context.Context, repoUrl, expectedVersion string) (*v1.PluginItem, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	owner, repo, err := parseRepo(repoUrl)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, err.Error())
	}

	// Fetch latest release metadata from GitHub API.
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	releaseData, err := pluginHTTPGet(ctx, apiURL)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, gerror.NewCode(gcode.CodeNotFound,
				fmt.Sprintf("repository %s/%s has no published releases — ask the author to publish a GitHub Release first", owner, repo))
		}
		return nil, gerror.NewCode(gcode.CodeInternalError,
			fmt.Sprintf("cannot reach GitHub for %s/%s: %v (check network or GitHub status)", owner, repo, err))
	}

	var release struct {
		TagName    string `json:"tag_name"`
		ZipballURL string `json:"zipball_url"`
		Assets     []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err = json.Unmarshal(releaseData, &release); err != nil || release.ZipballURL == "" {
		return nil, gerror.NewCode(gcode.CodeInternalError,
			fmt.Sprintf("GitHub returned malformed release data for %s/%s — the release may be a draft or pre-release", owner, repo))
	}

	if expectedVersion != "" {
		got := strings.TrimPrefix(release.TagName, "v")
		want := strings.TrimPrefix(expectedVersion, "v")
		if got != want {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter,
				fmt.Sprintf(
					"registry version mismatch: registry lists %s/%s@%s but latest GitHub release is %s. "+
						"The registry may be stale, or the author hasn't published the new release yet.",
					owner, repo, expectedVersion, release.TagName))
		}
	}

	// Prefer plugin.zip asset (contains built files) over the source zipball.
	downloadURL := release.ZipballURL
	for _, a := range release.Assets {
		if a.Name == "plugin.zip" {
			downloadURL = a.BrowserDownloadURL
			break
		}
	}

	zipData, err := pluginHTTPGet(ctx, downloadURL)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError,
			fmt.Sprintf("cannot download release zip (%s): %v", release.TagName, err))
	}

	return s.installFromZipBytes(ctx, zipData, "https://github.com/"+owner+"/"+repo)
}

// ── Update ────────────────────────────────────────────────────────────────

// Update pulls the latest version of a plugin from its GitHub repo URL.
func (s *sPlugin) Update(ctx context.Context, id string) (*v1.PluginItem, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	type updateRow struct {
		RepoUrl  string `orm:"repo_url"`
		Manifest string `orm:"manifest"`
	}
	var r updateRow
	if err := g.DB().Ctx(ctx).Model("plugins").Where("id", id).
		Fields("repo_url, manifest").Scan(&r); err != nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "plugin not found")
	}

	repoUrl := r.RepoUrl

	// Fallback: extract homepage from stored manifest
	if repoUrl == "" && r.Manifest != "" {
		var mf struct {
			Homepage string `json:"homepage"`
		}
		if json.Unmarshal([]byte(r.Manifest), &mf) == nil && mf.Homepage != "" {
			repoUrl = mf.Homepage
		}
	}

	if repoUrl == "" {
		return nil, gerror.NewCode(gcode.CodeNotFound,
			"plugin has no repo_url or homepage — cannot update automatically")
	}
	return s.Install(ctx, repoUrl, "")
}

// ── InstallZip ────────────────────────────────────────────────────────────

// InstallZip installs a plugin from raw archive bytes (local upload).
func (s *sPlugin) InstallZip(ctx context.Context, zipData []byte) (*v1.PluginItem, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	return s.installFromZipBytes(ctx, zipData, "")
}

// installFromZipBytes is the shared core: parse archive → validate → upsert DB → load engine.
func (s *sPlugin) installFromZipBytes(ctx context.Context, zipData []byte, repoUrl string) (*v1.PluginItem, error) {
	ar, err := parseArchive(zipData)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, err.Error())
	}
	manifest := ar.Manifest

	// Use homepage from package.json as repo_url fallback
	if repoUrl == "" && ar.Homepage != "" {
		repoUrl = ar.Homepage
	}

	// ── JS (Goja) plugin path ───────────────────────────────────────────────
	if ar.IsJS {
		return s.installJSPlugin(ctx, ar, repoUrl)
	}

	// ── Builtin (compiled Go) plugin path ───────────────────────────────────
	if ar.IsBuiltin {
		return s.installBuiltinPlugin(ctx, ar, repoUrl)
	}

	// ── Legacy script plugin path ───────────────────────────────────────────
	script := ar.Script

	// Phase 2.8: save frontend assets to data/plugins/{id}/
	if len(ar.Assets) > 0 {
		if saveErr := savePluginAssets(manifest.Name, ar.Assets); saveErr != nil {
			g.Log().Warningf(ctx, "[plugin] save assets for %s: %v", manifest.Name, saveErr)
		}
	}

	if loadErr := eng.Load(manifest.Name, script, *manifest); loadErr != nil {
		eng.Unload(manifest.Name)
		return nil, gerror.NewCode(gcode.CodeInternalError,
			fmt.Sprintf("plugin script error: %v", loadErr))
	}

	// Phase 4.1: run pending DB migrations
	if len(manifest.Migrations) > 0 {
		if n, migErr := eng.RunMigrations(ctx, manifest.Name, manifest.Migrations); migErr != nil {
			g.Log().Warningf(ctx, "[plugin] migration error for %s: %v", manifest.Name, migErr)
		} else if n > 0 {
			g.Log().Infof(ctx, "[plugin] applied %d migration(s) for %s", n, manifest.Name)
		}
	}

	now := gtime.Now()
	schemaJSON := manifestSettingsJSON(manifest)
	capsJSON := manifestCapabilitiesJSON(manifest)
	manifestJSON := manifestFullJSON(manifest)

	cnt, _ := g.DB().Ctx(ctx).Model("plugins").Where("id", manifest.Name).Count()
	if cnt > 0 {
		_, err = g.DB().Ctx(ctx).Model("plugins").Where("id", manifest.Name).Data(g.Map{
			"title":           manifest.Title,
			"description":     manifest.Description,
			"version":         manifest.Version,
			"author":          manifest.Author,
			"icon":            orDefault(manifest.Icon, "i-tabler-plug"),
			"repo_url":        repoUrl,
			"script":          script,
			"styles":          manifest.CSS,
			"priority":        manifest.Priority,
			"settings_schema": schemaJSON,
			"capabilities":    capsJSON,
			"manifest":        manifestJSON,
			"enabled":         1,
			"updated_at":      now,
		}).Update()
	} else {
		_, err = g.DB().Ctx(ctx).Model("plugins").Data(g.Map{
			"id":              manifest.Name,
			"title":           manifest.Title,
			"description":     manifest.Description,
			"version":         manifest.Version,
			"author":          manifest.Author,
			"icon":            orDefault(manifest.Icon, "i-tabler-plug"),
			"repo_url":        repoUrl,
			"script":          script,
			"styles":          manifest.CSS,
			"priority":        manifest.Priority,
			"settings":        manifestDefaultSettings(manifest),
			"settings_schema": schemaJSON,
			"capabilities":    capsJSON,
			"manifest":        manifestJSON,
			"enabled":         1,
			"installed_at":    now,
			"updated_at":      now,
		}).Insert()
	}
	if err != nil {
		eng.Unload(manifest.Name)
		return nil, err
	}

	return &v1.PluginItem{
		Id:          manifest.Name,
		Title:       manifest.Title,
		Description: manifest.Description,
		Version:     manifest.Version,
		Author:      manifest.Author,
		Icon:        orDefault(manifest.Icon, "i-tabler-plug"),
		RepoUrl:     repoUrl,
		Enabled:     true,
		InstalledAt: now.String(),
		Type:        manifest.Type,
	}, nil
}

// ── BatchUpdate ──────────────────────────────────────────────────────────

func (s *sPlugin) BatchUpdate(ctx context.Context, ids []string) (*v1.PluginBatchUpdateRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	res := &v1.PluginBatchUpdateRes{}
	for _, id := range ids {
		if _, err := s.Update(ctx, id); err != nil {
			res.Failed = append(res.Failed, id)
		} else {
			res.Succeeded = append(res.Succeeded, id)
		}
	}
	return res, nil
}

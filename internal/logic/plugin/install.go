package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
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

// installBuiltinPlugin handles "type: builtin" plugins installed from a zip archive.
// It extracts all plugin files (Go source, yaml, web/dist/) to builtin/{dirName}/
// preserving the original directory structure, then regenerates builtin/plugins.go.
// The plugin activates after the next server rebuild + restart.
func (s *sPlugin) installBuiltinPlugin(ctx context.Context, ar *parseArchiveResult, repoUrl string) (*v1.PluginItem, error) {
	manifest := ar.Manifest
	pluginID := manifest.Name
	prefix := detectArchivePrefix(ar.AllFiles)

	// Detect Go package name from plugin.go (needed for the blank import)
	goPkgName := ""
	for archivePath, content := range ar.AllFiles {
		relPath := strings.TrimPrefix(archivePath, prefix)
		if path.Base(relPath) == "plugin.go" && !strings.Contains(relPath, "/") {
			goPkgName = detectGoPkgName(string(content))
			break
		}
	}
	if goPkgName == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter,
			"builtin plugin archive must contain plugin.go with a Go package declaration")
	}

	// Use plugin ID as directory name: nuxtblog-plugin-view-counter -> nuxtblog_plugin_view_counter
	dirName := strings.ReplaceAll(pluginID, "-", "_")

	// Resolve builtin/ directory
	builtinDir := builtinSourceDir()
	if builtinDir == "" {
		return nil, gerror.NewCode(gcode.CodeInternalError,
			"cannot locate builtin/ source directory — is this a development environment?")
	}

	// Extract files to builtin/{dirName}/, preserving original directory structure
	targetDir := filepath.Join(builtinDir, dirName)
	// Clean previous installation
	_ = os.RemoveAll(targetDir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError,
			fmt.Sprintf("create builtin dir: %v", err))
	}

	extracted := 0
	for archivePath, content := range ar.AllFiles {
		relPath := archivePath
		if prefix != "" {
			relPath = strings.TrimPrefix(archivePath, prefix)
		}
		if relPath == "" {
			continue
		}

		// Skip non-runtime files
		if isSkippedBuiltinPath(relPath) {
			continue
		}

		destPath := filepath.Join(targetDir, filepath.FromSlash(relPath))
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return nil, gerror.NewCode(gcode.CodeInternalError,
				fmt.Sprintf("create dir %s: %v", destDir, err))
		}
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			return nil, gerror.NewCode(gcode.CodeInternalError,
				fmt.Sprintf("write %s: %v", relPath, err))
		}
		extracted++
	}

	g.Log().Infof(ctx, "[plugin] extracted builtin plugin %s to %s (%d files)", pluginID, targetDir, extracted)

	// Regenerate builtin/plugins.go with updated imports
	if err := regenerateBuiltinImports(builtinDir); err != nil {
		g.Log().Warningf(ctx, "[plugin] regenerate builtin imports: %v", err)
	}

	// Run DB migrations
	if len(manifest.Migrations) > 0 {
		if n, migErr := eng.RunMigrations(ctx, pluginID, manifest.Migrations); migErr != nil {
			g.Log().Warningf(ctx, "[plugin] migration error for %s: %v", pluginID, migErr)
		} else if n > 0 {
			g.Log().Infof(ctx, "[plugin] applied %d migration(s) for %s", n, pluginID)
		}
	}

	// Upsert DB record
	now := gtime.Now()
	schemaJSON := manifestSettingsJSON(manifest)
	capsJSON := manifestCapabilitiesJSON(manifest)
	manifestJSON := manifestFullJSON(manifest)

	cnt, _ := g.DB().Ctx(ctx).Model("plugins").Where("id", pluginID).Count()
	var err error
	if cnt > 0 {
		_, err = g.DB().Ctx(ctx).Model("plugins").Where("id", pluginID).Data(g.Map{
			"title":           manifest.Title,
			"description":     manifest.Description,
			"version":         manifest.Version,
			"author":          manifest.Author,
			"icon":            orDefault(manifest.Icon, "i-tabler-plug"),
			"repo_url":        repoUrl,
			"priority":        manifest.Priority,
			"settings_schema": schemaJSON,
			"capabilities":    capsJSON,
			"manifest":        manifestJSON,
			"source":          "builtin",
			"enabled":         1,
			"updated_at":      now,
		}).Update()
	} else {
		_, err = g.DB().Ctx(ctx).Model("plugins").Where("id", pluginID).Data(g.Map{
			"id":              pluginID,
			"title":           manifest.Title,
			"description":     manifest.Description,
			"version":         manifest.Version,
			"author":          manifest.Author,
			"icon":            orDefault(manifest.Icon, "i-tabler-plug"),
			"repo_url":        repoUrl,
			"priority":        manifest.Priority,
			"settings":        manifestDefaultSettings(manifest),
			"settings_schema": schemaJSON,
			"capabilities":    capsJSON,
			"manifest":        manifestJSON,
			"source":          "builtin",
			"enabled":         1,
			"installed_at":    now,
			"updated_at":      now,
		}).Insert()
	}
	if err != nil {
		return nil, err
	}

	return &v1.PluginItem{
		Id:          pluginID,
		Title:       manifest.Title,
		Description: manifest.Description,
		Version:     manifest.Version,
		Author:      manifest.Author,
		Icon:        orDefault(manifest.Icon, "i-tabler-plug"),
		RepoUrl:     repoUrl,
		Enabled:     true,
		InstalledAt: now.String(),
		Source:      "builtin",
		Type:        "builtin",
		NeedRestart: true,
	}, nil
}

// isSkippedBuiltinPath returns true for files not needed in builtin/.
func isSkippedBuiltinPath(relPath string) bool {
	base := path.Base(relPath)
	switch base {
	case "go.mod", "go.sum", ".gitignore", ".gitattributes":
		return true
	}
	parts := strings.Split(relPath, "/")
	for _, p := range parts {
		switch p {
		case ".git", ".github", "node_modules", ".idea", ".vscode":
			return true
		}
	}
	return false
}

// builtinSourceDir locates the builtin/ source directory.
// It walks up from the current working directory looking for the builtin/plugins.go file.
func builtinSourceDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	// In dev mode the server typically runs from the nuxtblog/ (server) directory.
	// builtin/ is a direct child of that directory.
	candidate := filepath.Join(cwd, "builtin")
	if _, err := os.Stat(filepath.Join(candidate, "plugins.go")); err == nil {
		return candidate
	}
	// Walk up a few levels in case cwd is deeper
	dir := cwd
	for i := 0; i < 4; i++ {
		dir = filepath.Dir(dir)
		candidate = filepath.Join(dir, "builtin")
		if _, err := os.Stat(filepath.Join(candidate, "plugins.go")); err == nil {
			return candidate
		}
	}
	return ""
}

// detectGoPkgName extracts the Go package name from source code.
func detectGoPkgName(src string) string {
	for _, line := range strings.Split(src, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "package ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	}
	return ""
}

// regenerateBuiltinImports scans builtin/ for plugin subdirectories (each containing
// plugin.go) and regenerates builtin/plugins.go with the corresponding blank imports.
// The directory name is used as the import path segment, while the actual Go package
// name inside may differ (Go allows this).
func regenerateBuiltinImports(builtinDir string) error {
	const modulePath = "github.com/nuxtblog/nuxtblog"

	entries, err := os.ReadDir(builtinDir)
	if err != nil {
		return err
	}

	var dirs []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		// Check if directory contains a plugin.go file
		if _, err := os.Stat(filepath.Join(builtinDir, e.Name(), "plugin.go")); err == nil {
			dirs = append(dirs, e.Name())
		}
	}

	var buf strings.Builder
	buf.WriteString("// Code generated by plugin installer; DO NOT EDIT.\n\n")
	buf.WriteString("package builtin\n")
	if len(dirs) > 0 {
		buf.WriteString("\nimport (\n")
		for _, dir := range dirs {
			fmt.Fprintf(&buf, "\t_ \"%s/builtin/%s\"\n", modulePath, dir)
		}
		buf.WriteString(")\n")
	}

	return os.WriteFile(filepath.Join(builtinDir, "plugins.go"), []byte(buf.String()), 0644)
}

// installingSetKey is a context key whose value is a map[string]bool tracking
// plugin IDs currently being installed on the call stack. Used to detect
// circular dependencies during recursive auto-install.
type installingSetKey struct{}

// ensureDependenciesInstalled walks the `depends:` list of the plugin being
// installed and, for each non-optional dependency that isn't already present,
// looks up its repository URL in the marketplace cache and recursively calls
// Install. Circular dependencies are detected via a context-scoped visited set.
func (s *sPlugin) ensureDependenciesInstalled(ctx context.Context, parentID string, deps []pluginDependency) error {
	if len(deps) == 0 {
		return nil
	}
	mgr := eng.GetManager()
	if mgr == nil {
		return gerror.NewCode(gcode.CodeInternalError, "plugin manager not initialized")
	}

	// Inherit or create the installing set, then mark the parent as in-progress.
	set, _ := ctx.Value(installingSetKey{}).(map[string]bool)
	if set == nil {
		set = map[string]bool{}
	}
	set[parentID] = true
	ctx = context.WithValue(ctx, installingSetKey{}, set)

	for _, dep := range deps {
		if dep.ID == "" || dep.Optional {
			continue
		}
		if mgr.HasPlugin(dep.ID) {
			// Already installed — version check will be enforced later by
			// pluginsys.Manager.InstallJSPlugin against the resolved constraint.
			continue
		}
		if set[dep.ID] {
			return gerror.NewCode(gcode.CodeInternalError,
				fmt.Sprintf("circular dependency involving %s", dep.ID))
		}
		repoURL, ok := lookupMarketplaceRepoByID(ctx, dep.ID)
		if !ok {
			return gerror.NewCode(gcode.CodeInvalidParameter, fmt.Sprintf(
				"required dependency %q is not installed and not found in marketplace; "+
					"please sync the marketplace or install it manually first", dep.ID))
		}
		g.Log().Infof(ctx, "[plugin] auto-installing dependency %s (from %s) for %s",
			dep.ID, repoURL, parentID)
		if _, err := s.Install(ctx, repoURL, ""); err != nil {
			return gerror.NewCode(gcode.CodeInternalError,
				fmt.Sprintf("failed to auto-install dependency %s: %v", dep.ID, err))
		}
	}
	return nil
}

// lookupMarketplaceRepoByID searches the cached marketplace registry for an
// entry whose Name matches id and returns its repository URL. If the cache is
// empty, a sync is attempted inline. Returns ("", false) if not found.
func lookupMarketplaceRepoByID(ctx context.Context, id string) (string, bool) {
	cache := getOption(ctx, registryCacheKey)
	if cache == "" {
		if _, err := (&sPlugin{}).SyncMarketplace(ctx); err != nil {
			g.Log().Warningf(ctx, "[plugin] marketplace sync for dep lookup failed: %v", err)
			return "", false
		}
		cache = getOption(ctx, registryCacheKey)
	}
	if cache == "" {
		return "", false
	}
	var items []v1.MarketplaceItem
	if err := json.Unmarshal([]byte(cache), &items); err != nil {
		return "", false
	}
	for _, it := range items {
		if it.Name == id && it.Repo != "" {
			return it.Repo, true
		}
	}
	return "", false
}

// installJSPlugin extracts a JS (Goja) plugin to data/plugins/{id}/ and loads it.
func (s *sPlugin) installJSPlugin(ctx context.Context, ar *parseArchiveResult, repoUrl string) (*v1.PluginItem, error) {
	manifest := ar.Manifest
	pluginID := manifest.Name

	// Auto-install missing required dependencies from marketplace before extracting.
	// Doing this before writing any files lets us bail out cleanly on failure.
	if err := s.ensureDependenciesInstalled(ctx, pluginID, ar.Depends); err != nil {
		return nil, err
	}

	// Determine the directory prefix to strip from archive paths.
	// GitHub zipball wraps files in "{owner}-{repo}-{sha}/" directory.
	prefix := detectArchivePrefix(ar.AllFiles)

	// Extract all files to data/plugins/{id}/
	pluginDir := filepath.Join(PluginAssetsDir(), sanitizePluginPath(pluginID))
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError,
			fmt.Sprintf("create plugin dir: %v", err))
	}

	for archivePath, content := range ar.AllFiles {
		// Strip archive prefix (e.g., "owner-repo-sha/")
		relPath := archivePath
		if prefix != "" {
			relPath = strings.TrimPrefix(archivePath, prefix)
		}
		if relPath == "" {
			continue
		}

		// Skip .git directory, CI configs, and other non-runtime files
		if isSkippedArchivePath(relPath) {
			continue
		}

		destPath := filepath.Join(pluginDir, filepath.FromSlash(relPath))
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return nil, gerror.NewCode(gcode.CodeInternalError,
				fmt.Sprintf("create dir %s: %v", destDir, err))
		}
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			return nil, gerror.NewCode(gcode.CodeInternalError,
				fmt.Sprintf("write %s: %v", relPath, err))
		}
	}

	g.Log().Infof(ctx, "[plugin] extracted plugin %s to %s (%d files)", pluginID, pluginDir, len(ar.AllFiles))

	// Load via Goja
	mgr := eng.GetManager()
	if mgr == nil {
		return nil, gerror.NewCode(gcode.CodeInternalError, "plugin manager not initialized")
	}

	jsFile := ar.JSEntry
	if jsFile == "" {
		jsFile = "plugin.js"
	}
	if err := mgr.InstallJSPlugin(ctx, pluginDir, jsFile); err != nil {
		// Cleanup extracted files on failure
		_ = os.RemoveAll(pluginDir)
		return nil, gerror.NewCode(gcode.CodeInternalError,
			fmt.Sprintf("goja load failed: %v", err))
	}

	// Update repo_url in DB (ensureDBRecordExt doesn't set it)
	if repoUrl != "" {
		_, _ = g.DB().Ctx(ctx).Model("plugins").Where("id", pluginID).
			Data(g.Map{"repo_url": repoUrl}).Update()
	}

	now := gtime.Now()
	return &v1.PluginItem{
		Id:          pluginID,
		Title:       manifest.Title,
		Description: manifest.Description,
		Version:     manifest.Version,
		Author:      manifest.Author,
		Icon:        orDefault(manifest.Icon, "i-tabler-plug"),
		RepoUrl:     repoUrl,
		Enabled:     true,
		InstalledAt: now.String(),
		Source:      "external",
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

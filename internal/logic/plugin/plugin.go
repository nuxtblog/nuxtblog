package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/mholt/archives"
	"gopkg.in/yaml.v3"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sPlugin struct{}

func New() service.IPlugin  { return &sPlugin{} }
func init()                  { service.RegisterPlugin(New()) }

// ── List ──────────────────────────────────────────────────────────────────

func (s *sPlugin) List(ctx context.Context) (*v1.PluginListRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	type row struct {
		Id           string `orm:"id"`
		Title        string `orm:"title"`
		Description  string `orm:"description"`
		Version      string `orm:"version"`
		Author       string `orm:"author"`
		Icon         string `orm:"icon"`
		RepoUrl      string `orm:"repo_url"`
		Enabled      int    `orm:"enabled"`
		InstalledAt  string `orm:"installed_at"`
		Capabilities string `orm:"capabilities"`
		Source       string `orm:"source"`
		Manifest     string `orm:"manifest"`
	}
	var rows []row
	if err := g.DB().Ctx(ctx).Model("plugins").
		Fields("id,title,description,version,author,icon,repo_url,enabled,installed_at,COALESCE(capabilities,'{}') as capabilities,COALESCE(source,'external') as source,COALESCE(manifest,'{}') as manifest").
		OrderDesc("installed_at").
		Scan(&rows); err != nil {
		return nil, err
	}

	items := make([]v1.PluginItem, 0, len(rows))
	for _, r := range rows {
		source := r.Source
		if source == "" {
			source = "external"
		}
		// Extract type from stored manifest JSON
		pluginType := ""
		if r.Manifest != "" {
			var mf struct {
				Type string `json:"type"`
			}
			if json.Unmarshal([]byte(r.Manifest), &mf) == nil {
				pluginType = mf.Type
			}
		}
		// Derive type from source if not in manifest
		if pluginType == "" {
			if source == "builtin" {
				pluginType = "builtin"
			} else {
				pluginType = "js" // default for external
			}
		}
		items = append(items, v1.PluginItem{
			Id:           r.Id,
			Title:        r.Title,
			Description:  r.Description,
			Version:      r.Version,
			Author:       r.Author,
			Icon:         r.Icon,
			RepoUrl:      r.RepoUrl,
			Enabled:      r.Enabled == 1,
			InstalledAt:  r.InstalledAt,
			Capabilities: r.Capabilities,
			Source:       source,
			Type:         pluginType,
		})
	}

	// Merge Layer 0 (YAML declarative plugins) — always enabled, no DB record
	for _, yp := range eng.GetAllYAMLPlugins() {
		items = append(items, v1.PluginItem{
			Id:          yp.ID,
			Title:       yp.Title,
			Description: yp.Description,
			Version:     yp.Version,
			Author:      yp.Author,
			Icon:        yp.Icon,
			Enabled:     true,
			Type:        "yaml",
		})
	}

	return &v1.PluginListRes{Items: items}, nil
}

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

// installJSPlugin extracts a JS (Goja) plugin to data/plugins/{id}/ and loads it.
func (s *sPlugin) installJSPlugin(ctx context.Context, ar *parseArchiveResult, repoUrl string) (*v1.PluginItem, error) {
	manifest := ar.Manifest
	pluginID := manifest.Name

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

// isSkippedArchivePath returns true for files that should NOT be extracted to data/plugins/.
func isSkippedArchivePath(relPath string) bool {
	parts := strings.Split(relPath, "/")
	for _, p := range parts {
		switch p {
		case ".git", ".github", "node_modules", ".idea", ".vscode":
			return true
		}
	}
	base := path.Base(relPath)
	switch base {
	case ".gitignore", ".gitattributes", "go.mod", "go.sum":
		return true
	}
	return false
}

// detectArchivePrefix finds the common directory prefix in archive paths.
// GitHub zipball wraps all files under "{owner}-{repo}-{sha}/".
func detectArchivePrefix(files map[string][]byte) string {
	var prefix string
	for name := range files {
		idx := strings.IndexByte(name, '/')
		if idx < 0 {
			return "" // file at root — no prefix
		}
		candidate := name[:idx+1]
		if prefix == "" {
			prefix = candidate
		} else if candidate != prefix {
			return "" // inconsistent — no common prefix
		}
	}
	return prefix
}

// parseArchive reads the plugin manifest and entry script from any archive format
// supported by mholt/archives (zip, tar.gz, tar.bz2, tar.xz, 7z, rar, …).
//
// Manifest resolution order:
//  1. plugin.yaml — new declarative format (preferred)
//  2. package.json with a "plugin" field — legacy npm-style manifest
//
// parseArchiveResult holds all data extracted from a plugin archive.
type parseArchiveResult struct {
	Manifest  *eng.Manifest
	Script    string
	Assets    map[string][]byte
	AllFiles  map[string][]byte // all files in archive (for JS/builtin plugins that need full extraction)
	Homepage  string            // from package.json "homepage" field, used as repo_url fallback
	IsJS      bool              // true when plugin type is "js" or "full"
	IsBuiltin bool              // true when plugin type is "builtin" (compiled Go)
	JSEntry   string            // JS entry file name from plugin.yaml (default "plugin.js")
}

// pluginYAMLManifest is the YAML structure for plugin.yaml files.
type pluginYAMLManifest struct {
	ID          string `yaml:"id"`
	Title       string `yaml:"title"`
	Version     string `yaml:"version"`
	Icon        string `yaml:"icon"`
	Author      string `yaml:"author"`
	Description string `yaml:"description"`
	License     string `yaml:"license"`
	TrustLevel  string `yaml:"trust_level"`
	SDKVersion  string `yaml:"sdk_version"`
	Homepage    string `yaml:"homepage"`

	Type    string `yaml:"type"`     // builtin, js, yaml, ui, full
	Runtime string `yaml:"runtime"` // compiled | interpreted
	Bundled bool   `yaml:"bundled"` // included in official prebuilt binary
	JSEntry string `yaml:"js_entry"` // JS entry file (default "plugin.js")
	Layer    any    `yaml:"layer"`    // int or []int
	Entry    string `yaml:"entry"`
	AdminJS  string `yaml:"admin_js"`
	PublicJS string `yaml:"public_js"`
	CSS      string `yaml:"css"`

	Settings []struct {
		Key         string   `yaml:"key"`
		Label       string   `yaml:"label"`
		Type        string   `yaml:"type"`
		Required    bool     `yaml:"required"`
		Default     any      `yaml:"default"`
		Placeholder string   `yaml:"placeholder"`
		Description string   `yaml:"description"`
		Options     []string `yaml:"options"`
		Group       string   `yaml:"group"`
		Shared      bool     `yaml:"shared"`
	} `yaml:"settings"`

	Pages []struct {
		Path      string `yaml:"path"`
		Slot      string `yaml:"slot"`
		Component string `yaml:"component"`
		Title     string `yaml:"title"`
		Nav       *struct {
			Group string `yaml:"group"`
			Icon  string `yaml:"icon"`
			Order int    `yaml:"order"`
		} `yaml:"nav"`
	} `yaml:"pages"`

	Contributes *struct {
		Commands []struct {
			ID      string `yaml:"id"`
			Title   string `yaml:"title"`
			TitleEn string `yaml:"title_en"`
			Icon    string `yaml:"icon"`
		} `yaml:"commands"`
		Menus map[string][]struct {
			Command string `yaml:"command"`
		} `yaml:"menus"`
	} `yaml:"contributes"`

	Migrations []struct {
		Version int    `yaml:"version"`
		Up      string `yaml:"up"`
		Down    string `yaml:"down"`
	} `yaml:"migrations"`

	Routes []struct {
		Method      string `yaml:"method"`
		Path        string `yaml:"path"`
		Auth        string `yaml:"auth"`
		Fn          string `yaml:"fn"`
		Description string `yaml:"description"`
	} `yaml:"routes"`

	Depends []struct {
		ID       string `yaml:"id"`
		Version  string `yaml:"version"`
		Optional bool   `yaml:"optional"`
	} `yaml:"depends"`
}

func parseArchive(data []byte) (*parseArchiveResult, error) {
	ctx := context.Background()

	format, reader, err := archives.Identify(ctx, "", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("unsupported archive format: %w", err)
	}

	ex, ok := format.(archives.Extractor)
	if !ok {
		return nil, fmt.Errorf("archive format %q cannot be extracted", format.Extension())
	}

	// Collect every file in the archive into a map (clean path → content).
	collected := make(map[string][]byte)
	if err = ex.Extract(ctx, reader, func(ctx context.Context, f archives.FileInfo) error {
		if f.IsDir() {
			return nil
		}
		rc, e := f.Open()
		if e != nil {
			return fmt.Errorf("open %s: %w", f.NameInArchive, e)
		}
		defer rc.Close()
		b, e := io.ReadAll(rc)
		if e != nil {
			return fmt.Errorf("read %s: %w", f.NameInArchive, e)
		}
		collected[path.Clean(f.NameInArchive)] = b
		return nil
	}); err != nil {
		return nil, fmt.Errorf("extract archive: %w", err)
	}

	// ── Strategy 1: Try plugin.yaml first ────────────────────────────────────

	var m eng.Manifest
	var homepage string
	var yamlJSEntry string
	var yamlType string
	resolved := false

	for name, content := range collected {
		base := path.Base(name)
		if base != "plugin.yaml" && base != "plugin.yml" {
			continue
		}
		// Parse raw YAML to extract js_entry and type before converting to Manifest
		var raw pluginYAMLManifest
		if yaml.Unmarshal(content, &raw) == nil {
			yamlJSEntry = raw.JSEntry
			yamlType = raw.Type
		}
		if mf, hp, err := parsePluginYAML(content); err == nil {
			m = *mf
			homepage = hp
			resolved = true
			break
		}
	}

	// ── Strategy 2: Fallback to package.json "plugin" field ──────────────────

	if !resolved {
		for name, content := range collected {
			if path.Base(name) != "package.json" {
				continue
			}
			var pkg struct {
				Name        string `json:"name"`
				Version     string `json:"version"`
				Description string `json:"description"`
				Author      any    `json:"author"`
				Homepage    string `json:"homepage"`
				Plugin      *struct {
					Title            string                 `json:"title"`
					Icon             string                 `json:"icon"`
					Entry            string                 `json:"entry"`
					CSS              string                 `json:"css"`
					Priority         int                    `json:"priority"`
					Settings         []eng.SettingField     `json:"settings"`
					Capabilities     eng.PluginCapabilities `json:"capabilities"`
					Webhooks         []eng.WebhookDef       `json:"webhooks"`
					Pipelines        []eng.PipelineDef      `json:"pipelines"`
					MinHostVersion   string                 `json:"minHostVersion"`
					TrustLevel       eng.TrustLevel         `json:"trust_level"`
					ActivationEvents []string               `json:"activationEvents"`
					AdminJS          string                 `json:"admin_js"`
					PublicJS         string                 `json:"public_js"`
					Routes           []eng.RouteDef         `json:"routes"`
					Contributes      *eng.Contributes       `json:"contributes"`
					Migrations       []eng.MigrationDef     `json:"migrations"`
					Pages            []eng.PageDef          `json:"pages"`
					Service          *eng.ServiceDef        `json:"service"`
				} `json:"plugin"`
			}
			if err = json.Unmarshal(content, &pkg); err != nil || pkg.Plugin == nil {
				continue
			}
			m.Name = pkg.Name
			m.Version = pkg.Version
			m.Description = pkg.Description
			switch v := pkg.Author.(type) {
			case string:
				m.Author = v
			case map[string]any:
				if s, ok := v["name"].(string); ok {
					m.Author = s
				}
			}
			m.Title = pkg.Plugin.Title
			m.Icon = pkg.Plugin.Icon
			m.Entry = pkg.Plugin.Entry
			m.CSS = pkg.Plugin.CSS
			m.Priority = pkg.Plugin.Priority
			m.Settings = pkg.Plugin.Settings
			m.Capabilities = pkg.Plugin.Capabilities
			m.Webhooks = pkg.Plugin.Webhooks
			m.Pipelines = pkg.Plugin.Pipelines
			m.MinHostVersion = pkg.Plugin.MinHostVersion
			m.TrustLevel = pkg.Plugin.TrustLevel
			m.ActivationEvents = pkg.Plugin.ActivationEvents
			m.AdminJS = pkg.Plugin.AdminJS
			m.PublicJS = pkg.Plugin.PublicJS
			m.Routes = pkg.Plugin.Routes
			m.Contributes = pkg.Plugin.Contributes
			m.Migrations = pkg.Plugin.Migrations
			m.Pages = pkg.Plugin.Pages
			m.Service = pkg.Plugin.Service
			homepage = pkg.Homepage
			resolved = true
			break
		}
	}

	if !resolved || m.Name == "" {
		return nil, fmt.Errorf("manifest not found: archive must contain plugin.yaml or package.json with a \"plugin\" field")
	}

	// Determine if this is a JS (Goja) plugin
	isJS := yamlType == "js" || yamlType == "full"
	if !isJS {
		// Also detect by presence of plugin.js in archive
		for name := range collected {
			if path.Base(name) == "plugin.js" {
				isJS = true
				break
			}
		}
	}

	// For JS/full plugins, preserve all files for extraction to disk
	if isJS {
		return &parseArchiveResult{
			Manifest: &m,
			AllFiles: collected,
			Homepage: homepage,
			IsJS:     true,
			JSEntry:  yamlJSEntry,
		}, nil
	}

	// Builtin (compiled Go) plugins: preserve all files so frontend assets
	// can be extracted to data/plugins/{id}/. The Go backend logic is compiled
	// into the binary and will be activated on next server restart.
	if yamlType == "builtin" {
		return &parseArchiveResult{
			Manifest:  &m,
			AllFiles:  collected,
			Homepage:  homepage,
			IsBuiltin: true,
		}, nil
	}

	// ── JS/YAML plugin path: locate entry script ─────────────────────────────
	if m.Entry == "" {
		m.Entry = "dist/index.js"
	}

	entryBase := path.Base(m.Entry)
	var script string
	for name, content := range collected {
		if name == m.Entry || path.Base(name) == entryBase {
			script = string(content)
			break
		}
	}
	// Layer 2/3 plugins may have no Goja entry script — that's OK
	if script == "" && m.AdminJS == "" && m.PublicJS == "" {
		return nil, fmt.Errorf("entry file %s not found in archive", m.Entry)
	}

	// ── Collect frontend asset files ─────────────────────────────────────────
	assets := make(map[string][]byte)
	allowedAssetExts := map[string]bool{".js": true, ".mjs": true, ".css": true, ".json": true}
	for name, content := range collected {
		ext := path.Ext(name)
		if !allowedAssetExts[ext] {
			continue
		}
		base := path.Base(name)
		if base == entryBase || base == "package.json" || base == "plugin.yaml" || base == "plugin.yml" {
			continue
		}
		assets[base] = content
	}

	return &parseArchiveResult{
		Manifest: &m,
		Script:   script,
		Assets:   assets,
		Homepage: homepage,
	}, nil
}

// parsePluginYAML converts a plugin.yaml content into an engine Manifest.
func parsePluginYAML(data []byte) (*eng.Manifest, string, error) {
	var y pluginYAMLManifest
	if err := yaml.Unmarshal(data, &y); err != nil {
		return nil, "", err
	}
	if y.ID == "" {
		return nil, "", fmt.Errorf("plugin.yaml missing required 'id' field")
	}

	m := &eng.Manifest{
		Name:        y.ID,
		Title:       y.Title,
		Version:     y.Version,
		Icon:        y.Icon,
		Author:      y.Author,
		Description: y.Description,
		TrustLevel:  eng.TrustLevel(y.TrustLevel),
		Type:        y.Type,
		AdminJS:     y.AdminJS,
		PublicJS:    y.PublicJS,
		CSS:         y.CSS,
		Entry:       y.Entry,
	}

	// Settings
	for _, s := range y.Settings {
		sf := eng.SettingField{
			Key:         s.Key,
			Label:       s.Label,
			Type:        eng.SettingType(s.Type),
			Required:    s.Required,
			Default:     s.Default,
			Placeholder: s.Placeholder,
			Description: s.Description,
			Options:     s.Options,
			Group:       s.Group,
			Shared:      s.Shared,
		}
		m.Settings = append(m.Settings, sf)
	}

	// Pages
	for _, p := range y.Pages {
		pd := eng.PageDef{
			Path:      p.Path,
			Slot:      p.Slot,
			Component: p.Component,
			Title:     p.Title,
		}
		if p.Nav != nil {
			pd.Nav = &eng.NavDef{
				Group: p.Nav.Group,
				Icon:  p.Nav.Icon,
				Order: p.Nav.Order,
			}
		}
		m.Pages = append(m.Pages, pd)
	}

	// Contributes
	if y.Contributes != nil {
		c := &eng.Contributes{}
		for _, cmd := range y.Contributes.Commands {
			c.Commands = append(c.Commands, eng.CommandDef{
				ID:    cmd.ID,
				Title: cmd.Title,
				Icon:  cmd.Icon,
			})
		}
		if len(y.Contributes.Menus) > 0 {
			c.Menus = make(map[string][]eng.MenuEntry)
			for slot, entries := range y.Contributes.Menus {
				for _, e := range entries {
					c.Menus[slot] = append(c.Menus[slot], eng.MenuEntry{Command: e.Command})
				}
			}
		}
		m.Contributes = c
	}

	// Migrations
	for _, mig := range y.Migrations {
		m.Migrations = append(m.Migrations, eng.MigrationDef{
			Version: mig.Version,
			Up:      mig.Up,
			Down:    mig.Down,
		})
	}

	// Routes
	for _, r := range y.Routes {
		m.Routes = append(m.Routes, eng.RouteDef{
			Method:      r.Method,
			Path:        r.Path,
			Auth:        r.Auth,
			Fn:          r.Fn,
			Description: r.Description,
		})
	}

	return m, y.Homepage, nil
}

// ── Plugin Assets (Phase 2.8) ──────────────────────────────────────────────

// PluginAssetsDir returns the directory where plugin frontend assets are stored.
func PluginAssetsDir() string { return filepath.Join("data", "plugins") }

// savePluginAssets writes frontend asset files to data/plugins/{sanitized-id}/.
func savePluginAssets(pluginID string, assets map[string][]byte) error {
	dir := filepath.Join(PluginAssetsDir(), sanitizePluginPath(pluginID))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create assets dir: %w", err)
	}
	for name, content := range assets {
		fp := filepath.Join(dir, name)
		if err := os.WriteFile(fp, content, 0644); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}
	return nil
}

// sanitizePluginPath converts a plugin ID like "nuxtblog/ai-polish" into
// a safe directory name "nuxtblog--ai-polish" to avoid nested directories.
func sanitizePluginPath(id string) string {
	return strings.ReplaceAll(id, "/", "--")
}

// ── GetSettings ───────────────────────────────────────────────────────────

func (s *sPlugin) GetSettings(ctx context.Context, id string) (*v1.PluginGetSettingsRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	type row struct {
		SettingsSchema string `orm:"settings_schema"`
		Settings       string `orm:"settings"`
	}
	var r row
	if err := g.DB().Ctx(ctx).Model("plugins").
		Where("id", id).
		Fields("settings_schema, settings").
		Scan(&r); err != nil {
		return nil, err
	}

	var schema []v1.SettingField
	if r.SettingsSchema != "" && r.SettingsSchema != "[]" {
		_ = json.Unmarshal([]byte(r.SettingsSchema), &schema)
	}
	if schema == nil {
		schema = []v1.SettingField{}
	}

	var values map[string]interface{}
	if r.Settings != "" && r.Settings != "{}" {
		_ = json.Unmarshal([]byte(r.Settings), &values)
	}
	if values == nil {
		values = map[string]interface{}{}
	}

	return &v1.PluginGetSettingsRes{Schema: schema, Values: values}, nil
}

// ── UpdateSettings ────────────────────────────────────────────────────────

func (s *sPlugin) UpdateSettings(ctx context.Context, id string, values map[string]interface{}) error {
	if err := requireAdmin(ctx); err != nil {
		return err
	}
	raw, err := json.Marshal(values)
	if err != nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, "invalid values")
	}
	_, err = g.DB().Ctx(ctx).Model("plugins").Where("id", id).
		Data(g.Map{"settings": string(raw), "updated_at": gtime.Now()}).Update()
	return err
}

// ── Uninstall ─────────────────────────────────────────────────────────────

func (s *sPlugin) Uninstall(ctx context.Context, id string) (bool, error) {
	if err := requireAdmin(ctx); err != nil {
		return false, err
	}

	// Look up source before deleting the row so we know whether this is a
	// builtin plugin that needs symmetric cleanup of builtin/{dirName}/.
	srcVal, _ := g.DB().Ctx(ctx).Model("plugins").Where("id", id).Value("source")
	src := srcVal.String()
	isBuiltin := src == "builtin"

	// Unload from Goja engine (legacy JS plugins)
	eng.Unload(id)

	// Unload from plugin manager (JS/builtin plugins) — cascade unloads dependents
	if mgr := eng.GetManager(); mgr != nil {
		unloaded := mgr.UnloadPlugin(id)
		if len(unloaded) > 1 {
			g.Log().Infof(ctx, "[plugin] cascade unloaded %d plugins: %v", len(unloaded), unloaded)
		}
	}

	// Remove plugin directory from data/plugins/{id}/
	pluginDir := filepath.Join(PluginAssetsDir(), sanitizePluginPath(id))
	if info, err := os.Stat(pluginDir); err == nil && info.IsDir() {
		if rmErr := os.RemoveAll(pluginDir); rmErr != nil {
			g.Log().Warningf(ctx, "[plugin] failed to remove plugin dir %s: %v", pluginDir, rmErr)
		}
	}

	// Rollback database migrations (drop plugin-created tables)
	s.rollbackPluginMigrations(ctx, id)

	// Also remove per-plugin KV store entries
	_, _ = g.DB().Ctx(ctx).Model("options").
		WhereLike("key", "plugin_store:"+id+":%").Delete()

	// Symmetric cleanup for builtin plugins: remove the source tree under
	// builtin/{dirName}/ and regenerate builtin/plugins.go so the next build
	// no longer compiles this plugin in. The currently running binary still
	// contains the compiled code — the warning tells the user to rebuild.
	if src == "builtin" {
		dirName := strings.ReplaceAll(id, "-", "_") // mirrors installBuiltinPlugin
		builtinDir := builtinSourceDir()
		if builtinDir != "" {
			target := filepath.Join(builtinDir, dirName)
			if info, err := os.Stat(target); err == nil && info.IsDir() {
				if rmErr := os.RemoveAll(target); rmErr != nil {
					g.Log().Warningf(ctx, "[plugin] failed to remove builtin source %s: %v", target, rmErr)
				}
			}
			if regenErr := regenerateBuiltinImports(builtinDir); regenErr != nil {
				g.Log().Warningf(ctx, "[plugin] failed to regenerate builtin imports: %v", regenErr)
			}
			g.Log().Warningf(ctx, "[plugin] builtin plugin %s source removed; rebuild and restart required for the compiled code to disappear", id)
		}
	}

	_, err := g.DB().Ctx(ctx).Model("plugins").Where("id", id).Delete()
	if err != nil {
		return false, err
	}
	return isBuiltin, nil
}

// rollbackPluginMigrations reads the plugin's manifest from the DB and runs
// down migrations to drop plugin-created tables, then cleans up the
// plugin_migrations tracking records.
func (s *sPlugin) rollbackPluginMigrations(ctx context.Context, id string) {
	manifestVal, _ := g.DB().Ctx(ctx).Model("plugins").
		Where("id", id).Value("manifest")
	manifestStr := manifestVal.String()
	if manifestStr == "" || manifestStr == "{}" {
		return
	}

	// Try to extract migrations from the stored manifest JSON.
	// The manifest may be either sdk.Manifest format (yaml-based) or
	// pluginsys.Manifest format (package.json-based). Both have a
	// "migrations" array with "version", "up", "down" fields.
	var mf struct {
		Migrations []eng.MigrationDef `json:"migrations"`
	}
	if err := json.Unmarshal([]byte(manifestStr), &mf); err != nil || len(mf.Migrations) == 0 {
		// No migrations to roll back — just clean up tracking records
		_, _ = g.DB().Ctx(ctx).Model("plugin_migrations").
			Where("plugin_id", id).Delete()
		return
	}

	if err := eng.RollbackMigrations(ctx, id, mf.Migrations); err != nil {
		g.Log().Warningf(ctx, "[plugin] %s migration rollback error: %v", id, err)
	} else {
		g.Log().Infof(ctx, "[plugin] %s: rolled back %d migration(s)", id, len(mf.Migrations))
	}
}

// ── Toggle ────────────────────────────────────────────────────────────────

func (s *sPlugin) Toggle(ctx context.Context, id string, enabled bool) error {
	if err := requireAdmin(ctx); err != nil {
		return err
	}

	enabledInt := 0
	if enabled {
		enabledInt = 1
	}
	_, err := g.DB().Ctx(ctx).Model("plugins").Where("id", id).
		Data(g.Map{"enabled": enabledInt, "updated_at": gtime.Now()}).Update()
	if err != nil {
		return err
	}

	if enabled {
		// Try to reload as JS (Goja) plugin first (check for plugin dir on disk)
		pluginDir := filepath.Join(PluginAssetsDir(), sanitizePluginPath(id))
		if mgr := eng.GetManager(); mgr != nil {
			if info, statErr := os.Stat(filepath.Join(pluginDir, "plugin.js")); statErr == nil && !info.IsDir() {
				// This is a JS (Goja) plugin — reload from disk
				_ = mgr.InstallJSPlugin(ctx, pluginDir, "plugin.js")
				return nil
			}
		}

		// Fallback: reload legacy script from DB
		type row struct {
			Script       string `orm:"script"`
			Priority     int    `orm:"priority"`
			Capabilities string `orm:"capabilities"`
			Manifest     string `orm:"manifest"`
		}
		var r row
		if e := g.DB().Ctx(ctx).Model("plugins").Where("id", id).
			Fields("script, priority, COALESCE(capabilities,'{}') as capabilities, COALESCE(manifest,'{}') as manifest").Scan(&r); e == nil && r.Script != "" {
			var caps eng.PluginCapabilities
			_ = json.Unmarshal([]byte(r.Capabilities), &caps)
			var mf eng.Manifest
			_ = json.Unmarshal([]byte(r.Manifest), &mf)
			mf.Priority = r.Priority
			mf.Capabilities = caps
			_ = eng.Load(id, r.Script, mf)
		}
	} else {
		// Unload from both engines — cascade unloads dependents
		normalized, _ := eng.NormalizePluginID(id)
		eng.Unload(normalized)
		if mgr := eng.GetManager(); mgr != nil {
			unloaded := mgr.UnloadPlugin(id)
			if len(unloaded) > 1 {
				g.Log().Infof(ctx, "[plugin] cascade unloaded %d plugins on disable: %v", len(unloaded), unloaded)
			}
		}
	}
	return nil
}

// ── Unload Impact ─────────────────────────────────────────────────────────

func (s *sPlugin) UnloadImpact(ctx context.Context, id string) (*v1.PluginUnloadImpactRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	var willUnload []string
	if mgr := eng.GetManager(); mgr != nil {
		willUnload = mgr.UnloadImpact(id)
	}
	if willUnload == nil {
		willUnload = []string{}
	}
	return &v1.PluginUnloadImpactRes{WillUnload: willUnload}, nil
}

// ── Marketplace ───────────────────────────────────────────────────────────

const (
	defaultRegistryURL  = "https://raw.githubusercontent.com/nuxtblog/registry/main/registry.json"
	registryCacheKey    = "plugin_registry"
	registrySyncedAtKey = "plugin_registry_synced_at"
	registryTTL         = time.Hour
)

// GetMarketplace returns the marketplace plugin list, auto-syncing if the cache is stale.
func (s *sPlugin) GetMarketplace(ctx context.Context, keyword, pluginType string) (*v1.MarketplaceListRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	syncedAt := getOption(ctx, registrySyncedAtKey)
	needsSync := syncedAt == ""
	if !needsSync {
		t, err := time.Parse(time.RFC3339, syncedAt)
		needsSync = err != nil || time.Since(t) > registryTTL
	}
	if needsSync {
		if res, err := s.SyncMarketplace(ctx); err == nil {
			syncedAt = res.SyncedAt
		}
	}

	var items []v1.MarketplaceItem
	if cache := getOption(ctx, registryCacheKey); cache != "" {
		_ = json.Unmarshal([]byte(cache), &items)
	}

	if keyword != "" || pluginType != "" {
		q := strings.ToLower(keyword)
		out := items[:0]
		for _, item := range items {
			if pluginType != "" && item.Type != pluginType {
				continue
			}
			if q != "" &&
				!strings.Contains(strings.ToLower(item.Title), q) &&
				!strings.Contains(strings.ToLower(item.Description), q) &&
				!strings.Contains(strings.ToLower(item.Author), q) {
				continue
			}
			out = append(out, item)
		}
		items = out
	}

	if items == nil {
		items = []v1.MarketplaceItem{}
	}
	return &v1.MarketplaceListRes{Items: items, SyncedAt: syncedAt}, nil
}

// SyncMarketplace fetches the latest registry.json and stores it in the options table.
func (s *sPlugin) SyncMarketplace(ctx context.Context) (*v1.MarketplaceSyncRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	regURL := g.Cfg().MustGet(ctx, "plugin.registry_url", defaultRegistryURL).String()
	body, err := pluginHTTPGet(ctx, regURL, true)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError, "fetch registry failed: "+err.Error())
	}

	var items []v1.MarketplaceItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError, "invalid registry format: "+err.Error())
	}

	now := time.Now().UTC().Format(time.RFC3339)
	setOption(ctx, registryCacheKey, string(body))
	setOption(ctx, registrySyncedAtKey, now)

	return &v1.MarketplaceSyncRes{Count: len(items), SyncedAt: now}, nil
}

func getOption(ctx context.Context, key string) string {
	type row struct{ Value string `orm:"value"` }
	var r row
	_ = g.DB().Ctx(ctx).Model("options").Where("key", key).Fields("value").Scan(&r)
	return r.Value
}

func setOption(ctx context.Context, key, value string) {
	cnt, _ := g.DB().Ctx(ctx).Model("options").Where("key", key).Count()
	if cnt > 0 {
		_, _ = g.DB().Ctx(ctx).Model("options").Where("key", key).
			Data(g.Map{"value": value}).Update()
	} else {
		_, _ = g.DB().Ctx(ctx).Model("options").
			Data(g.Map{"key": key, "value": value}).Insert()
	}
}

// ── GitHub proxy helpers ──────────────────────────────────────────────────

// applyGitHubProxy rewrites a GitHub URL using the configured mirror prefix.
// e.g. "https://api.github.com/..." → "https://ghproxy.net/https://api.github.com/..."
func applyGitHubProxy(ctx context.Context, rawURL string) string {
	proxy := strings.TrimRight(getOption(ctx, "plugin_github_proxy"), "/")
	if proxy == "" {
		return rawURL
	}
	for _, prefix := range []string{
		"https://api.github.com",
		"https://github.com",
		"https://raw.githubusercontent.com",
		"https://objects.githubusercontent.com",
	} {
		if strings.HasPrefix(rawURL, prefix) {
			return proxy + "/" + rawURL
		}
	}
	return rawURL
}

// pluginHTTPGet performs an HTTP GET like githubGet, but honours the
// plugin_http_proxy and plugin_github_proxy options.
func pluginHTTPGet(ctx context.Context, rawURL string, noCache ...bool) ([]byte, error) {
	targetURL := applyGitHubProxy(ctx, rawURL)

	transport := http.DefaultTransport.(*http.Transport).Clone()
	if proxyAddr := getOption(ctx, "plugin_http_proxy"); proxyAddr != "" {
		proxyURL, err := url.Parse(proxyAddr)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}
	client := &http.Client{Timeout: 60 * time.Second, Transport: transport}

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "nuxtblog-plugin-installer/1.0")
	req.Header.Set("Accept", "application/vnd.github+json")
	if len(noCache) > 0 && noCache[0] {
		req.Header.Set("Cache-Control", "no-cache")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("not found (404)")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// ── GetStyles ─────────────────────────────────────────────────────────────

// GetStyles returns the concatenated CSS of all enabled plugins.
// This endpoint is public — no auth required.
func (s *sPlugin) GetStyles(ctx context.Context) (string, error) {
	type row struct {
		Styles string `orm:"styles"`
	}
	var rows []row
	if err := g.DB().Ctx(ctx).Model("plugins").
		Fields("styles").Where("enabled", 1).
		Scan(&rows); err != nil {
		return "", err
	}
	var sb strings.Builder
	for _, r := range rows {
		if r.Styles != "" {
			sb.WriteString(r.Styles)
			sb.WriteByte('\n')
		}
	}
	return sb.String(), nil
}

// ── Stats / Errors (4.5-B1 / 4.5-B2 / 4.5-B3) ────────────────────────────

func (s *sPlugin) GetStats(ctx context.Context, id string) (*v1.PluginGetStatsRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	snap := eng.GetStats(id)
	if snap == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "plugin not loaded: "+id)
	}
	history := eng.GetHistory(id)
	buckets := make([]v1.WindowBucket, 0, len(history))
	for _, b := range history {
		buckets = append(buckets, v1.WindowBucket{
			Minute:      b.Minute.UTC().Format(time.RFC3339),
			Invocations: b.Invocations,
			Errors:      b.Errors,
		})
	}
	res := &v1.PluginGetStatsRes{
		PluginID:      snap.PluginID,
		Invocations:   snap.Invocations,
		Errors:        snap.Errors,
		AvgDurationMs: snap.AvgDurationMs,
		LastError:     snap.LastError,
		History:       buckets,
	}
	if !snap.LastErrorAt.IsZero() {
		res.LastErrorAt = snap.LastErrorAt.UTC().Format(time.RFC3339)
	}
	return res, nil
}

func (s *sPlugin) GetErrors(ctx context.Context, id string) (*v1.PluginGetErrorsRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	entries := eng.GetErrors(id)
	if entries == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "plugin not loaded: "+id)
	}
	items := make([]v1.PluginErrorEntry, 0, len(entries))
	for _, e := range entries {
		items = append(items, v1.PluginErrorEntry{
			At:        e.At.UTC().Format(time.RFC3339),
			EventName: e.EventName,
			Message:   e.Message,
			InputDiff: e.InputDiff,
		})
	}
	return &v1.PluginGetErrorsRes{Items: items}, nil
}

// ── Preview ───────────────────────────────────────────────────────────────

// Preview fetches plugin.yaml (or package.json) from the GitHub repo's default branch
// and returns parsed manifest metadata without downloading or installing the plugin.
func (s *sPlugin) Preview(ctx context.Context, repo string) (*v1.PluginPreviewRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	owner, repoName, err := parseRepo(repo)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, err.Error())
	}

	// Try plugin.yaml first
	yamlURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/HEAD/plugin.yaml", owner, repoName)
	if yamlData, yamlErr := pluginHTTPGet(ctx, yamlURL); yamlErr == nil {
		mf, _, parseErr := parsePluginYAML(yamlData)
		if parseErr == nil {
			res := buildPreviewFromManifest(mf)
			// Parse depends from raw YAML (eng.Manifest doesn't carry depends)
			var raw pluginYAMLManifest
			if yaml.Unmarshal(yamlData, &raw) == nil && len(raw.Depends) > 0 {
				for _, d := range raw.Depends {
					res.Depends = append(res.Depends, v1.DependencyPreview{
						ID: d.ID, Version: d.Version, Optional: d.Optional,
					})
				}
			}
			return res, nil
		}
	}

	// Fallback: package.json
	rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/HEAD/package.json", owner, repoName)
	data, err := pluginHTTPGet(ctx, rawURL)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError,
			fmt.Sprintf("cannot fetch plugin.yaml or package.json from %s/%s: %v", owner, repoName, err))
	}

	var pkg struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
		Author      any    `json:"author"`
		Plugin      *struct {
			Title        string                 `json:"title"`
			Icon         string                 `json:"icon"`
			CSS          string                 `json:"css"`
			Priority     int                    `json:"priority"`
			Settings     []eng.SettingField     `json:"settings"`
			Capabilities eng.PluginCapabilities `json:"capabilities"`
			Webhooks     []eng.WebhookDef       `json:"webhooks"`
			Pipelines    []eng.PipelineDef      `json:"pipelines"`
		} `json:"plugin"`
	}
	if err = json.Unmarshal(data, &pkg); err != nil || pkg.Plugin == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter,
			"package.json must contain a \"plugin\" field")
	}

	author := ""
	switch v := pkg.Author.(type) {
	case string:
		author = v
	case map[string]any:
		if s, ok2 := v["name"].(string); ok2 {
			author = s
		}
	}

	settings := make([]v1.SettingField, 0, len(pkg.Plugin.Settings))
	for _, sf := range pkg.Plugin.Settings {
		settings = append(settings, v1.SettingField{
			Key:         sf.Key,
			Label:       sf.Label,
			Type:        string(sf.Type),
			Required:    sf.Required,
			Default:     sf.Default,
			Placeholder: sf.Placeholder,
			Description: sf.Description,
			Options:     sf.Options,
			Shared:      sf.Shared,
		})
	}

	webhooks := make([]v1.WebhookPreview, 0, len(pkg.Plugin.Webhooks))
	for _, wh := range pkg.Plugin.Webhooks {
		webhooks = append(webhooks, v1.WebhookPreview{URL: wh.URL, Events: wh.Events})
	}

	pipelines := make([]v1.PipelinePreview, 0, len(pkg.Plugin.Pipelines))
	for _, p := range pkg.Plugin.Pipelines {
		pipelines = append(pipelines, v1.PipelinePreview{
			Name:      p.Name,
			Trigger:   p.Trigger,
			StepCount: len(p.Steps),
		})
	}

	caps := v1.CapabilitiesPreview{}
	if pkg.Plugin.Capabilities.HTTP != nil {
		caps.HTTP = &v1.HTTPCapPreview{
			Allow:     pkg.Plugin.Capabilities.HTTP.Allow,
			TimeoutMs: pkg.Plugin.Capabilities.HTTP.TimeoutMs,
		}
	}
	if pkg.Plugin.Capabilities.Store != nil {
		caps.Store = &v1.StoreCapPreview{
			Read:  pkg.Plugin.Capabilities.Store.Read,
			Write: pkg.Plugin.Capabilities.Store.Write,
		}
	}
	if pkg.Plugin.Capabilities.Events != nil {
		caps.Events = &v1.EventsCapPreview{
			Subscribe: pkg.Plugin.Capabilities.Events.Subscribe,
		}
	}
	if pkg.Plugin.Capabilities.DB != nil {
		caps.DB = convertDBCapPreview(pkg.Plugin.Capabilities.DB)
	}

	return &v1.PluginPreviewRes{
		Name:         pkg.Name,
		Title:        pkg.Plugin.Title,
		Description:  pkg.Description,
		Version:      pkg.Version,
		Author:       author,
		Icon:         orDefault(pkg.Plugin.Icon, "i-tabler-plug"),
		Priority:     pkg.Plugin.Priority,
		HasCSS:       pkg.Plugin.CSS != "",
		Capabilities: caps,
		Settings:     settings,
		Webhooks:     webhooks,
		Pipelines:    pipelines,
	}, nil
}

// buildPreviewFromManifest converts an engine Manifest into a PluginPreviewRes.
func buildPreviewFromManifest(mf *eng.Manifest) *v1.PluginPreviewRes {
	settings := make([]v1.SettingField, 0, len(mf.Settings))
	for _, sf := range mf.Settings {
		settings = append(settings, v1.SettingField{
			Key:         sf.Key,
			Label:       sf.Label,
			Type:        string(sf.Type),
			Required:    sf.Required,
			Default:     sf.Default,
			Placeholder: sf.Placeholder,
			Description: sf.Description,
			Options:     sf.Options,
			Shared:      sf.Shared,
		})
	}

	caps := v1.CapabilitiesPreview{}
	if mf.Capabilities.DB != nil {
		caps.DB = convertDBCapPreview(mf.Capabilities.DB)
	}

	return &v1.PluginPreviewRes{
		Name:         mf.Name,
		Title:        mf.Title,
		Description:  mf.Description,
		Version:      mf.Version,
		Author:       mf.Author,
		Icon:         orDefault(mf.Icon, "i-tabler-plug"),
		HasCSS:       mf.CSS != "",
		Capabilities: caps,
		Settings:     settings,
	}
}

// convertDBCapPreview converts engine DBCap to API DBCapPreview.
func convertDBCapPreview(d *eng.DBCap) *v1.DBCapPreview {
	if d == nil {
		return nil
	}
	r := &v1.DBCapPreview{
		Own: d.Own,
		Raw: d.Raw,
	}
	for _, t := range d.Tables {
		r.Tables = append(r.Tables, v1.DBTablePreview{
			Table: t.Table,
			Ops:   t.Ops,
		})
	}
	return r
}

// ── helpers ───────────────────────────────────────────────────────────────

func requireAdmin(ctx context.Context) error {
	if middleware.GetCurrentUserRole(ctx) < middleware.RoleAdmin {
		return gerror.NewCode(gcode.CodeNotAuthorized,
			g.I18n().T(ctx, "error.forbidden"))
	}
	return nil
}

// parseRepo normalises various GitHub URL formats to (owner, repo).
func parseRepo(raw string) (string, string, error) {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "https://")
	raw = strings.TrimPrefix(raw, "http://")
	raw = strings.TrimPrefix(raw, "github.com/")
	parts := strings.SplitN(raw, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repo format, expected owner/repo or github.com/owner/repo")
	}
	// Strip .git suffix and any trailing path
	repo := strings.TrimSuffix(strings.Split(parts[1], "/")[0], ".git")
	return parts[0], repo, nil
}


var (
	httpClient   = &http.Client{Timeout: 20 * time.Second}
	githubClient = &http.Client{Timeout: 60 * time.Second}
)

// httpGet is a plain HTTP GET used for non-GitHub URLs (e.g. registry sync).
func httpGet(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// githubGet makes a GitHub API / CDN request with the required User-Agent header.
// It follows redirects (needed for zipball URLs) and uses a longer timeout.
func githubGet(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "nuxtblog-plugin-installer/1.0")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := githubClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("not found (404)")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// manifestSettingsJSON serialises the plugin's settings schema to a JSON string.
func manifestSettingsJSON(m *eng.Manifest) string {
	if len(m.Settings) == 0 {
		return "[]"
	}
	b, _ := json.Marshal(m.Settings)
	return string(b)
}

// manifestDefaultSettings builds a JSON object pre-filled with each field's
// default value so plugins can call blog.settings.get(key) immediately.
func manifestDefaultSettings(m *eng.Manifest) string {
	defaults := make(map[string]interface{}, len(m.Settings))
	for _, f := range m.Settings {
		if f.Default != nil {
			defaults[f.Key] = f.Default
		}
	}
	b, _ := json.Marshal(defaults)
	return string(b)
}

// manifestCapabilitiesJSON serialises the plugin's capability declarations.
func manifestCapabilitiesJSON(m *eng.Manifest) string {
	b, _ := json.Marshal(m.Capabilities)
	if b == nil {
		return "{}"
	}
	return string(b)
}

// manifestFullJSON serialises the complete plugin manifest to a JSON string.
func manifestFullJSON(m *eng.Manifest) string {
	b, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// ── GetManifest / UpdateManifest (P-B11) ─────────────────────────────────────

func (s *sPlugin) GetManifest(ctx context.Context, id string) (*v1.PluginGetManifestRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	type row struct {
		Manifest string `orm:"manifest"`
	}
	var r row
	if err := g.DB().Ctx(ctx).Model("plugins").Where("id", id).
		Fields("COALESCE(manifest,'{}') as manifest").Scan(&r); err != nil {
		return nil, err
	}
	return &v1.PluginGetManifestRes{Manifest: r.Manifest}, nil
}

func (s *sPlugin) UpdateManifest(ctx context.Context, id string, manifest string) error {
	if err := requireAdmin(ctx); err != nil {
		return err
	}
	// Validate: must be valid JSON that deserialises to a Manifest
	var mf eng.Manifest
	if err := json.Unmarshal([]byte(manifest), &mf); err != nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, "manifest must be valid JSON: "+err.Error())
	}
	_, err := g.DB().Ctx(ctx).Model("plugins").Where("id", id).
		Data(g.Map{"manifest": manifest, "updated_at": gtime.Now()}).Update()
	if err != nil {
		return err
	}
	// Reload the plugin in the engine to apply pipeline/webhook changes
	type row struct {
		Script       string `orm:"script"`
		Priority     int    `orm:"priority"`
		Capabilities string `orm:"capabilities"`
	}
	var r row
	if e := g.DB().Ctx(ctx).Model("plugins").Where("id", id).
		Fields("script, priority, COALESCE(capabilities,'{}') as capabilities").Scan(&r); e == nil && r.Script != "" {
		var caps eng.PluginCapabilities
		_ = json.Unmarshal([]byte(r.Capabilities), &caps)
		mf.Priority = r.Priority
		mf.Capabilities = caps
		_ = eng.Load(id, r.Script, mf)
	}
	return nil
}

// ── ClientList (Phase 2.4) ────────────────────────────────────────────────

func (s *sPlugin) ClientList(ctx context.Context) (*v1.PluginClientListRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	type row struct {
		Id       string `orm:"id"`
		Title    string `orm:"title"`
		Icon     string `orm:"icon"`
		Manifest string `orm:"manifest"`
	}
	var rows []row
	if err := g.DB().Ctx(ctx).Model("plugins").
		Where("enabled", 1).
		Fields("id, title, icon, COALESCE(manifest,'{}') as manifest").
		Scan(&rows); err != nil {
		return nil, err
	}

	items := make([]v1.PluginClientItem, 0, len(rows))
	for _, r := range rows {
		var mf eng.Manifest
		_ = json.Unmarshal([]byte(r.Manifest), &mf)

		item := v1.PluginClientItem{
			ID:         r.Id,
			Title:      r.Title,
			Icon:       r.Icon,
			TrustLevel: string(mf.TrustLevel),
			AdminJS:    mf.AdminJS,
			PublicJS:   mf.PublicJS,
		}
		if item.TrustLevel == "" {
			item.TrustLevel = "community"
		}
		if mf.Contributes != nil {
			cb, _ := json.Marshal(mf.Contributes)
			item.Contributes = string(cb)
		}
		if len(mf.Pages) > 0 {
			pb, _ := json.Marshal(mf.Pages)
			item.Pages = string(pb)
		}
		items = append(items, item)
	}
	return &v1.PluginClientListRes{Items: items}, nil
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

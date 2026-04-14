package plugin

import (
	"context"
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
	candidate := filepath.Join(cwd, "builtin")
	if _, err := os.Stat(filepath.Join(candidate, "plugins.go")); err == nil {
		return candidate
	}
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

package plugin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// installJSPlugin extracts a JS (Goja) plugin to data/plugins/{id}/ and loads it.
func (s *sPlugin) installJSPlugin(ctx context.Context, ar *parseArchiveResult, repoUrl string) (*v1.PluginItem, error) {
	manifest := ar.Manifest
	pluginID := manifest.Name

	// Auto-install missing required dependencies from marketplace before extracting.
	if err := s.ensureDependenciesInstalled(ctx, pluginID, ar.Depends); err != nil {
		return nil, err
	}

	// Determine the directory prefix to strip from archive paths.
	prefix := detectArchivePrefix(ar.AllFiles)

	// Extract all files to data/plugins/{id}/
	pluginDir := filepath.Join(PluginAssetsDir(), sanitizePluginPath(pluginID))
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError,
			fmt.Sprintf("create plugin dir: %v", err))
	}

	for archivePath, content := range ar.AllFiles {
		relPath := archivePath
		if prefix != "" {
			relPath = strings.TrimPrefix(archivePath, prefix)
		}
		if relPath == "" {
			continue
		}
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

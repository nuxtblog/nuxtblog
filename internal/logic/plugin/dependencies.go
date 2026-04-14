package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

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

package oauth

import (
	"context"
)

var registry = map[string]Provider{}

// Register adds a builtin provider to the compiled registry. Call from each provider's init().
func Register(p Provider) {
	registry[p.Name()] = p
}

// GetProvider returns a provider by name, checking the compiled registry first,
// then falling back to generic (DB-configured) providers.
func GetProvider(ctx context.Context, name string) (Provider, bool) {
	// Builtin providers take priority
	if p, ok := registry[name]; ok {
		return p, ok
	}
	// Fall back to generic provider from DB
	if cfg := loadGenericConfig(ctx, name); cfg != nil {
		return &genericProvider{slug: name}, true
	}
	return nil, false
}

// IsEnabled checks whether a provider is enabled (reads DB options first, falls back to config.yaml).
func IsEnabled(ctx context.Context, name string) bool {
	return GetConfig(ctx, name).Enabled
}

// Enabled returns the names of all enabled providers (builtin + generic from DB).
func Enabled(ctx context.Context) []string {
	seen := map[string]bool{}
	var names []string

	// 1. Builtin compiled providers
	for name := range registry {
		if IsEnabled(ctx, name) {
			names = append(names, name)
			seen[name] = true
		}
	}

	// 2. Generic providers listed in oauth_providers option
	for _, slug := range GenericProviderSlugs(ctx) {
		if seen[slug] {
			continue
		}
		if IsEnabled(ctx, slug) {
			names = append(names, slug)
		}
	}

	return names
}

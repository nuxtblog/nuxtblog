package sdk

import (
	"fmt"
	"sync"

	"gopkg.in/yaml.v3"
)

// ParseManifest parses raw YAML bytes into a Manifest.
func ParseManifest(data []byte) (*Manifest, error) {
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse plugin.yaml: %w", err)
	}
	if m.ID == "" {
		return nil, fmt.Errorf("plugin.yaml: missing 'id' field")
	}
	if m.Priority == 0 {
		m.Priority = 10
	}
	return &m, nil
}

// MustParseManifest is like ParseManifest but panics on error.
// Intended for use with //go:embed in builtin plugins.
func MustParseManifest(data []byte) Manifest {
	m, err := ParseManifest(data)
	if err != nil {
		panic(err)
	}
	return *m
}

// manifestCache caches parsed manifests for //go:embed usage.
// Builtin plugins call ParseManifestCached to avoid re-parsing on every Manifest() call.
var (
	manifestCacheMu sync.RWMutex
	manifestCache   = make(map[string]*Manifest)
)

// ParseManifestCached parses and caches a manifest by a cache key.
// Subsequent calls with the same key return the cached result.
func ParseManifestCached(key string, data []byte) Manifest {
	manifestCacheMu.RLock()
	if m, ok := manifestCache[key]; ok {
		manifestCacheMu.RUnlock()
		return *m
	}
	manifestCacheMu.RUnlock()

	m := MustParseManifest(data)
	manifestCacheMu.Lock()
	manifestCache[key] = &m
	manifestCacheMu.Unlock()
	return m
}

package sdk

import "sync"

var (
	registryMu sync.Mutex
	registry   []Plugin
)

// Register adds a built-in plugin to the static registry.
// Called from init() in plugin packages imported via `import _`.
func Register(p Plugin) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry = append(registry, p)
}

// GetStatic returns all statically registered plugins.
func GetStatic() []Plugin {
	registryMu.Lock()
	defer registryMu.Unlock()
	out := make([]Plugin, len(registry))
	copy(out, registry)
	return out
}

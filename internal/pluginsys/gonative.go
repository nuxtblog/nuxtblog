package pluginsys

import "context"

// GoNativeManager is the interface that the plugin manager must implement.
// Registered via RegisterGoNativeManager() before startup. Manages both
// compiled Go (builtin) plugins and interpreted JS (Goja) plugins.
type GoNativeManager interface {
	// HasPlugin reports whether a Go-native plugin with the given ID is registered.
	HasPlugin(id string) bool
	// FanOutEvent dispatches a fire-and-forget event to all Go-native plugins.
	FanOutEvent(ctx context.Context, event string, data map[string]any)
	// RunFilter runs Go-native plugin filters for the given event.
	RunFilter(ctx context.Context, event string, data, meta map[string]any) error
}

var goNativeMgr GoNativeManager

// RegisterGoNativeManager registers the Go-native plugin manager.
// Must be called before LoadAll().
func RegisterGoNativeManager(mgr GoNativeManager) {
	goNativeMgr = mgr
}

// isGoNativePlugin reports whether a plugin ID has a Go-native implementation.
func isGoNativePlugin(id string) bool {
	if goNativeMgr == nil {
		return false
	}
	return goNativeMgr.HasPlugin(id)
}

// Package pluginsys provides a multi-layer plugin system:
//
//   - Layer 0: YAML declarative plugins (yaml.go)
//   - Layer 1: JavaScript plugins via Goja (goja.go, goja_bridge.go)
//   - Layer 2: Go-native compiled plugins (manager.go)
//
// The former Layer 1 (Goja JavaScript engine) has been removed.
package pluginsys

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// normalizePluginID ensures the plugin ID uses the canonical "owner/slug" format.
// If id has no "/" separator, it is prefixed with "unknown/" and the second
// return value is true to signal that the caller should emit a deprecation warning.
func NormalizePluginID(id string) (string, bool) {
	if !strings.Contains(id, "/") {
		return "unknown/" + id, true
	}
	return id, false
}

// Load is a no-op retained for backward compatibility with callers that
// previously loaded Goja plugins. The Goja engine has been removed.
func Load(id, script string, mf Manifest) error {
	return nil
}

// Unload is a no-op retained for backward compatibility with callers that
// previously unloaded Goja plugins. The Goja engine has been removed.
func Unload(id string) {}

// ─── Public Filter API ────────────────────────────────────────────────────────

// Filter runs all plugin filter handlers for eventName synchronously.
// It builds a PluginCtx internally; the public API remains map[string]any.
// Returns modified data or an error if any handler aborts / throws / times out.
// Runs YAML and Go-native plugin filters in order.
func Filter(ctx context.Context, eventName string, data map[string]any) (map[string]any, error) {
	// YAML declarative plugin filters (Layer 0)
	if err := RunYAMLFilters(ctx, eventName, data); err != nil {
		return nil, err
	}
	// Go-native plugin filters (Layer 2)
	if goNativeMgr != nil {
		meta := make(map[string]any)
		if err := goNativeMgr.RunFilter(ctx, eventName, data, meta); err != nil {
			return nil, err
		}
	}
	return data, nil
}

// ─── Event dispatch ───────────────────────────────────────────────────────────

// fanOut dispatches an event to YAML webhooks and Go-native plugins.
func fanOut(eventName string, payloadMap map[string]any) {
	// Dispatch YAML webhooks (Layer 0)
	FanOutYAMLWebhooks(eventName, payloadMap)

	// Dispatch to Go-native plugins (Layer 2)
	if goNativeMgr != nil {
		goNativeMgr.FanOutEvent(context.Background(), eventName, payloadMap)
	}
}

// deepCopyMap returns a deep copy of m via JSON round-trip.
func deepCopyMap(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	b, _ := json.Marshal(m)
	var out map[string]any
	_ = json.Unmarshal(b, &out)
	return out
}

// isVersionCompatible reports whether current >= required (simple semver comparison).
// Both strings must be in "MAJOR.MINOR.PATCH" form; malformed strings return false.
func isVersionCompatible(current, required string) bool {
	parse := func(s string) (major, minor, patch int, ok bool) {
		_, err := fmt.Sscanf(s, "%d.%d.%d", &major, &minor, &patch)
		return major, minor, patch, err == nil
	}
	cMaj, cMin, cPat, ok1 := parse(current)
	rMaj, rMin, rPat, ok2 := parse(required)
	if !ok1 || !ok2 {
		return false
	}
	if cMaj != rMaj {
		return cMaj > rMaj
	}
	if cMin != rMin {
		return cMin > rMin
	}
	return cPat >= rPat
}

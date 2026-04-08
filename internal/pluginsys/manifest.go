package pluginsys

import (
	"github.com/nuxtblog/nuxtblog/sdk"
)

// parsePluginYAML parses raw YAML bytes into the unified sdk.Manifest.
func parsePluginYAML(data []byte) (*sdk.Manifest, error) {
	return sdk.ParseManifest(data)
}

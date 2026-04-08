package payload

// PluginInstalled is delivered when a plugin is installed or upgraded.
type PluginInstalled struct {
	PluginID string
	Title    string
	Version  string
	Author   string
}

// PluginUninstalled is delivered when a plugin is uninstalled.
type PluginUninstalled struct {
	PluginID string
}

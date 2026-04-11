package service

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
)

type IPlugin interface {
	List(ctx context.Context) (*v1.PluginListRes, error)
	Install(ctx context.Context, repoUrl, expectedVersion string) (*v1.PluginItem, error)
	InstallZip(ctx context.Context, zipData []byte) (*v1.PluginItem, error)
	Update(ctx context.Context, id string) (*v1.PluginItem, error)
	BatchUpdate(ctx context.Context, ids []string) (*v1.PluginBatchUpdateRes, error)
	Uninstall(ctx context.Context, id string) (needRestart bool, err error)
	Toggle(ctx context.Context, id string, enabled bool) error
	GetSettings(ctx context.Context, id string) (*v1.PluginGetSettingsRes, error)
	UpdateSettings(ctx context.Context, id string, values map[string]interface{}) error
	GetStyles(ctx context.Context) (string, error)
	GetMarketplace(ctx context.Context, keyword, pluginType string) (*v1.MarketplaceListRes, error)
	SyncMarketplace(ctx context.Context) (*v1.MarketplaceSyncRes, error)
	GetStats(ctx context.Context, id string) (*v1.PluginGetStatsRes, error)
	GetErrors(ctx context.Context, id string) (*v1.PluginGetErrorsRes, error)
	GetManifest(ctx context.Context, id string) (*v1.PluginGetManifestRes, error)
	UpdateManifest(ctx context.Context, id string, manifest string) error
	Preview(ctx context.Context, repo string) (*v1.PluginPreviewRes, error)
	ClientList(ctx context.Context) (*v1.PluginClientListRes, error)
	PublicClientList(ctx context.Context) (*v1.PluginPublicClientRes, error)
	UnloadImpact(ctx context.Context, id string) (*v1.PluginUnloadImpactRes, error)
}

var _plugin IPlugin

func Plugin() IPlugin           { return _plugin }
func RegisterPlugin(p IPlugin) { _plugin = p }

package plugin

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
)

type IPlugin interface {
	PluginList(ctx context.Context, req *v1.PluginListReq) (res *v1.PluginListRes, err error)
	PluginInstall(ctx context.Context, req *v1.PluginInstallReq) (res *v1.PluginInstallRes, err error)
	PluginUploadZip(ctx context.Context, req *v1.PluginUploadZipReq) (res *v1.PluginUploadZipRes, err error)
	PluginUninstall(ctx context.Context, req *v1.PluginUninstallReq) (res *v1.PluginUninstallRes, err error)
	PluginToggle(ctx context.Context, req *v1.PluginToggleReq) (res *v1.PluginToggleRes, err error)
	PluginGetSettings(ctx context.Context, req *v1.PluginGetSettingsReq) (res *v1.PluginGetSettingsRes, err error)
	PluginUpdateSettings(ctx context.Context, req *v1.PluginUpdateSettingsReq) (res *v1.PluginUpdateSettingsRes, err error)
	PluginGetManifest(ctx context.Context, req *v1.PluginGetManifestReq) (res *v1.PluginGetManifestRes, err error)
	PluginUpdateManifest(ctx context.Context, req *v1.PluginUpdateManifestReq) (res *v1.PluginUpdateManifestRes, err error)
	PluginClientList(ctx context.Context, req *v1.PluginClientListReq) (res *v1.PluginClientListRes, err error)
}

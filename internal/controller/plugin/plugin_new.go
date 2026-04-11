package plugin

import apiplugin "github.com/nuxtblog/nuxtblog/api/plugin"

func New() apiplugin.IPlugin       { return &ControllerV1{} }
func NewPublic() apiplugin.IPluginPublic { return &ControllerV1{} }

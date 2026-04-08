package plugin

import apiplugin "github.com/nuxtblog/nuxtblog/api/plugin"

func New() apiplugin.IPlugin { return &ControllerV1{} }

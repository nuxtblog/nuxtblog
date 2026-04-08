package nav_menu

import api "github.com/nuxtblog/nuxtblog/api/nav_menu"

type ControllerV1 struct{}

func NewV1() api.INavMenuV1 { return &ControllerV1{} }

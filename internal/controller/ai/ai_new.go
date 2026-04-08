package ai

import apiai "github.com/nuxtblog/nuxtblog/api/ai"

func New() apiai.IAI { return &ControllerV1{} }

type ControllerV1 struct{}

package history

import apihistory "github.com/nuxtblog/nuxtblog/api/history"

func New() apihistory.IHistory { return &ControllerV1{} }

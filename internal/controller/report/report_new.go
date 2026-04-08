package report

import apireport "github.com/nuxtblog/nuxtblog/api/report"

func New() apireport.IReport { return &ControllerV1{} }

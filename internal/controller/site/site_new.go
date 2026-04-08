package site

import apisite "github.com/nuxtblog/nuxtblog/api/site"

type ControllerV1 struct{}

func NewV1() apisite.ISiteV1 { return &ControllerV1{} }

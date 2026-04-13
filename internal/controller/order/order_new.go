package order

import apiorder "github.com/nuxtblog/nuxtblog/api/order"

func New() apiorder.IOrderV1 { return &ControllerV1{} }

type ControllerV1 struct{}

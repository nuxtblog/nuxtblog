package token

import apitoken "github.com/nuxtblog/nuxtblog/api/token"

type ControllerV1 struct{}

func NewV1() apitoken.ITokenV1 {
	return &ControllerV1{}
}

package verifycode

import apiverifycode "github.com/nuxtblog/nuxtblog/api/verifycode"

type ControllerV1 struct{}

func NewV1() apiverifycode.IVerifycodeV1 {
	return &ControllerV1{}
}

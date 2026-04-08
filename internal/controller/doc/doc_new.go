package doc

import "github.com/nuxtblog/nuxtblog/api/doc"

func NewV1() doc.IDocV1 {
	return &ControllerV1{}
}

package system

import (
	"github.com/nuxtblog/nuxtblog/api/system"
)

type ControllerV1 struct{}

func NewV1() system.ISystemV1 {
	return &ControllerV1{}
}

package reaction

import apireaction "github.com/nuxtblog/nuxtblog/api/reaction"

type ControllerV1 struct{}

func NewV1() apireaction.IReactionV1 {
	return &ControllerV1{}
}

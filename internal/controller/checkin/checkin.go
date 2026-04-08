package checkin

import apicheckin "github.com/nuxtblog/nuxtblog/api/checkin"

type ControllerV1 struct{}

func NewV1() apicheckin.ICheckinV1 {
	return &ControllerV1{}
}

package moment

import "github.com/nuxtblog/nuxtblog/api/moment"

func NewV1() moment.IMomentV1 {
	return &ControllerV1{}
}

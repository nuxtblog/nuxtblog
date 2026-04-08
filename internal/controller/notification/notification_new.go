package notification

import "github.com/nuxtblog/nuxtblog/api/notification"

type ControllerV1 struct{}

func NewV1() notification.INotificationV1 {
	return &ControllerV1{}
}

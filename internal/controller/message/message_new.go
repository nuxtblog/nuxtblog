package message

import apimessage "github.com/nuxtblog/nuxtblog/api/message"

func New() apimessage.IMessage { return &ControllerV1{} }

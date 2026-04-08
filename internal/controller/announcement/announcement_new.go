package announcement

import apiannounc "github.com/nuxtblog/nuxtblog/api/announcement"

func NewAdmin() apiannounc.IAnnouncementAdmin   { return &AdminControllerV1{} }
func NewPublic() apiannounc.IAnnouncementPublic { return &PublicControllerV1{} }

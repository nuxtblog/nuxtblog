package friendlink

import apifriendlink "github.com/nuxtblog/nuxtblog/api/friendlink"

func NewAdmin() apifriendlink.IFriendlinkAdmin   { return &AdminControllerV1{} }
func NewPublic() apifriendlink.IFriendlinkPublic { return &PublicControllerV1{} }

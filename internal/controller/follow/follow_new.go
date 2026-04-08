package follow

import apifollow "github.com/nuxtblog/nuxtblog/api/follow"

// PublicControllerV1 handles routes that don't require authentication.
type PublicControllerV1 struct{}

// AuthControllerV1 handles routes that require a valid JWT.
type AuthControllerV1 struct{}

func NewPublicV1() apifollow.IFollowPublicV1 { return &PublicControllerV1{} }
func NewAuthV1() apifollow.IFollowAuthV1     { return &AuthControllerV1{} }

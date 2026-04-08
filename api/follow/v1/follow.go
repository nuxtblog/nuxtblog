package v1

import "github.com/gogf/gf/v2/frame/g"

// FollowUserItem is the user representation returned in follower/following lists.
type FollowUserItem struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	DisplayName    string `json:"display_name"`
	Avatar         string `json:"avatar"`
	Bio            string `json:"bio"`
	FollowedAt     string `json:"followed_at"`
	ArticleCount   int    `json:"article_count"`
	FollowerCount  int    `json:"follower_count"`
	IsFollowingBack bool  `json:"is_following_back"`
}

// ----------------------------------------------------------------
//  Public — Followers list
// ----------------------------------------------------------------

type UserFollowersReq struct {
	g.Meta `path:"/users/{id}/followers" method:"get" tags:"Follow" summary:"Get user's followers"`
	Id     int64 `p:"id"   v:"required|min:1"`
	Page   int   `p:"page" d:"1"`
	Size   int   `p:"size" d:"20"`
}
type UserFollowersRes struct {
	List  []FollowUserItem `json:"list"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

// ----------------------------------------------------------------
//  Public — Following list
// ----------------------------------------------------------------

type UserFollowingReq struct {
	g.Meta `path:"/users/{id}/following" method:"get" tags:"Follow" summary:"Get users followed by a user"`
	Id     int64 `p:"id"   v:"required|min:1"`
	Page   int   `p:"page" d:"1"`
	Size   int   `p:"size" d:"20"`
}
type UserFollowingRes struct {
	List  []FollowUserItem `json:"list"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

// ----------------------------------------------------------------
//  Auth — Follow
// ----------------------------------------------------------------

type UserFollowReq struct {
	g.Meta `path:"/users/{id}/follow" method:"post" tags:"Follow" summary:"Follow a user"`
	Id     int64 `p:"id" v:"required|min:1"`
}
type UserFollowRes struct {
	Following bool `json:"following"`
}

// ----------------------------------------------------------------
//  Auth — Unfollow
// ----------------------------------------------------------------

type UserUnfollowReq struct {
	g.Meta `path:"/users/{id}/follow" method:"delete" tags:"Follow" summary:"Unfollow a user"`
	Id     int64 `p:"id" v:"required|min:1"`
}
type UserUnfollowRes struct {
	Following bool `json:"following"`
}

// ----------------------------------------------------------------
//  Auth — Remove follower
// ----------------------------------------------------------------

type UserRemoveFollowerReq struct {
	g.Meta `path:"/users/{id}/follower" method:"delete" tags:"Follow" summary:"Remove a user from your followers"`
	Id     int64 `p:"id" v:"required|min:1"`
}
type UserRemoveFollowerRes struct {
	Removed bool `json:"removed"`
}

// ----------------------------------------------------------------
//  Auth — Follow status
// ----------------------------------------------------------------

type UserFollowStatusReq struct {
	g.Meta `path:"/users/{id}/follow-status" method:"get" tags:"Follow" summary:"Get follow relationship status"`
	Id     int64 `p:"id" v:"required|min:1"`
}
type UserFollowStatusRes struct {
	Following      bool `json:"following"`
	FollowerCount  int  `json:"follower_count"`
	FollowingCount int  `json:"following_count"`
}

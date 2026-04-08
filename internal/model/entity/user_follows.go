// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import "github.com/gogf/gf/v2/os/gtime"

// UserFollows is the golang structure for table user_follows.
type UserFollows struct {
	FollowerId  int64       `json:"followerId"  orm:"follower_id"  description:""`
	FollowingId int64       `json:"followingId" orm:"following_id" description:""`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""`
}

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ----------------------------------------------------------------
//  Enums
// ----------------------------------------------------------------

type CommentReviewStatus int

const (
	CommentReviewPending  CommentReviewStatus = 1
	CommentReviewApproved CommentReviewStatus = 2
	CommentReviewSpam     CommentReviewStatus = 3
	CommentReviewTrashed  CommentReviewStatus = 4
)

// ----------------------------------------------------------------
//  Shared output
// ----------------------------------------------------------------

type CommentItem struct {
	Id          int64               `json:"id"`
	ObjectId    int64               `json:"object_id"`
	ObjectType  string              `json:"object_type"`
	ParentId    *int64              `json:"parent_id"`
	UserId      *int64              `json:"user_id"`
	AuthorName  string              `json:"author_name"`
	AuthorEmail string              `json:"author_email,omitempty"`
	Content     string              `json:"content"`
	Status      CommentReviewStatus `json:"status"`
	CreatedAt   *gtime.Time         `json:"created_at"`
	Replies     []*CommentItem      `json:"replies,omitempty"`
}

// CommentAdminItem is the enriched comment item for admin endpoints
type CommentAdminItem struct {
	CommentId int64               `json:"comment_id"`
	PostId    int64               `json:"post_id"`
	Content   string              `json:"content"`
	ParentId  *int64              `json:"parent_id,omitempty"`
	Status    string              `json:"status"`
	IpAddress string              `json:"ip_address,omitempty"`
	UserAgent string              `json:"user_agent,omitempty"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
	Author    *CommentAdminAuthor `json:"author,omitempty"`
	Children  []*CommentAdminItem `json:"children,omitempty"`
}

type CommentAdminAuthor struct {
	Id     *int64 `json:"id,omitempty"`
	Name   string `json:"name"`
	Email  string `json:"email,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// ----------------------------------------------------------------
//  Create (public — visitor submit)
// ----------------------------------------------------------------

type CommentCreateReq struct {
	g.Meta      `path:"/comments" method:"post" tags:"Comment" summary:"Submit a comment"`
	ObjectId    int64  `v:"required|min:1"                  dc:"object id, e.g. post id"`
	ObjectType  string `v:"required|in:post,product,video"  dc:"object type"`
	ParentId    *int64 `v:"min:1"                           dc:"parent comment id for replies"`
	AuthorName  string `v:"required|length:1,100"           dc:"commenter name"`
	AuthorEmail string `v:"required|email"                  dc:"commenter email, not shown publicly"`
	Content     string `v:"required|length:1,2000"          dc:"comment content"`
}
type CommentCreateRes struct {
	Id     int64               `json:"id"`
	Status CommentReviewStatus `json:"status" dc:"pending if moderation is on"`
}

// ----------------------------------------------------------------
//  Delete
// ----------------------------------------------------------------

type CommentDeleteReq struct {
	g.Meta `path:"/comments/{id}" method:"delete" tags:"Comment" summary:"Delete comment"`
	Id     int64 `v:"required|min:1" dc:"comment id"`
}
type CommentDeleteRes struct{}

// ----------------------------------------------------------------
//  Update status (admin — moderate)
// ----------------------------------------------------------------

type CommentUpdateStatusReq struct {
	g.Meta `path:"/comments/{id}/status" method:"put" tags:"Comment" summary:"Update comment status"`
	Id     int64               `v:"required|min:1"  dc:"comment id"`
	Status CommentReviewStatus `v:"required|in:1,2,3,4" dc:"1=pending 2=approved 3=spam 4=trash"`
}
type CommentUpdateStatusRes struct{}

// ----------------------------------------------------------------
//  Update content (admin)
// ----------------------------------------------------------------

type CommentUpdateContentReq struct {
	g.Meta  `path:"/comments/{id}" method:"put" tags:"Comment" summary:"Admin: update comment content"`
	Id      int64  `v:"required|min:1"        dc:"comment id"`
	Content string `v:"required|length:1,2000" dc:"new comment content"`
}
type CommentUpdateContentRes struct{}

// ----------------------------------------------------------------
//  Get list for object (public)
// ----------------------------------------------------------------

type CommentGetListReq struct {
	g.Meta     `path:"/comments" method:"get" tags:"Comment" summary:"Get comments for an object"`
	ObjectId   int64  `v:"required|min:1"                 dc:"object id"`
	ObjectType string `v:"required|in:post,product,video" dc:"object type"`
	Page       int    `v:"min:1"                          dc:"page number" d:"1"`
	Size       int    `v:"between:1,100"                  dc:"page size"   d:"20"`
}
type CommentGetListRes struct {
	List  []*CommentItem `json:"list"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

// ----------------------------------------------------------------
//  Get list for admin (all objects, filterable)
// ----------------------------------------------------------------

type CommentAdminGetListReq struct {
	g.Meta     `path:"/admin/comments" method:"get" tags:"Comment" summary:"Admin: get all comments"`
	ObjectType *string `v:"in:post,product,video" dc:"filter by object type"`
	Status     *string `                          dc:"filter: pending/approved/spam/trash or 1/2/3/4"`
	Keyword    *string `                          dc:"search content or author"`
	Page       int     `v:"min:1"                 dc:"page number" d:"1"`
	Size       int     `v:"between:1,100"         dc:"page size"   d:"20"`
}
type CommentAdminGetListRes struct {
	List  []*CommentAdminItem `json:"list"`
	Total int                 `json:"total"`
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
}

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ----------------------------------------------------------------
//  Like
// ----------------------------------------------------------------

type LikePostReq struct {
	g.Meta `path:"/posts/{id}/like" method:"post" tags:"Reaction" summary:"Like a post"`
	Id     int64 `v:"required|min:1" in:"path" dc:"post id"`
}
type LikePostRes struct{}

type UnlikePostReq struct {
	g.Meta `path:"/posts/{id}/like" method:"delete" tags:"Reaction" summary:"Unlike a post"`
	Id     int64 `v:"required|min:1" in:"path" dc:"post id"`
}
type UnlikePostRes struct{}

// ----------------------------------------------------------------
//  Bookmark
// ----------------------------------------------------------------

type BookmarkPostReq struct {
	g.Meta `path:"/posts/{id}/bookmark" method:"post" tags:"Reaction" summary:"Bookmark a post"`
	Id     int64 `v:"required|min:1" in:"path" dc:"post id"`
}
type BookmarkPostRes struct{}

type UnbookmarkPostReq struct {
	g.Meta `path:"/posts/{id}/bookmark" method:"delete" tags:"Reaction" summary:"Unbookmark a post"`
	Id     int64 `v:"required|min:1" in:"path" dc:"post id"`
}
type UnbookmarkPostRes struct{}

// ----------------------------------------------------------------
//  Get reaction status for current user
// ----------------------------------------------------------------

type GetReactionReq struct {
	g.Meta `path:"/posts/{id}/reaction" method:"get" tags:"Reaction" summary:"Get current user reaction status"`
	Id     int64 `v:"required|min:1" in:"path" dc:"post id"`
}
type GetReactionRes struct {
	Liked      bool `json:"liked"`
	Bookmarked bool `json:"bookmarked"`
}

// ----------------------------------------------------------------
//  Bookmark list item
// ----------------------------------------------------------------

type BookmarkPostItem struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Excerpt   string `json:"excerpt"`
	CreatedAt string `json:"created_at"`
}

// ----------------------------------------------------------------
//  Get bookmarks list (user's saved posts)
// ----------------------------------------------------------------

type GetBookmarksReq struct {
	g.Meta `path:"/users/me/bookmarks" method:"get" tags:"Reaction" summary:"Get current user's bookmarked posts"`
	Page   int `v:"min:1" d:"1"  dc:"page number"`
	Size   int `v:"min:1|max:50" d:"20" dc:"page size"`
}
type GetBookmarksRes struct {
	List  []*BookmarkPostItem `json:"list"`
	Total int                 `json:"total"`
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
}

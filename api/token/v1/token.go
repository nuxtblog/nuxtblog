package v1

import "github.com/gogf/gf/v2/frame/g"

// UserTokenItem is a safe token representation (hash is never exposed).
type UserTokenItem struct {
	Id         int64   `json:"id"`
	Name       string  `json:"name"`
	Prefix     string  `json:"prefix"`
	ExpiresAt  *string `json:"expires_at"`
	LastUsedAt *string `json:"last_used_at"`
	CreatedAt  string  `json:"created_at"`
}

// ----------------------------------------------------------------
//  List
// ----------------------------------------------------------------

type UserTokenListReq struct {
	g.Meta `path:"/users/tokens" method:"get" tags:"Token" summary:"List personal API tokens"`
}
type UserTokenListRes struct {
	List []UserTokenItem `json:"list"`
}

// ----------------------------------------------------------------
//  Create
// ----------------------------------------------------------------

type UserTokenCreateReq struct {
	g.Meta        `path:"/users/tokens" method:"post" tags:"Token" summary:"Create a personal API token"`
	Name          string `p:"name"            v:"required"  dc:"token label"`
	ExpiresInDays int    `p:"expires_in_days"               dc:"validity in days; 0 = never expires"`
}
type UserTokenCreateRes struct {
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
	Prefix    string  `json:"prefix"`
	Token     string  `json:"token"      dc:"shown only once at creation"`
	ExpiresAt *string `json:"expires_at"`
	CreatedAt string  `json:"created_at"`
}

// ----------------------------------------------------------------
//  Delete
// ----------------------------------------------------------------

type UserTokenDeleteReq struct {
	g.Meta `path:"/users/tokens/{id}" method:"delete" tags:"Token" summary:"Delete a personal API token"`
	Id     int64 `p:"id" v:"required|min:1"`
}
type UserTokenDeleteRes struct{}

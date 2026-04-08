package v1

import "github.com/gogf/gf/v2/frame/g"

// ----------------------------------------------------------------
//  Shared output
// ----------------------------------------------------------------

type OptionItem struct {
	Key      string `json:"key"`
	Value    string `json:"value"` // raw JSON value
	Autoload int    `json:"autoload"`
}

// ----------------------------------------------------------------
//  Get one
// ----------------------------------------------------------------

type OptionGetReq struct {
	g.Meta `path:"/options/{key}" method:"get" tags:"Option" summary:"Get option by key"`
	Key    string `v:"required|length:1,191" dc:"option key"`
}
type OptionGetRes struct {
	*OptionItem `dc:"option"`
}

// ----------------------------------------------------------------
//  Get batch (autoload options — called on app init)
// ----------------------------------------------------------------

type OptionGetAutoloadReq struct {
	g.Meta `path:"/options/autoload" method:"get" tags:"Option" summary:"Get all autoload options"`
}
type OptionGetAutoloadRes struct {
	Options map[string]string `json:"options" dc:"key → raw JSON value"`
}

// ----------------------------------------------------------------
//  Set (upsert)
// ----------------------------------------------------------------

type OptionSetReq struct {
	g.Meta   `path:"/options/{key}" method:"put" tags:"Option" summary:"Set option value"`
	Key      string `v:"required|length:1,191" dc:"option key"`
	Value    string `v:"required"              dc:"JSON-encoded value"`
	Autoload *int   `v:"in:0,1"                dc:"0=no 1=yes"`
}
type OptionSetRes struct{}

// ----------------------------------------------------------------
//  Delete
// ----------------------------------------------------------------

type OptionDeleteReq struct {
	g.Meta `path:"/options/{key}" method:"delete" tags:"Option" summary:"Delete option"`
	Key    string `v:"required|length:1,191" dc:"option key"`
}
type OptionDeleteRes struct{}

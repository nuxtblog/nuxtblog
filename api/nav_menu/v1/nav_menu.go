package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// NavMenuItemInput is sent by the client when saving a menu.
// Items are in display order. ParentIdx=-1 means root; >=0 means index of parent in the same array.
type NavMenuItemInput struct {
	ParentIdx  int    `json:"parent_idx"`  // -1 = root
	ObjectType string `json:"object_type"` // custom|page|category
	ObjectId   int64  `json:"object_id"`
	Label      string `json:"label"`
	Url        string `json:"url"`
	Target     string `json:"target"` // '' or '_blank'
	CssClasses string `json:"css_classes"`
}

type NavMenuItemOutput struct {
	Id         int64  `json:"id"`
	MenuId     int64  `json:"menu_id"`
	ParentId   int64  `json:"parent_id"`
	ObjectType string `json:"object_type"`
	ObjectId   int64  `json:"object_id"`
	Label      string `json:"label"`
	Url        string `json:"url"`
	Target     string `json:"target"`
	CssClasses string `json:"css_classes"`
	SortOrder  int    `json:"sort_order"`
}

type NavMenuOutput struct {
	Id          int64                `json:"id"`
	Name        string               `json:"name"`
	Location    string               `json:"location"`
	Description string               `json:"description"`
	Items       []*NavMenuItemOutput `json:"items"`
	CreatedAt   *gtime.Time          `json:"created_at"`
	UpdatedAt   *gtime.Time          `json:"updated_at"`
}

// List
type NavMenuListReq struct {
	g.Meta `path:"/nav-menus" method:"get" tags:"NavMenu" summary:"List all menus with items"`
}
type NavMenuListRes struct {
	List []*NavMenuOutput `json:"list"`
}

// Create
type NavMenuCreateReq struct {
	g.Meta      `path:"/nav-menus" method:"post" tags:"NavMenu" summary:"Create a menu"`
	Name        string `v:"required|length:1,100" json:"name"`
	Location    string `v:"max-length:50"          json:"location"`
	Description string `v:"max-length:500"         json:"description"`
}
type NavMenuCreateRes struct {
	*NavMenuOutput
}

// Get one
type NavMenuGetOneReq struct {
	g.Meta `path:"/nav-menus/{id}" method:"get" tags:"NavMenu" summary:"Get menu with items"`
	Id     int64 `v:"required|min:1" dc:"menu id"`
}
type NavMenuGetOneRes struct {
	*NavMenuOutput
}

// Update (meta + replace items)
type NavMenuUpdateReq struct {
	g.Meta      `path:"/nav-menus/{id}" method:"put" tags:"NavMenu" summary:"Update menu and replace all items"`
	Id          int64               `v:"required|min:1"  dc:"menu id"`
	Name        *string             `v:"max-length:100"  json:"name"`
	Location    *string             `v:"max-length:50"   json:"location"`
	Description *string             `v:"max-length:500"  json:"description"`
	Items       []*NavMenuItemInput `json:"items"`
}
type NavMenuUpdateRes struct {
	*NavMenuOutput
}

// Delete
type NavMenuDeleteReq struct {
	g.Meta `path:"/nav-menus/{id}" method:"delete" tags:"NavMenu" summary:"Delete a menu"`
	Id     int64 `v:"required|min:1" dc:"menu id"`
}
type NavMenuDeleteRes struct{}

// Get by location (public GET — allowed by AdminWriteRequired)
type NavMenuByLocationReq struct {
	g.Meta   `path:"/nav-menus/location" method:"get" tags:"NavMenu" summary:"Get menu by location"`
	Location string `v:"required" p:"location" dc:"header|footer|sidebar"`
}
type NavMenuByLocationRes struct {
	*NavMenuOutput
}

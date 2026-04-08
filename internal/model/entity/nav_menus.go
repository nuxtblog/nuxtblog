package entity

import "github.com/gogf/gf/v2/os/gtime"

type NavMenus struct {
	Id          int         `json:"id"          orm:"id"`
	Name        string      `json:"name"        orm:"name"`
	Location    string      `json:"location"    orm:"location"`
	Description string      `json:"description" orm:"description"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"`
}

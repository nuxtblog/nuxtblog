package v1

import "github.com/gogf/gf/v2/frame/g"

type AdminStatsReq struct {
	g.Meta `path:"/admin/stats" method:"get" tags:"Admin" summary:"Get dashboard stats"`
}

type AdminStatsPostStats struct {
	Total     int `json:"total"`
	Published int `json:"published"`
	Draft     int `json:"draft"`
}

type AdminStatsCommentStats struct {
	Total   int `json:"total"`
	Pending int `json:"pending"`
}

type AdminStatsUserStats struct {
	Total  int `json:"total"`
	Active int `json:"active"`
}

type AdminStatsRes struct {
	Posts    AdminStatsPostStats    `json:"posts"`
	Comments AdminStatsCommentStats `json:"comments"`
	Users    AdminStatsUserStats    `json:"users"`
	Views    int64                  `json:"views"`
}

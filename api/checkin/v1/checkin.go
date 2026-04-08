package v1

import "github.com/gogf/gf/v2/frame/g"

// ----------------------------------------------------------------
//  Do check-in
// ----------------------------------------------------------------

type DoCheckinReq struct {
	g.Meta `path:"/checkin" method:"post" tags:"Checkin" summary:"Daily check-in"`
}
type DoCheckinRes struct {
	AlreadyCheckedIn bool `json:"already_checked_in"`
	Streak           int  `json:"streak"`
}

// ----------------------------------------------------------------
//  Get today's check-in status
// ----------------------------------------------------------------

type GetCheckinStatusReq struct {
	g.Meta `path:"/checkin/status" method:"get" tags:"Checkin" summary:"Get today's check-in status and streak"`
}
type GetCheckinStatusRes struct {
	CheckedInToday bool `json:"checked_in_today"`
	Streak         int  `json:"streak"`
}

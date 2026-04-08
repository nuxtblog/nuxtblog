// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import "github.com/gogf/gf/v2/os/gtime"

// Reports is the golang structure for table reports.
type Reports struct {
	Id         int64       `json:"id"         orm:"id"          description:""`
	ReporterId int64       `json:"reporterId" orm:"reporter_id" description:""`
	TargetType string      `json:"targetType" orm:"target_type" description:""`
	TargetId   int64       `json:"targetId"   orm:"target_id"   description:""`
	Reason     string      `json:"reason"     orm:"reason"      description:""`
	Detail     string      `json:"detail"     orm:"detail"      description:""`
	Status     string      `json:"status"     orm:"status"      description:""`
	Notes      string      `json:"notes"      orm:"notes"       description:""`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`
	ResolvedAt *gtime.Time `json:"resolvedAt" orm:"resolved_at" description:""`
}

package v1

import "github.com/gogf/gf/v2/frame/g"

type ReportCreateReq struct {
	g.Meta     `path:"/reports" method:"post" tags:"Report" summary:"Submit a report" auth:"true"`
	TargetType string `v:"required|in:post,comment,user" dc:"target type"`
	TargetId   int64  `v:"required|min:1"                dc:"target id"`
	Reason     string `v:"required|length:1,100"         dc:"reason"`
	Detail     string `v:"max-length:500"                dc:"optional detail"`
}
type ReportCreateRes struct {
	Id int64 `json:"id"`
}

type ReportItem struct {
	Id           int64   `json:"id"`
	ReporterId   int64   `json:"reporter_id"`
	ReporterName string  `json:"reporter_name"`
	TargetType   string  `json:"target_type"`
	TargetId     int64   `json:"target_id"`
	TargetName   string  `json:"target_name"`
	Reason       string  `json:"reason"`
	Detail       string  `json:"detail"`
	Status       string  `json:"status"`
	Notes        string  `json:"notes"`
	CreatedAt    string  `json:"created_at"`
	ResolvedAt   *string `json:"resolved_at,omitempty"`
}

type ReportListReq struct {
	g.Meta `path:"/reports" method:"get" tags:"Report" summary:"Admin: list reports" auth:"true"`
	Status string `d:"pending" dc:"pending|resolved|dismissed|all"`
	Page   int    `d:"1"  v:"min:1"`
	Size   int    `d:"20" v:"min:1|max:100"`
}
type ReportListRes struct {
	Items []ReportItem `json:"items"`
	Total int          `json:"total"`
}

type ReportHandleReq struct {
	g.Meta `path:"/reports/{id}" method:"put" tags:"Report" summary:"Admin: handle report" auth:"true"`
	Id     int64  `v:"required|min:1" dc:"report id"`
	Status string `v:"required|in:resolved,dismissed" dc:"new status"`
	Notes  string `v:"max-length:500"                  dc:"admin notes"`
}
type ReportHandleRes struct{}

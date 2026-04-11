package v1

import "github.com/gogf/gf/v2/frame/g"

// ----------------------------------------------------------------
//  Memory info
// ----------------------------------------------------------------

type MemoryInfo struct {
	Alloc      uint64 `json:"alloc"`
	TotalAlloc uint64 `json:"total_alloc"`
	Sys        uint64 `json:"sys"`
	NumGC      uint32 `json:"num_gc"`
	AllocStr   string `json:"alloc_str"`
	TotalStr   string `json:"total_alloc_str"`
	SysStr     string `json:"sys_str"`
}

// ----------------------------------------------------------------
//  System info
// ----------------------------------------------------------------

type InfoReq struct {
	g.Meta `path:"/system/info" method:"get" tags:"System" summary:"Get system runtime info"`
}

type InfoRes struct {
	OS          string     `json:"os"`
	Arch        string     `json:"arch"`
	GoVersion   string     `json:"go_version"`
	NumCPU      int        `json:"num_cpu"`
	Goroutines  int        `json:"goroutines"`
	Memory      MemoryInfo `json:"memory"`
	Uptime      string     `json:"uptime"`
	UptimeSec   int64      `json:"uptime_sec"`
	DBSizeBytes int64      `json:"db_size_bytes"`
	DBSize      string     `json:"db_size"`
}

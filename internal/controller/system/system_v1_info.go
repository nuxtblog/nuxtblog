package system

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	v1 "github.com/nuxtblog/nuxtblog/api/system/v1"
)

var startTime = time.Now()

func (c *ControllerV1) Info(ctx context.Context, req *v1.InfoReq) (res *v1.InfoRes, err error) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	uptime := time.Since(startTime)

	res = &v1.InfoRes{
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		GoVersion:  runtime.Version(),
		NumCPU:     runtime.NumCPU(),
		Goroutines: runtime.NumGoroutine(),
		Memory: v1.MemoryInfo{
			Alloc:      mem.Alloc,
			TotalAlloc: mem.TotalAlloc,
			Sys:        mem.Sys,
			NumGC:      mem.NumGC,
			AllocStr:   formatBytes(mem.Alloc),
			TotalStr:   formatBytes(mem.TotalAlloc),
			SysStr:     formatBytes(mem.Sys),
		},
		Uptime:    formatDuration(uptime),
		UptimeSec: int64(uptime.Seconds()),
	}

	// Try to get SQLite DB file size
	dbLink, _ := g.Cfg().Get(ctx, "database.default.link")
	if link := dbLink.String(); link != "" {
		if dbPath := parseSQLitePath(link); dbPath != "" {
			if info, e := os.Stat(dbPath); e == nil {
				res.DBSizeBytes = info.Size()
				res.DBSize = formatBytes(uint64(info.Size()))
			}
		}
	}

	return
}

// parseSQLitePath extracts the file path from a GoFrame SQLite link string.
// Format: "sqlite::@file(./blog.db)" or "sqlite:user:pass@file(/path/to/db)"
func parseSQLitePath(link string) string {
	if !strings.HasPrefix(link, "sqlite") {
		return ""
	}
	start := strings.Index(link, "@file(")
	if start < 0 {
		return ""
	}
	start += len("@file(")
	end := strings.Index(link[start:], ")")
	if end < 0 {
		return ""
	}
	return link[start : start+end]
}

func formatBytes(b uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case b >= GB:
		return fmt.Sprintf("%.2f GB", float64(b)/float64(GB))
	case b >= MB:
		return fmt.Sprintf("%.2f MB", float64(b)/float64(MB))
	case b >= KB:
		return fmt.Sprintf("%.2f KB", float64(b)/float64(KB))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

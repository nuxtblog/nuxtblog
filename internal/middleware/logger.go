package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorGray   = "\033[90m"
)

// AccessLogger is a colorful HTTP access log middleware.
func AccessLogger(r *ghttp.Request) {
	start := time.Now()
	r.Middleware.Next()

	status := r.Response.Status
	method := r.Method
	path := r.URL.Path
	elapsed := time.Since(start)
	ip := r.GetClientIp()

	// Skip static/internal paths
	if strings.HasPrefix(path, "/uploads/") ||
		strings.HasPrefix(path, "/swagger") ||
		path == "/api.json" ||
		path == "/favicon.ico" {
		return
	}

	statusColor := colorStatusCode(status)
	methodColor := colorMethod(method)
	elapsedColor := colorElapsed(elapsed)

	fmt.Printf(
		"%s%s%s  %s%s%s  %s%-7s%s  %-45s  %s%s%s  %s%s%s\n",
		colorGray, time.Now().Format("01-02 15:04:05"), colorReset,
		statusColor, fmt.Sprintf("%d", status), colorReset,
		methodColor, method, colorReset,
		path,
		elapsedColor, formatElapsed(elapsed), colorReset,
		colorGray, ip, colorReset,
	)
}

func colorStatusCode(status int) string {
	switch {
	case status >= 500:
		return colorRed + colorBold
	case status >= 400:
		return colorYellow
	case status >= 300:
		return colorCyan
	case status >= 200:
		return colorGreen
	default:
		return colorWhite
	}
}

func colorMethod(method string) string {
	switch method {
	case "GET":
		return colorCyan
	case "POST":
		return colorBlue
	case "PUT", "PATCH":
		return colorYellow
	case "DELETE":
		return colorRed
	case "OPTIONS":
		return colorGray
	default:
		return colorWhite
	}
}

func colorElapsed(d time.Duration) string {
	switch {
	case d > 1*time.Second:
		return colorRed
	case d > 200*time.Millisecond:
		return colorYellow
	default:
		return colorGray
	}
}

func formatElapsed(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

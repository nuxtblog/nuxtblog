package middleware

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/glog"
)

// InitLogger registers a custom glog handler that outputs concise, colored log
// lines matching the AccessLogger style: HH:MM:SS  LEVEL  message
//
// Stdout is handled directly by this handler (with ANSI colors).
// File logging is handled by glog's built-in pipeline via in.Next —
// we disable glog's own stdout so it only writes to file.
func InitLogger() {
	// Disable glog's built-in stdout; we handle it ourselves with colors.
	glog.SetStdoutPrint(false)

	glog.SetDefaultHandler(func(ctx context.Context, in *glog.HandlerInput) {
		// Ensure the logger instance doesn't double-print to stdout;
		// we handle stdout ourselves below.
		in.Logger.SetStdoutPrint(false)

		ts := in.Time.Format("01-02 15:04:05")
		msg := in.Content + in.ValuesContent()

		// Colored line to stdout
		levelColor := logLevelColor(in.Level)
		fmt.Printf(
			"%s%s%s  %s%-4s%s   %s\n",
			colorGray, ts, colorReset,
			levelColor, in.LevelFormat, colorReset,
			msg,
		)

		// Plain format into buffer for file logging
		in.Buffer.Reset()
		in.Buffer.WriteString(fmt.Sprintf("%s  %-4s   %s\n", ts, in.LevelFormat, msg))
		in.Next(ctx)
	})
}

func logLevelColor(level int) string {
	switch level {
	case glog.LEVEL_DEBU:
		return colorGray
	case glog.LEVEL_INFO:
		return colorGreen
	case glog.LEVEL_NOTI:
		return colorCyan
	case glog.LEVEL_WARN:
		return colorYellow
	case glog.LEVEL_ERRO, glog.LEVEL_CRIT, glog.LEVEL_PANI, glog.LEVEL_FATA:
		return colorRed + colorBold
	default:
		return colorWhite
	}
}

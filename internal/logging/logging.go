package logging

import (
	"io"
	"os"
	"strings"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/klog"
)

func Init(level string, out io.Writer) {
	if out == nil {
		out = os.Stdout
	}

	hlog.SetOutput(out)
	klog.SetOutput(out)

	hlog.SetLevel(parseHertzLevel(level))
	klog.SetLevel(parseKitexLevel(level))
}

func parseHertzLevel(level string) hlog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "trace":
		return hlog.LevelTrace
	case "debug":
		return hlog.LevelDebug
	case "notice":
		return hlog.LevelNotice
	case "warn", "warning":
		return hlog.LevelWarn
	case "error":
		return hlog.LevelError
	case "fatal":
		return hlog.LevelFatal
	default:
		return hlog.LevelInfo
	}
}

func parseKitexLevel(level string) klog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "trace":
		return klog.LevelTrace
	case "debug":
		return klog.LevelDebug
	case "notice":
		return klog.LevelNotice
	case "warn", "warning":
		return klog.LevelWarn
	case "error":
		return klog.LevelError
	case "fatal":
		return klog.LevelFatal
	default:
		return klog.LevelInfo
	}
}

package logger

import (
	"github.com/eldius/file-server/internal/config"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

var logger *slog.Logger

func GetLogger() *slog.Logger {
	if logger == nil {
		setupLogger()
	}
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return logger.With("caller", details.Name())
	}
	return logger
}

func setupLogger() {
	l := slog.LevelError
	if config.GetDebugModeEnabled() {
		l = slog.LevelDebug
	}
	var handler slog.Handler
	if strings.ToLower(config.GetLogFormat()) == "json" {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: l})
	} else {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: l})
	}
	logger = slog.New(handler)
	slog.SetDefault(logger)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.LstdFlags)
}

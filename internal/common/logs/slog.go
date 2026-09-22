package logs

import (
	"log/slog"
	"os"
	"strconv"
)

// Init sets the process-wide default slog logger: JSON output in
// production, human-readable text when LOCAL_ENV is set.
func Init() {
	slog.SetDefault(slog.New(newHandler()))
}

func newHandler() slog.Handler {
	if isLocalEnv, _ := strconv.ParseBool(os.Getenv("LOCAL_ENV")); isLocalEnv {
		return slog.NewTextHandler(os.Stdout, nil)
	}

	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: renameFields,
	})
}

func renameFields(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.LevelKey:
		a.Key = "severity"
	case slog.MessageKey:
		a.Key = "message"
	}
	return a
}

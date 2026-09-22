package logs

import (
	"context"
	"log/slog"
)

func LogCommandExecution(ctx context.Context, commandName string, cmd interface{}, err error) {
	if err == nil {
		slog.InfoContext(ctx, commandName + " command succeeded", "cmd", cmd)
		return
	}

	slog.ErrorContext(ctx, commandName + " command failed", "cmd", cmd, "error", err)
}

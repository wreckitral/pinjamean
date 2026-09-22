package decorator

import (
	"context"
	"fmt"
	"log/slog"
)

type commandLoggingDecorator[C any] struct {
	base CommandHandler[C]
	logger *slog.Logger
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	logger := d.logger.With(
		"command", generateActionName(cmd),
		"command_body", fmt.Sprintf("%#v", cmd),
	)

	logger.DebugContext(ctx, "Executing command")
	defer func() {
		if err == nil {
			logger.InfoContext(ctx, "Command executed successfully")
		} else {
			logger.ErrorContext(ctx, "Failed to execute command", "error", err)
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type queryLoggingDecorator[C any, R any] struct {
	base QueryHandler[C, R]
	logger *slog.Logger
}

func (d queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	logger := d.logger.With(
		"query", generateActionName(cmd),
		"query_body", fmt.Sprintf("%#v", cmd),
	)

	logger.DebugContext(ctx, "Executing query")
	defer func() {
		if err == nil {
			logger.InfoContext(ctx, "Query executed successfully")
		} else {
			logger.ErrorContext(ctx, "Failed to execute query", "error", err)
		}
	}()

	return d.base.Handle(ctx, cmd)
}


func generateActionName(cmd interface{}) string {
	return fmt.Sprintf("%T", cmd)
}

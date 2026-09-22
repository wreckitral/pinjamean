package decorator

import (
	"context"
	"log/slog"
)

func ApplyCommandDecorators[H any](
	handler CommandHandler[H],
	logger *slog.Logger,
) CommandHandler[H] {
	return commandLoggingDecorator[H]{
		base: handler,
		logger: logger,
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

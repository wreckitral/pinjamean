package decorator

import (
	"context"
	"log/slog"
)

func ApplyQueryDecorators[H any, R any](
	handler QueryHandler[H, R],
	logger *slog.Logger,
) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		base: handler,
		logger: logger,
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}

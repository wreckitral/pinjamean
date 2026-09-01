package application

import (
	"context"
	"errors"
)

var ErrorNotFound = errors.New("application not found")

type Repository interface {
	Add(ctx context.Context, app *Application) error
	Get(ctx context.Context, uuid string) (*Application, error)
	Update(ctx context.Context, uuid string, updateFn func(app *Application) (*Application, error)) error
}

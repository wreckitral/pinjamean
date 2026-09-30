package query

import (
	"context"
	"log/slog"

	"github.com/wreckitral/pinjamean/internal/common/decorator"
	"github.com/wreckitral/pinjamean/internal/identity/domain/account"
)

type GetAccountByEmail struct {
	Email string
}

type AccountView struct {
	UUID  string
	Email string
	Role  account.Role
}

type getAccountByEmailHandler struct {
	repo account.Repository
}

type GetAccountByEmailHandler decorator.QueryHandler[GetAccountByEmail, AccountView]

func NewGetAccountByEmailHandler(repo account.Repository, logger *slog.Logger) GetAccountByEmailHandler {
	if repo == nil {
		panic("nil repo")
	}
	return decorator.ApplyQueryDecorators[GetAccountByEmail, AccountView](
		getAccountByEmailHandler{repo: repo}, logger,
	)
}

func (h getAccountByEmailHandler) Handle(ctx context.Context, query GetAccountByEmail) (AccountView, error) {
	a, err := h.repo.GetAccountByEmail(ctx, query.Email)
	if err != nil {
		return AccountView{}, err
	}
	return AccountView{UUID: a.UUID(), Email: a.Email(), Role: a.Role()}, nil
}

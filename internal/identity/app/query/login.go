package query

import (
	"context"
	"log/slog"

	"github.com/wreckitral/pinjamean/internal/common/auth"
	"github.com/wreckitral/pinjamean/internal/common/decorator"
	commonerrors "github.com/wreckitral/pinjamean/internal/common/errors"
	"github.com/wreckitral/pinjamean/internal/identity/domain/account"
)

type Login struct {
	Email string
	Password string
}

type LoginView struct {
	Token string
}

type loginHandler struct {
	repo account.Repository
	secret []byte
}

type LoginHandler decorator.QueryHandler[Login, LoginView]

func NewLoginHandler(repo account.Repository, secret []byte, logger *slog.Logger) LoginHandler {
	if repo == nil { panic("nil repo") }
	return decorator.ApplyQueryDecorators[Login, LoginView](loginHandler{repo: repo, secret: secret}, logger)
}

func (h loginHandler) Handle(ctx context.Context, query Login) (LoginView, error) {
	a, err := h.repo.GetAccountByEmail(ctx, query.Email)
	if err != nil { return LoginView{}, err }
	if a == nil || !a.VerifyPassword(query.Password) {
		return LoginView{}, commonerrors.NewIncorrectInputError("invalid email or password", "invalid-credentials")
	}
	token, err := auth.IssueJWT(a.UUID(), string(a.Role()), h.secret)
	if err != nil { return LoginView{}, err }
	return LoginView{Token: token}, nil
}

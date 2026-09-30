package command

import (
	"context"
	"errors"
	"log/slog"

	"github.com/wreckitral/pinjamean/internal/common/decorator"
	commonerrors "github.com/wreckitral/pinjamean/internal/common/errors"
	"github.com/wreckitral/pinjamean/internal/identity/domain/account"
	"golang.org/x/crypto/bcrypt"
)

type SaveAccount struct {
	UUID     string
	Email    string
	Password string
	Role     account.Role
}

type saveAccountHandler struct {
	repo account.Repository
}

type SaveAccountHandler decorator.CommandHandler[SaveAccount]

func NewSaveAccountHandler(repo account.Repository, logger *slog.Logger) SaveAccountHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[SaveAccount](
		saveAccountHandler{repo: repo},
		logger,
	)
}

func (h saveAccountHandler) Handle(ctx context.Context, cmd SaveAccount) (err error) {
	if cmd.Password == "" {
		return errors.New("password is required")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a, err := account.NewAccount(cmd.UUID, cmd.Email, string(passwordHash), cmd.Role)
	if err != nil {
		return err
	}

	existingAccount, err := h.repo.GetAccountByEmail(ctx, cmd.Email)
	if err != nil {
		return err
	}
	if existingAccount != nil {
		return commonerrors.NewIncorrectInputError("An account with this email already exists", "account-already-exists")
	}

	if err := h.repo.SaveAccount(ctx, a); err != nil {
		return err
	}

	return nil
}

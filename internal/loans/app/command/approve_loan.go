package command

import (
	"context"
	"log/slog"

	"github.com/wreckitral/pinjamean/internal/common/decorator"
	"github.com/wreckitral/pinjamean/internal/loans/domain/loan"
)

type ApproveLoan struct {
	LoanUUID string

	ApproverUUID string
}

type approveLoanHandler struct {
	repo loan.Repository
}

type ApproveLoanHandler decorator.CommandHandler[ApproveLoan]

func NewApproveLoanHandler(repo loan.Repository, logger *slog.Logger) ApproveLoanHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[ApproveLoan](
		approveLoanHandler{repo: repo},
		logger,
	)
}

func (h approveLoanHandler) Handle(ctx context.Context, cmd ApproveLoan) (err error) {
	l, err := h.repo.GetLoanByID(ctx, cmd.LoanUUID)
	if err != nil {
		return err
	}

	if err := l.Approve(cmd.ApproverUUID); err != nil {
		return err
	}

	return h.repo.SaveLoan(ctx, l)
}

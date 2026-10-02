package command

import (
	"context"
	"log/slog"

	"github.com/wreckitral/pinjamean/internal/common/decorator"
	"github.com/wreckitral/pinjamean/internal/loans/domain/loan"
)

type ReassignLoan struct {
	LoanUUID string

	ReassignToUUID string
}

type reassignLoanHandler struct {
	repo loan.Repository
}

type ReassignLoanHandler decorator.CommandHandler[ReassignLoan]

func NewReassignLoanHandler(repo loan.Repository, logger *slog.Logger) ReassignLoanHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[ReassignLoan](
		reassignLoanHandler{repo: repo},
		logger,
	)
}

func (h reassignLoanHandler) Handle(ctx context.Context, cmd ReassignLoan) (err error) {
	l, err := h.repo.GetLoanByID(ctx, cmd.LoanUUID)
	if err != nil {
		return err
	}

	if err := l.Reassign(cmd.ReassignToUUID); err != nil {
		return err
	}

	return h.repo.SaveLoan(ctx, l)
}

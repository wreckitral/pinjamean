package command

import (
	"context"
	"log/slog"

	"github.com/wreckitral/pinjamean/internal/common/decorator"
	"github.com/wreckitral/pinjamean/internal/loans/domain/loan"
)

type RejectLoan struct {
	LoanUUID    string
	AnalystUUID string
	Reason      string
}

type rejectLoanHandler struct {
	repo loan.Repository
}

type RejectLoanHandler decorator.CommandHandler[RejectLoan]

func NewRejectLoanHandler(repo loan.Repository, logger *slog.Logger) RejectLoanHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[RejectLoan](rejectLoanHandler{repo: repo}, logger)
}

func (h rejectLoanHandler) Handle(ctx context.Context, cmd RejectLoan) error {
	l, err := h.repo.GetLoanByID(ctx, cmd.LoanUUID)
	if err != nil {
		return err
	}
	if err := l.Reject(cmd.AnalystUUID, cmd.Reason); err != nil {
		return err
	}
	return h.repo.SaveLoan(ctx, l)
}

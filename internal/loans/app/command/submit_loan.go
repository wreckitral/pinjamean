package command

import (
	"context"
	"log/slog"

	"github.com/wreckitral/pinjamean/internal/common/decorator"
	"github.com/wreckitral/pinjamean/internal/loans/domain/loan"
)

type SubmitLoan struct {
	LoanUUID string

	BorrowerUUID string
	LoanAmountIDR int64
	TermMonths int
	LoanType loan.LoanType
}

type submitLoanHandler struct {
	repo loan.Repository
}

type SubmitLoanHandler decorator.CommandHandler[SubmitLoan]

func NewSubmitLoanHandler(repo loan.Repository, logger *slog.Logger) SubmitLoanHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[SubmitLoan](
		submitLoanHandler{repo: repo},
		logger,
	)
}

func (h submitLoanHandler) Handle(ctx context.Context, cmd SubmitLoan) (err error) {
	l, err := loan.NewLoan(cmd.LoanUUID, cmd.BorrowerUUID, cmd.LoanAmountIDR, cmd.TermMonths, cmd.LoanType)
	if err != nil {
		return err
	}

	if err := h.repo.SaveLoan(ctx, l); err != nil {
		return err
	}

	return nil
}

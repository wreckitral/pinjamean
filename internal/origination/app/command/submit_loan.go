package command

import (
	"context"

	"github.com/wreckitral/pinjamean/internal/origination/domain/application"
)

type SubmitLoan struct {
	UUID string
	BorrowerUUID string
	LoanAmountIDR int64
	TermMonths int
	LoanType string
}

type SubmitLoanHandler interface {
	Handle(ctx context.Context, cmd SubmitLoan) error
}

type submitLoanHandler struct {
	appRepo application.Repository
}

func NewSubmitLoanHandler(appRepo application.Repository) SubmitLoanHandler {
	if appRepo == nil {
		panic("nil appRepo")
	}

	return submitLoanHandler{appRepo: appRepo}
}

func (h submitLoanHandler) Handle(ctx context.Context, cmd SubmitLoan) error {
	app, err := application.NewApplication(
		cmd.UUID,
		cmd.BorrowerUUID,
		cmd.LoanAmountIDR,
		cmd.TermMonths,
		application.LoanType(cmd.LoanType),
	)
	if err != nil {
		return err
	}

	if err := h.appRepo.Add(ctx, app); err != nil {
		return err
	}

	return nil
}

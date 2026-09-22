package query

import (
	"context"
	"log/slog"

	"github.com/wreckitral/pinjamean/internal/common/decorator"
	"github.com/wreckitral/pinjamean/internal/loans/domain/loan"
)

type GetLoan struct {
	LoanUUID string
}

type LoanView struct {
	UUID string

	BorrowerUUID string

	LoanAmountIDR   int64
	TermMonths      int
	LoanType        loan.LoanType
	InterestRateAPR float64
	Status          loan.LoanStatus
}

type GetLoanByIDHandler decorator.QueryHandler[GetLoan, LoanView]

type getLoanByIDHandler struct {
	repo loan.Repository
}

func NewGetLoanByIDHandler(repo loan.Repository, logger *slog.Logger) GetLoanByIDHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyQueryDecorators[GetLoan, LoanView](
		getLoanByIDHandler{repo: repo},
		logger,
	)
}

func (h getLoanByIDHandler) Handle(ctx context.Context, query GetLoan) (LoanView, error) {
	l, err := h.repo.GetLoanByID(ctx, query.LoanUUID)
	if err != nil {
		return LoanView{}, err
	}

	return LoanView{
		UUID:            l.UUID(),
		BorrowerUUID:    l.BorrowerUUID(),
		LoanAmountIDR:   l.AmountIDR(),
		TermMonths:      l.TermMonths(),
		LoanType:        l.Type(),
		InterestRateAPR: l.InterestRateAPR(),
		Status:          l.Status(),
	}, nil
}

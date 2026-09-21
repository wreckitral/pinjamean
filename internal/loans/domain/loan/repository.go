package loan

import "context"

type Repository interface {
	SaveLoan(ctx context.Context, loan *Loan) error
	GetLoanByID(ctx context.Context, LoanUUID string) (*Loan, error)
}

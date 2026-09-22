package adapters

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wreckitral/pinjamean/internal/loans/adapters/sqlcgen"
	"github.com/wreckitral/pinjamean/internal/loans/domain/loan"
)

type postgresLoan struct {
    UUID            string    `db:"uuid"`
    BorrowerUUID    string    `db:"borrower_uuid"`
    AmountIDR       int64     `db:"amount_idr"`
    TermMonths      int       `db:"term_months"`
    LoanType        string    `db:"loan_type"`
    InterestRateAPR float64   `db:"interest_rate_apr"`
    Status          string    `db:"status"`
    CreatedAt       time.Time `db:"created_at"`
    UpdatedAt       time.Time `db:"updated_at"`
}

type PostgresLoanRepository struct {
	q *sqlcgen.Queries
}

func NewPostgresLoanRepository(db *pgxpool.Pool) *PostgresLoanRepository {
	return &PostgresLoanRepository{
		q: sqlcgen.New(db),
	}
}

func (r *PostgresLoanRepository) SaveLoan(ctx context.Context, l *loan.Loan) error {
	loanUuid, err := uuid.Parse(l.UUID())
	if err != nil {
		return fmt.Errorf("invalid loan UUID %q: %w", l.UUID(), err)
	}

	borrowerUuid, err := uuid.Parse(l.BorrowerUUID())
	if err != nil {
		return fmt.Errorf("invalid borrower UUID %q: %w", l.BorrowerUUID(), err)
	}

	return r.q.SaveLoan(ctx, sqlcgen.SaveLoanParams{
		Uuid:            loanUuid,
		BorrowerUuid:    borrowerUuid,
		AmountIdr:       l.AmountIDR(),
		TermMonths:      int32(l.TermMonths()),
		LoanType:        string(l.Type()),
		InterestRateApr: l.InterestRateAPR(),
		Status:          string(l.Status()),
		CreatedAt:       l.CreatedAt(),
		UpdatedAt:       l.UpdatedAt(),
	})
}

func (r *PostgresLoanRepository) GetLoanByID(ctx context.Context, loanUUID string) (*loan.Loan, error) {
	loanUuid, err := uuid.Parse(loanUUID)
	if err != nil {
		return nil, fmt.Errorf("invalid loan UUID %q: %w", loanUUID, err)
	}

	row, err := r.q.GetLoanByID(ctx, loanUuid)
	if err != nil {
		return nil, err
	}

	return loan.UnmarshalLoanFromDatabase(
		row.Uuid.String(),
		row.BorrowerUuid.String(),
		row.AmountIdr,
		int(row.TermMonths),
		loan.LoanType(row.LoanType),
		row.InterestRateApr,
		loan.LoanStatus(row.Status),
		row.CreatedAt,
		row.UpdatedAt,
	)
}

func NewPostgresSQLConnection() (*pgxpool.Pool, error) {
	ctx := context.Background()
	return pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
}

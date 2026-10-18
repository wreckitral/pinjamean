package adapters

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wreckitral/pinjamean/internal/loans/adapters/sqlcgen"
	"github.com/wreckitral/pinjamean/internal/loans/domain/loan"
)

type postgresLoan struct {
	UUID            string    `db:"uuid"`
	BorrowerUUID    string    `db:"borrower_uuid"`
	OfficerUUID     string    `db:"officer_uuid"`
	AssignedToUUID  string    `db:"assigned_to_uuid"`
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
	officerUuid, err := uuid.Parse(l.OfficerUUID())
	if err != nil {
		return fmt.Errorf("invalid officer UUID %q: %w", l.OfficerUUID(), err)
	}
	assignedToUuid, err := uuid.Parse(l.AssignedToUUID())
	if err != nil {
		return fmt.Errorf("invalid assigned-to UUID %q: %w", l.AssignedToUUID(), err)
	}
	var analystUuid pgtype.UUID
	if l.AnalystUUID() != "" {
		parsed, err := uuid.Parse(l.AnalystUUID())
		if err != nil {
			return fmt.Errorf("invalid analyst UUID %q: %w", l.AnalystUUID(), err)
		}
		analystUuid = pgtype.UUID{Bytes: parsed, Valid: true}
	}

	return r.q.SaveLoan(ctx, sqlcgen.SaveLoanParams{
		Uuid:            loanUuid,
		BorrowerUuid:    borrowerUuid,
		OfficerUuid:     officerUuid,
		AssignedToUuid:  assignedToUuid,
		AmountIdr:       l.AmountIDR(),
		TermMonths:      int32(l.TermMonths()),
		LoanType:        string(l.Type()),
		InterestRateApr: l.InterestRateAPR(),
		Status:          string(l.Status()),
		CreatedAt:       l.CreatedAt(),
		UpdatedAt:       l.UpdatedAt(),
		AnalystUuid:     analystUuid,
		RejectionReason: pgtype.Text{String: l.RejectionReason(), Valid: l.RejectionReason() != ""},
	})
}

func (r *PostgresLoanRepository) GetLoanByID(ctx context.Context, loanUUID string) (*loan.Loan, error) {
	loanUuid, err := uuid.Parse(loanUUID)
	if err != nil {
		return nil, fmt.Errorf("invalid loan UUID %q: %w", loanUUID, err)
	}

	row, err := r.q.GetLoanByID(ctx, loanUuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, loan.ErrLoanNotFound
		}
		return nil, err
	}

	return loan.UnmarshalLoanFromDatabase(
		row.Uuid.String(),
		row.BorrowerUuid.String(),
		row.OfficerUuid.String(),
		row.AssignedToUuid.String(),
		nullableUUIDString(row.AnalystUuid),
		row.AmountIdr,
		int(row.TermMonths),
		loan.LoanType(row.LoanType),
		row.InterestRateApr,
		loan.LoanStatus(row.Status),
		nullableTextString(row.RejectionReason),
		row.CreatedAt,
		row.UpdatedAt,
	)
}

func nullableUUIDString(value pgtype.UUID) string {
	if !value.Valid {
		return ""
	}
	return uuid.UUID(value.Bytes).String()
}

func nullableTextString(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

func NewPostgresSQLConnection() (*pgxpool.Pool, error) {
	ctx := context.Background()
	return pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
}

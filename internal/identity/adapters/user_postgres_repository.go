package adapters

import (
	"context"
	"fmt"
	"time"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wreckitral/pinjamean/internal/identity/adapters/sqlcgen"
	"github.com/wreckitral/pinjamean/internal/identity/domain/account"
)

type postgresUser struct {
	UUID string `db:"uuid"`
	Email string `db:"email"`
	PasswordHashed string `db:"passwordHashed"`
	Role string `db:"role"`
    CreatedAt       time.Time `db:"created_at"`
    UpdatedAt       time.Time `db:"updated_at"`
}

type PostgresUserRepository struct {
	q *sqlcgen.Queries
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		q: sqlcgen.New(db),
	}
}

func (r *PostgresUserRepository) SaveAccount(ctx context.Context, a *account.Account) error {
	userUUID, err := uuid.Parse(a.UUID())
	if err != nil {
		return fmt.Errorf("invalid user UUID %q: %w", a.UUID(), err)
	}

	return r.q.SaveAccount(ctx, sqlcgen.SaveAccountParams{
		Uuid:           userUUID,
		Email:          a.Email(),
		Passwordhashed: a.PasswordHash(),
		Role:           string(a.Role()),
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
}

func (r *PostgresUserRepository) GetAccountByEmail(ctx context.Context, email string) (*account.Account, error) {
	row, err := r.q.GetAccountByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return nil, nil }
		return nil, err
	}
	return account.NewAccount(row.Uuid.String(), row.Email, row.Passwordhashed, account.Role(row.Role))
}

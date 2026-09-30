package account

import "context"

type Repository interface {
	SaveAccount(ctx context.Context, a *Account) error
	GetAccountByEmail(ctx context.Context, email string) (*Account, error)
}

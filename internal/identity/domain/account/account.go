package account

import (
	"errors"
	"time"
)

type Role string

const (
	RoleOfficer Role = "officer"
	RoleCreditAnalyst Role = "credit_analyst"
	RoleSupervisor Role = "supervisor"
)

type Account struct {
	uuid string
	email string
	passwordHash string
	role Role
	createdAt time.Time
	updatedAt time.Time
}

func NewAccount(uuid, email, passwordHash string, role Role) (*Account, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if passwordHash == "" {
		return nil, errors.New("hashed password is required")
	}

	switch role {
	case RoleOfficer, RoleCreditAnalyst, RoleSupervisor:
		// valid
	default:
		return nil, errors.New("unsupported role")
	}

	return &Account{
		uuid:         uuid,
		email:        email,
		passwordHash: passwordHash,
		role:         role,
	}, nil
}

func (a *Account) UUID() string         { return a.uuid }
func (a *Account) Email() string        { return a.email }
func (a *Account) PasswordHash() string { return a.passwordHash }
func (a *Account) Role() Role           { return a.role }
func (a *Account) CreatedAt() time.Time { return a.createdAt }
func (a *Account) UpdatedAt() time.Time { return a.updatedAt }

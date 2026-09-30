package account

import "golang.org/x/crypto/bcrypt"

func (a *Account) VerifyPassword(candidate string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.passwordHash), []byte(candidate))

	return err == nil
}

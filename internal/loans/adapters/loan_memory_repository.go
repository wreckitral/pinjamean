package adapters

import (
	"context"
	"errors"
	"sync"

	"github.com/wreckitral/pinjamean/internal/loans/domain/loan"
)

type MemoryLoanRepository struct {
    mu    sync.RWMutex
    loans map[string]loan.Loan
}

func NewMemoryLoanRepository() *MemoryLoanRepository {
    return &MemoryLoanRepository{
        loans: make(map[string]loan.Loan),
    }
}

func (r *MemoryLoanRepository) SaveLoan(ctx context.Context, l *loan.Loan) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.loans[l.UUID()] = *l
    return nil
}

func (r *MemoryLoanRepository) GetLoanByID(ctx context.Context, LoanUUID string) (*loan.Loan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if l, exists := r.loans[LoanUUID]; exists {
		return &l, nil
	}

	return nil, errors.New("loan not found on memory")
}

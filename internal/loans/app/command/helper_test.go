package command_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/wreckitral/pinjamean/internal/loans/domain/loan"
)

var discardLogger = slog.New(slog.DiscardHandler)

func seedPendingLoan(t *testing.T, repo loan.Repository) *loan.Loan {
	t.Helper()

	officer := uuid.NewString()

	l, err := loan.NewLoan(uuid.NewString(), uuid.NewString(), officer, officer, 50_000_000, 24, loan.TypeKTA)

	require.NoError(t, err)
	require.NoError(t, repo.SaveLoan(context.Background(), l))

	return l
}

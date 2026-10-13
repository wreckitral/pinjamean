package command_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/wreckitral/pinjamean/internal/loans/adapters"
	"github.com/wreckitral/pinjamean/internal/loans/app/command"
	"github.com/wreckitral/pinjamean/internal/loans/domain/loan"
)


func TestApproveLoan_ApprovePendingLoan(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := adapters.NewMemoryLoanRepository()
	h := command.NewApproveLoanHandler(repo, discardLogger)

	l := seedPendingLoan(t, repo)
	approver := uuid.NewString()

	require.NoError(t, h.Handle(ctx, command.ApproveLoan{
		LoanUUID: l.UUID(),

		ApproverUUID: approver,
	}))

	got, err := repo.GetLoanByID(ctx, l.UUID())
	require.NoError(t, err)
	require.Equal(t, loan.LoanStatusApproved, got.Status())
	require.Equal(t, approver, got.AnalystUUID())
}


func TestApproveLoan_AlreadyDecidedLeavesLoanUnchanged(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := adapters.NewMemoryLoanRepository()
	h := command.NewApproveLoanHandler(repo, discardLogger)

	l := seedPendingLoan(t, repo)
	firstAnalyst := uuid.NewString()
	require.NoError(t, h.Handle(ctx, command.ApproveLoan{LoanUUID: l.UUID(), ApproverUUID: firstAnalyst}))

	err := h.Handle(ctx, command.ApproveLoan{LoanUUID: l.UUID(), ApproverUUID: uuid.NewString()})
	require.ErrorIs(t, err, loan.ErrLoanStatusNotPending)

	got, err := repo.GetLoanByID(ctx, l.UUID())
	require.NoError(t, err)
	require.Equal(t, loan.LoanStatusApproved, got.Status())
	require.Equal(t, firstAnalyst, got.AnalystUUID(), "a failed approval must not overwrite the first decision")
}

func TestApproveLoan_UnknownLoan(t *testing.T) {
	t.Parallel()

	h := command.NewApproveLoanHandler(adapters.NewMemoryLoanRepository(), discardLogger)

	err := h.Handle(context.Background(), command.ApproveLoan{
		LoanUUID:     uuid.NewString(),
		ApproverUUID: uuid.NewString(),
	})
	require.Error(t, err)
}

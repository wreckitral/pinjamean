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

func TestRejectLoan_RejectsPendingLoan(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := adapters.NewMemoryLoanRepository()
	h := command.NewRejectLoanHandler(repo, discardLogger)

	l := seedPendingLoan(t, repo)
	analyst := uuid.NewString()

	require.NoError(t, h.Handle(ctx, command.RejectLoan{
		LoanUUID: l.UUID(),
		AnalystUUID: analyst,
		Reason: "insufficient income",
	}))

	got, err := repo.GetLoanByID(ctx, l.UUID())
	require.NoError(t, err)
	require.Equal(t, loan.LoanStatusRejected, got.Status())
	require.Equal(t, analyst, got.AnalystUUID())
	require.Equal(t, "insufficient income", got.RejectionReason())
}

func TestRejectLoan_AlreadyDecidedLeavesLoanUnchanged(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := adapters.NewMemoryLoanRepository()
	approve := command.NewApproveLoanHandler(repo, discardLogger)
	reject := command.NewRejectLoanHandler(repo, discardLogger)

	l := seedPendingLoan(t, repo)
	firstAnalyst := uuid.NewString()
	require.NoError(t, approve.Handle(ctx, command.ApproveLoan{LoanUUID: l.UUID(), ApproverUUID: firstAnalyst}))

	err := reject.Handle(ctx, command.RejectLoan{
		LoanUUID:     l.UUID(),
		AnalystUUID: uuid.NewString(),
		Reason:       "changed my mind",
	})
	require.ErrorIs(t, err, loan.ErrLoanStatusNotPending)

	got, err := repo.GetLoanByID(ctx, l.UUID())
	require.NoError(t, err)
	require.Equal(t, loan.LoanStatusApproved, got.Status())
	require.Equal(t, firstAnalyst, got.AnalystUUID())
	require.Empty(t, got.RejectionReason())
}


func TestRejectLoan_MissingReasonLeavesLoanPending(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := adapters.NewMemoryLoanRepository()
	h := command.NewRejectLoanHandler(repo, discardLogger)

	l := seedPendingLoan(t, repo)

	err := h.Handle(ctx, command.RejectLoan{
		LoanUUID:     l.UUID(),
		AnalystUUID: uuid.NewString(),
		Reason:       "",
	})
	require.ErrorIs(t, err, loan.ErrMissingRejectionReason)

	got, err := repo.GetLoanByID(ctx, l.UUID())
	require.NoError(t, err)
	require.Equal(t, loan.LoanStatusPending, got.Status())
}

func TestRejectLoan_UnknownLoan(t *testing.T) {
	t.Parallel()

	h := command.NewRejectLoanHandler(adapters.NewMemoryLoanRepository(), discardLogger)

	err := h.Handle(context.Background(), command.RejectLoan{
		LoanUUID:     uuid.NewString(),
		AnalystUUID: uuid.NewString(),
		Reason:       "whatever",
	})
	require.Error(t, err)
}

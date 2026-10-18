package adapters_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/wreckitral/pinjamean/internal/loans/adapters"
	"github.com/wreckitral/pinjamean/internal/loans/domain/loan"
)

func testLoanRepositoryContract(t *testing.T, repo loan.Repository) {
	t.Run("SaveAndGetPendingLoan", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		loanUUID := uuid.NewString()
		borrowerUUID := uuid.NewString()
		officerUUID := uuid.NewString()
		assignedToUUID := uuid.NewString()
		amountIDR := int64(100_000_000)
		termMonths := 24
		loanType := loan.TypeKTA

		expected, err := loan.NewLoan(
			loanUUID,
			borrowerUUID,
			officerUUID,
			assignedToUUID,
			amountIDR,
			termMonths,
			loanType,
		)
		require.NoError(t, err)

		err = repo.SaveLoan(ctx, expected)
		require.NoError(t, err)

		got, err := repo.GetLoanByID(ctx, loanUUID)
		require.NoError(t, err)
		require.NotNil(t, got)

		require.Equal(t, expected.UUID(), got.UUID())
		require.Equal(t, expected.BorrowerUUID(), got.BorrowerUUID())
		require.Equal(t, expected.OfficerUUID(), got.OfficerUUID())
		require.Equal(t, expected.AssignedToUUID(), got.AssignedToUUID())
		require.Equal(t, expected.AmountIDR(), got.AmountIDR())
		require.Equal(t, expected.TermMonths(), got.TermMonths())
		require.Equal(t, expected.Type(), got.Type())
		require.Equal(t, expected.InterestRateAPR(), got.InterestRateAPR())
		require.Equal(t, expected.Status(), got.Status())
		require.Equal(t, expected.AnalystUUID(), got.AnalystUUID())
		require.Equal(t, expected.RejectionReason(), got.RejectionReason())
		require.WithinDuration(t, expected.CreatedAt(), got.CreatedAt(), time.Microsecond)
		require.WithinDuration(t, expected.UpdatedAt(), got.UpdatedAt(), time.Microsecond)
	})

	t.Run("SaveApproveAndSaveAgain", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		loanUUID := uuid.NewString()
		borrowerUUID := uuid.NewString()
		officerUUID := uuid.NewString()
		assignedToUUID := uuid.NewString()
		amountIDR := int64(75_000_000)
		termMonths := 12
		loanType := loan.TypeKKB

		l, err := loan.NewLoan(
			loanUUID,
			borrowerUUID,
			officerUUID,
			assignedToUUID,
			amountIDR,
			termMonths,
			loanType,
		)
		require.NoError(t, err)

		err = repo.SaveLoan(ctx, l)
		require.NoError(t, err)

		analystUUID := uuid.NewString()
		err = l.Approve(analystUUID)
		require.NoError(t, err)

		err = repo.SaveLoan(ctx, l)
		require.NoError(t, err)

		got, err := repo.GetLoanByID(ctx, loanUUID)
		require.NoError(t, err)
		require.NotNil(t, got)

		require.Equal(t, loan.LoanStatusApproved, got.Status())
		require.Equal(t, analystUUID, got.AnalystUUID())
		require.Empty(t, got.RejectionReason())

		require.Equal(t, l.UUID(), got.UUID())
		require.Equal(t, l.BorrowerUUID(), got.BorrowerUUID())
		require.Equal(t, l.OfficerUUID(), got.OfficerUUID())
		require.Equal(t, l.AssignedToUUID(), got.AssignedToUUID())
		require.Equal(t, l.AmountIDR(), got.AmountIDR())
		require.Equal(t, l.TermMonths(), got.TermMonths())
		require.Equal(t, l.Type(), got.Type())
		require.Equal(t, l.InterestRateAPR(), got.InterestRateAPR())
		require.WithinDuration(t, l.CreatedAt(), got.CreatedAt(), time.Microsecond)
		require.WithinDuration(t, l.UpdatedAt(), got.UpdatedAt(), time.Microsecond)
	})

	t.Run("GetUnknownLoanReturnsErrNotFound", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		unknownUUID := uuid.NewString()
		got, err := repo.GetLoanByID(ctx, unknownUUID)
		require.ErrorIs(t, err, loan.ErrLoanNotFound)
		require.Nil(t, got)
	})
}

func TestMemoryLoanRepository(t *testing.T) {
	t.Parallel()

	repo := adapters.NewMemoryLoanRepository()
	testLoanRepositoryContract(t, repo)
}

package application_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wreckitral/pinjamean/internal/origination/domain/application"
)

func TestNewApplication(t *testing.T) {
	tests := []struct {
		name         string
		borrowerUUID string
		loanAmountIDR int64
		termMonths   int
		loanType     application.LoanType
		expectedAPR  float64
		expectedErr  error
	}{
		{
			name:          "valid KKB loan assigns 6.5 APR",
			borrowerUUID:  "borrower-123",
			loanAmountIDR: 150_000_000,
			termMonths:    36,
			loanType:      application.TypeKKB,
			expectedAPR:   6.5,
			expectedErr:   nil,
		},
		{
			name:          "valid KUR exact assigns 6.0 APR",
			borrowerUUID:  "borrower-123",
			loanAmountIDR: 50_000_000,
			termMonths:    12,
			loanType:      application.TypeKUR,
			expectedAPR:   6.0,
			expectedErr:   nil,
		},
		{
			name:          "empty borrower ID",
			borrowerUUID:  "",
			loanAmountIDR: 10_000_000,
			termMonths:    12,
			loanType:      application.TypeKTA,
			expectedErr:   application.ErrMissingBorrower,
		},
		{
			name:          "zero amount",
			borrowerUUID:  "borrower-123",
			loanAmountIDR: 0,
			termMonths:    12,
			loanType:      application.TypeKTA,
			expectedErr:   application.ErrInvalidAmount,
		},
		{
			name:          "invalid term months",
			borrowerUUID:  "borrower-123",
			loanAmountIDR: 10_000_000,
			termMonths:    0,
			loanType:      application.TypeKTA,
			expectedErr:   application.ErrInvalidTerm,
		},
		{
			name:          "unsupported loan type",
			borrowerUUID:  "borrower-123",
			loanAmountIDR: 10_000_000,
			termMonths:    12,
			loanType:      "INVALID_TYPE",
			expectedErr:   application.ErrUnsupportedType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := application.NewApplication(
				"app-uuid-1",
				tt.borrowerUUID,
				tt.loanAmountIDR,
				tt.termMonths,
				tt.loanType,
			)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				require.Nil(t, app)
			} else {
				require.NoError(t, err)
				require.NotNil(t, app)
				require.Equal(t, tt.loanAmountIDR, app.LoanAmountIDR())
				require.Equal(t, application.ApplicationStatusPending, app.Status())

				require.Equal(t, tt.expectedAPR, app.InterestRateAPR())
			}
		})
	}
}

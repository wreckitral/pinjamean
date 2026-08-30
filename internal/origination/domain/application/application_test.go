package application_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wreckitral/pinjamean/internal/origination/domain/application"
)

func TestNewApplication(t *testing.T) {
	tests := []struct {
		name            string
		borrowerUUID    string
		loanAmountIDR       int64
		termMonths      int
		loanType        application.LoanType
		interestRateAPR float64
		expectedErr     error
	}{
		{
			name:            "valid KKB loan",
			borrowerUUID:    "borrower-123",
			loanAmountIDR:       150_000_000,
			termMonths:      36,
			loanType:        application.TypeKKB,
			interestRateAPR: 10.5,
			expectedErr:     nil,
		},
		{
			name:            "valid KUR exact APR",
			borrowerUUID:    "borrower-123",
			loanAmountIDR:       50_000_000,
			termMonths:      12,
			loanType:        application.TypeKUR,
			interestRateAPR: 6.0,
			expectedErr:     nil,
		},
		{
			name:            "empty borrower ID",
			borrowerUUID:    "",
			loanAmountIDR:       10_000_000,
			termMonths:      12,
			loanType:        application.TypeKTA,
			interestRateAPR: 20.0,
			expectedErr:     application.ErrMissingBorrower,
		},
		{
			name:            "zero amount",
			borrowerUUID:    "borrower-123",
			loanAmountIDR:       0,
			termMonths:      12,
			loanType:        application.TypeKTA,
			interestRateAPR: 20.0,
			expectedErr:     application.ErrInvalidAmount,
		},
		{
			name:            "invalid term months",
			borrowerUUID:    "borrower-123",
			loanAmountIDR:       10_000_000,
			termMonths:      0,
			loanType:        application.TypeKTA,
			interestRateAPR: 20.0,
			expectedErr:     application.ErrInvalidTerm,
		},
		{
			name:            "unsupported loan type",
			borrowerUUID:    "borrower-123",
			loanAmountIDR:       10_000_000,
			termMonths:      12,
			loanType:        "INVALID_TYPE",
			interestRateAPR: 10.0,
			expectedErr:     application.ErrUnsupportedType,
		},
		{
			name:            "KPR APR too low",
			borrowerUUID:    "borrower-123",
			loanAmountIDR:       500_000_000,
			termMonths:      120,
			loanType:        application.TypeKPR,
			interestRateAPR: 6.9, // Min is 7.0
			expectedErr:     application.ErrAPRTooLow,
		},
		{
			name:            "KPR APR exactly min (Boundary)",
			borrowerUUID:    "borrower-123",
			loanAmountIDR:       500_000_000,
			termMonths:      120,
			loanType:        application.TypeKPR,
			interestRateAPR: 7.0, // Min is 7.0
			expectedErr:     nil,
		},
		{
			name:            "KTA APR too high",
			borrowerUUID:    "borrower-123",
			loanAmountIDR:       20_000_000,
			termMonths:      24,
			loanType:        application.TypeKTA,
			interestRateAPR: 30.1, // Max is 30.0
			expectedErr:     application.ErrAPRTooHigh,
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
				tt.interestRateAPR,
			)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				require.Nil(t, app)
			} else {
				require.NoError(t, err)
				require.NotNil(t, app)
				require.Equal(t, tt.loanAmountIDR, app.LoanAmountIDR())
				require.Equal(t, application.ApplicationStatusPending, app.Status())
			}
		})
	}
}

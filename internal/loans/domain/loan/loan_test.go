package loan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLoan(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		uuid         string
		borrowerUUID string
		amountIDR    int64
		termMonths   int
		loanType     LoanType
		wantErr      bool
		wantAPR      float64
	}{
		{
			name:         "missing borrower",
			uuid:         "app-1",
			borrowerUUID: "",
			amountIDR:    10_000_000,
			termMonths:   12,
			loanType:     TypeKTA,
			wantErr:      true,
		},
		{
			name:         "zero amount",
			uuid:         "app-1",
			borrowerUUID: "borrower-1",
			amountIDR:    0,
			termMonths:   12,
			loanType:     TypeKTA,
			wantErr:      true,
		},
		{
			name:         "negative amount",
			uuid:         "app-1",
			borrowerUUID: "borrower-1",
			amountIDR:    -500_000,
			termMonths:   12,
			loanType:     TypeKTA,
			wantErr:      true,
		},
		{
			name:         "zero term",
			uuid:         "app-1",
			borrowerUUID: "borrower-1",
			amountIDR:    10_000_000,
			termMonths:   0,
			loanType:     TypeKTA,
			wantErr:      true,
		},
		{
			name:         "negative term",
			uuid:         "app-1",
			borrowerUUID: "borrower-1",
			amountIDR:    10_000_000,
			termMonths:   -6,
			loanType:     TypeKTA,
			wantErr:      true,
		},
		{
			name:         "unsupported loan type",
			uuid:         "app-1",
			borrowerUUID: "borrower-1",
			amountIDR:    10_000_000,
			termMonths:   12,
			loanType:     LoanType("XYZ"),
			wantErr:      true,
		},
		{
			name:         "valid KPR application",
			uuid:         "app-1",
			borrowerUUID: "borrower-1",
			amountIDR:    500_000_000,
			termMonths:   180,
			loanType:     TypeKPR,
			wantAPR:      8.5,
		},
		{
			name:         "valid KKB application",
			uuid:         "app-1",
			borrowerUUID: "borrower-1",
			amountIDR:    150_000_000,
			termMonths:   48,
			loanType:     TypeKKB,
			wantAPR:      6.5,
		},
		{
			name:         "valid KTA application",
			uuid:         "app-1",
			borrowerUUID: "borrower-1",
			amountIDR:    50_000_000,
			termMonths:   24,
			loanType:     TypeKTA,
			wantAPR:      18.0,
		},
		{
			name:         "valid KMG application",
			uuid:         "app-1",
			borrowerUUID: "borrower-1",
			amountIDR:    20_000_000,
			termMonths:   12,
			loanType:     TypeKMG,
			wantAPR:      11.0,
		},
		{
			name:         "valid KMK application",
			uuid:         "app-1",
			borrowerUUID: "borrower-1",
			amountIDR:    100_000_000,
			termMonths:   36,
			loanType:     TypeKMK,
			wantAPR:      10.0,
		},
		{
			name:         "valid KUR application",
			uuid:         "app-1",
			borrowerUUID: "borrower-1",
			amountIDR:    25_000_000,
			termMonths:   12,
			loanType:     TypeKUR,
			wantAPR:      6.0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewLoan(tc.uuid, tc.borrowerUUID, tc.amountIDR, tc.termMonths, tc.loanType)

			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, tc.wantAPR, got.interestRateAPR)
		})
	}
}

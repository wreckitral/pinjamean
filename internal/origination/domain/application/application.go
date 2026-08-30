package application

import (
	"time"

	commonerrors "github.com/wreckitral/pinjamean/internal/common/errors"
)

var (
	ErrMissingBorrower = commonerrors.NewIncorrectInputError("borrower ID is required", "missing-borrower")
	ErrInvalidAmount   = commonerrors.NewIncorrectInputError("loan amount must be greater than zero", "invalid-amount")
	ErrInvalidTerm     = commonerrors.NewIncorrectInputError("loan term must be at least 1 month", "invalid-term")
	ErrUnsupportedType = commonerrors.NewIncorrectInputError("unsupported loan type", "unsupported-type")
	ErrAPRTooLow       = commonerrors.NewIncorrectInputError("APR too low", "apr-too-low")
	ErrAPRTooHigh      = commonerrors.NewIncorrectInputError("APR too high", "apr-too-high")
)

type LoanType string

const (
	TypeKPR LoanType = "KPR" // Mortgage
	TypeKKB LoanType = "KKB" // Vehicle
	TypeKTA LoanType = "KTA" // Unsecured Personal (Consumptive)
	TypeKMG LoanType = "KMG" // Multipurpose
	TypeKMK LoanType = "KMK" // Working Capital (Productive)
	TypeKUR LoanType = "KUR" // Government Microloan (Productive)
)

type ApplicationStatus string

const (
	ApplicationStatusPending  ApplicationStatus = "PENDING"
	ApplicationStatusApproved ApplicationStatus = "APPROVED"
	ApplicationStatusRejected ApplicationStatus = "REJECTED"
)

type Application struct {
	uuid string

	borrowerUUID string

	loanAmountIDR   int64
	termMonths      int
	loanType        LoanType
	interestRateAPR float64
	status          ApplicationStatus
	createdAt       time.Time
	updatedAt       time.Time
}

func NewApplication(uuid, borrowerUUID string, amountIDR int64, termMonths int, loanType LoanType, interestRateAPR float64) (*Application, error) {
	if borrowerUUID == "" {
		return nil, ErrMissingBorrower
	}

	if amountIDR <= 0 {
		return nil, ErrInvalidAmount
	}

	if termMonths <= 0 {
		return nil, ErrInvalidTerm
	}

	if err := validateProductAPR(loanType, interestRateAPR); err != nil {
		return nil, err
	}

	return &Application{
		uuid: uuid,
		borrowerUUID: borrowerUUID,
		loanAmountIDR: amountIDR,
		termMonths: termMonths,
		loanType: loanType,
		interestRateAPR: interestRateAPR,
		status: ApplicationStatusPending,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

func (a Application) LoanAmountIDR() int64 {
	return a.loanAmountIDR
}

func (a Application) Status() ApplicationStatus {
	return a.status
}

func validateProductAPR(loanType LoanType, apr float64) error {
	switch loanType {
	case TypeKPR:
		if apr < 7.0 { return ErrAPRTooLow }
		if apr > 13.0 { return ErrAPRTooHigh }
	case TypeKKB, TypeKMG:
		if apr < 8.0 { return ErrAPRTooLow }
		if apr > 18.0 { return ErrAPRTooHigh }
	case TypeKTA:
		if apr < 15.0 { return ErrAPRTooLow }
		if apr > 30.0 { return ErrAPRTooHigh }
	case TypeKMK:
		if apr < 9.0 { return ErrAPRTooLow }
		if apr > 15.0 { return ErrAPRTooHigh }
	case TypeKUR:
		if apr != 6.0 { return ErrAPRTooHigh }
	default:
		return ErrUnsupportedType
	}
	return nil
}

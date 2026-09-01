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

func NewApplication(uuid, borrowerUUID string, amountIDR int64, termMonths int, loanType LoanType) (*Application, error) {
	if borrowerUUID == "" {
		return nil, ErrMissingBorrower
	}

	if amountIDR <= 0 {
		return nil, ErrInvalidAmount
	}

	if termMonths <= 0 {
		return nil, ErrInvalidTerm
	}

	var assignedAPR float64
	switch loanType {
	case TypeKPR:
		assignedAPR = 8.5
	case TypeKKB:
		assignedAPR = 6.5
	case TypeKTA:
		assignedAPR = 18.0
	case TypeKMG:
		assignedAPR = 11.0
	case TypeKMK:
		assignedAPR = 10.0
	case TypeKUR:
		assignedAPR = 6.0
	default:
		return nil, ErrUnsupportedType
	}

	return &Application{
		uuid:            uuid,
		borrowerUUID:    borrowerUUID,
		loanAmountIDR:   amountIDR,
		termMonths:      termMonths,
		loanType:        loanType,
		interestRateAPR: assignedAPR,
		status:          ApplicationStatusPending,
		createdAt:       time.Now(),
		updatedAt:       time.Now(),
	}, nil
}

func (a Application) LoanAmountIDR() int64 {
	return a.loanAmountIDR
}

func (a Application) Status() ApplicationStatus {
	return a.status
}

func (a Application) InterestRateAPR() float64 {
	return a.interestRateAPR
}

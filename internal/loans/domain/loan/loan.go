package loan

import (
	"time"
	"errors"

	commonerrors "github.com/wreckitral/pinjamean/internal/common/errors"
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

type LoanStatus string

const (
	LoanStatusPending  LoanStatus = "PENDING"
	LoanStatusApproved LoanStatus = "APPROVED"
	LoanStatusRejected LoanStatus = "REJECTED"
)

type Loan struct {
	uuid string

	borrowerUUID string

	loanAmountIDR   int64
	termMonths      int
	loanType        LoanType
	interestRateAPR float64
	status          LoanStatus
	createdAt       time.Time
	updatedAt       time.Time
}

func NewLoan(uuid, borrowerUUID string, amountIDR int64, termMonths int, loanType LoanType) (*Loan, error) {
	if borrowerUUID == "" {
		return nil, errors.New("borrowerUUID is required")
	}

	if amountIDR <= 0 {
		return nil, errors.New("loan amount must be greater than 0")
	}

	if termMonths <= 0 {
		return nil, errors.New("loan term must be at least 1 month")
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
		return nil, errors.New("unsupported loan type")
	}

	return &Loan{
		uuid:            uuid,
		borrowerUUID:    borrowerUUID,
		loanAmountIDR:   amountIDR,
		termMonths:      termMonths,
		loanType:        loanType,
		interestRateAPR: assignedAPR,
		status:          LoanStatusPending,
		createdAt:       time.Now(),
		updatedAt:       time.Now(),
	}, nil
}

func (l *Loan) UUID() string              { return l.uuid }
func (l *Loan) BorrowerUUID() string      { return l.borrowerUUID }
func (l *Loan) AmountIDR() int64          { return l.loanAmountIDR }
func (l *Loan) TermMonths() int           { return l.termMonths }
func (l *Loan) Type() LoanType            { return l.loanType }
func (l *Loan) InterestRateAPR() float64  { return l.interestRateAPR }
func (l *Loan) Status() LoanStatus        { return l.status }

var ErrLoanStatusNotPending = commonerrors.NewIncorrectInputError("loan must be pending", "not-pending")

// state-transition
func (a *Loan) Approve() error {
	if a.status != LoanStatusPending {
		return ErrLoanStatusNotPending
	}

	a.status = LoanStatusApproved
	a.updatedAt = time.Now()

	return nil
}

func (a *Loan) Reject() error {
	if a.status != LoanStatusPending {
		return ErrLoanStatusNotPending
	}

	a.status = LoanStatusRejected
	a.updatedAt = time.Now()

	return nil
}

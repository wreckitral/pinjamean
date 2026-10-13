package loan

import (
	"errors"
	"strings"
	"time"

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

	borrowerUUID   string
	officerUUID    string
	assignedToUUID string
	analystUUID    string

	loanAmountIDR   int64
	termMonths      int
	loanType        LoanType
	interestRateAPR float64
	status          LoanStatus
	rejectReason    string
	createdAt       time.Time
	updatedAt       time.Time
}

func NewLoan(uuid, borrowerUUID, officerUUID, assignedToUUID string, amountIDR int64, termMonths int, loanType LoanType) (*Loan, error) {
	if borrowerUUID == "" {
		return nil, errors.New("borrowerUUID is required")
	}

	if officerUUID == "" {
		return nil, errors.New("officerUUID is required")
	}

	if assignedToUUID == "" {
		return nil, errors.New("assignedToUUID is required")
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
		officerUUID:     officerUUID,
		assignedToUUID:  assignedToUUID,
		loanAmountIDR:   amountIDR,
		termMonths:      termMonths,
		loanType:        loanType,
		interestRateAPR: assignedAPR,
		status:          LoanStatusPending,
		createdAt:       time.Now(),
		updatedAt:       time.Now(),
	}, nil
}

func (l *Loan) UUID() string             { return l.uuid }
func (l *Loan) BorrowerUUID() string     { return l.borrowerUUID }
func (l *Loan) OfficerUUID() string      { return l.officerUUID }
func (l *Loan) AssignedToUUID() string   { return l.assignedToUUID }
func (l *Loan) AnalystUUID() string      { return l.analystUUID }
func (l *Loan) RejectionReason() string  { return l.rejectReason }
func (l *Loan) AmountIDR() int64         { return l.loanAmountIDR }
func (l *Loan) TermMonths() int          { return l.termMonths }
func (l *Loan) Type() LoanType           { return l.loanType }
func (l *Loan) InterestRateAPR() float64 { return l.interestRateAPR }
func (l *Loan) Status() LoanStatus       { return l.status }
func (l *Loan) CreatedAt() time.Time     { return l.createdAt }
func (l *Loan) UpdatedAt() time.Time     { return l.updatedAt }

var ErrLoanStatusNotPending = commonerrors.NewIncorrectInputError("loan must be pending", "not-pending")
var ErrMissingRejectionReason = commonerrors.NewIncorrectInputError("rejection reason is required", "reason-missing")

// state-transition
func (l *Loan) Approve(analystUUID string) error {
	if l.status != LoanStatusPending {
		return ErrLoanStatusNotPending
	}

	l.status = LoanStatusApproved
	l.analystUUID = analystUUID
	l.updatedAt = time.Now()

	return nil
}

func (l *Loan) Reject(analystUUID, reason string) error {
	if l.status != LoanStatusPending {
		return ErrLoanStatusNotPending
	}

	reason = strings.TrimSpace(reason)
	if reason == "" {
		return ErrMissingRejectionReason
	}

	l.status = LoanStatusRejected
	l.analystUUID = analystUUID
	l.rejectReason = reason
	l.updatedAt = time.Now()

	return nil
}

func (l *Loan) Reassign(ReassignToUUID string) error {
	if ReassignToUUID == "" {
		return errors.New("reassignToUUID is required")
	}

	l.assignedToUUID = ReassignToUUID
	l.updatedAt = time.Now()

	return nil
}

func UnmarshalLoanFromDatabase(
	uuid string,
	borrowerUUID string,
	officerUUID string,
	assignedToUUID string,
	analystUUID string,
	amountIDR int64,
	termMonths int,
	loanType LoanType,
	interestRateAPR float64,
	status LoanStatus,
	rejectReason string,
	createdAt time.Time,
	updatedAt time.Time,
) (*Loan, error) {
	l, err := NewLoan(
		uuid,
		borrowerUUID,
		officerUUID,
		assignedToUUID,
		amountIDR,
		termMonths,
		loanType,
	)
	if err != nil {
		return nil, err
	}

	l.analystUUID = analystUUID
	l.interestRateAPR = interestRateAPR
	l.status = status
	l.rejectReason = rejectReason
	l.createdAt = createdAt
	l.updatedAt = updatedAt

	return l, nil
}

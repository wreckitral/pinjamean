package app

import (
	"github.com/wreckitral/pinjamean/internal/loans/app/command"
	"github.com/wreckitral/pinjamean/internal/loans/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	SubmitLoan   command.SubmitLoanHandler
	ReassignLoan command.ReassignLoanHandler
	ApproveLoan  command.ApproveLoanHandler
	RejectLoan   command.RejectLoanHandler
}

type Queries struct {
	GetLoanByID query.GetLoanByIDHandler
}

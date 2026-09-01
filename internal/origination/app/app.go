package app

import (
	"github.com/wreckitral/pinjamean/internal/origination/app/command"
	"github.com/wreckitral/pinjamean/internal/origination/app/query"
)

type Application struct {
	Commands Commands
	Queries Queries

}

type Commands struct {
	SubmitLoan command.SubmitLoanHandler
}

type Queries struct {
	GetApplication query.GetApplicationHandler
}

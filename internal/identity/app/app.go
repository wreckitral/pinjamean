package app

import "github.com/wreckitral/pinjamean/internal/identity/app/command"
import "github.com/wreckitral/pinjamean/internal/identity/app/query"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	SaveAccount command.SaveAccountHandler
}

type Queries struct {
	GetAccountByEmail query.GetAccountByEmailHandler
	Login query.LoginHandler
}

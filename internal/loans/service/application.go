package service

import (
	"context"
	"log/slog"

	"github.com/wreckitral/pinjamean/internal/loans/adapters"
	"github.com/wreckitral/pinjamean/internal/loans/app"
	"github.com/wreckitral/pinjamean/internal/loans/app/command"
	"github.com/wreckitral/pinjamean/internal/loans/app/query"
)

func NewApplication(ctx context.Context) app.Application {
	db, err := adapters.NewPostgresSQLConnection()
	if err != nil {
		panic(err)
	}
	repo := adapters.NewPostgresLoanRepository(db)

	logger := slog.Default()

	return app.Application{
		Commands: app.Commands{
			SubmitLoan: command.NewSubmitLoanHandler(repo, logger),
		},
		Queries: app.Queries{
			GetLoanByID: query.NewGetLoanByIDHandler(repo, logger),
		},
	}
}

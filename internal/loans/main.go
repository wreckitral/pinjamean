package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wreckitral/pinjamean/internal/loans/adapters"
	"github.com/wreckitral/pinjamean/internal/loans/app"
	"github.com/wreckitral/pinjamean/internal/loans/app/command"
	"github.com/wreckitral/pinjamean/internal/loans/app/query"
	"github.com/wreckitral/pinjamean/internal/loans/ports"
)

func main() {
	ctx := context.Background()

	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	repo := adapters.NewPostgresLoanRepository(db)

	submitLoanHandler := command.NewSubmitLoanHandler(repo)
	getLoanHandler := query.NewGetLoanByIDHandler(repo)

	application := app.Application{
		Commands: app.Commands{
			SubmitLoan: submitLoanHandler,
		},
		Queries: app.Queries{
			GetLoanByID: getLoanHandler,
		},
	}

	router := chi.NewRouter()
	httpServer := ports.NewHttpServer(application)

	router.Mount("/", ports.HandlerFromMux(httpServer, router))

	if err := http.ListenAndServe(":7777", router); err != nil {
		panic(err)
	}
}

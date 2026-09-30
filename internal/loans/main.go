package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wreckitral/pinjamean/internal/common/auth"
	"github.com/wreckitral/pinjamean/internal/common/logs"
	"github.com/wreckitral/pinjamean/internal/common/server"
	identityadapters "github.com/wreckitral/pinjamean/internal/identity/adapters"
	identityapp "github.com/wreckitral/pinjamean/internal/identity/app"
	identitycommand "github.com/wreckitral/pinjamean/internal/identity/app/command"
	identityquery "github.com/wreckitral/pinjamean/internal/identity/app/query"
	identityports "github.com/wreckitral/pinjamean/internal/identity/ports"
	"github.com/wreckitral/pinjamean/internal/loans/ports"
	"github.com/wreckitral/pinjamean/internal/loans/service"
)

func main() {
	logs.Init()

	ctx := context.Background()

	application := service.NewApplication(ctx)
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		panic("JWT_SECRET must be set")
	}

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))

	switch serverType {
	case "http":
		server.RunHTTPServer(func(router chi.Router) http.Handler {
			db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
			if err != nil {
				panic(err)
			}

			userRepo := identityadapters.NewPostgresUserRepository(db)
			login := identityquery.NewLoginHandler(userRepo, jwtSecret, slog.Default())
			router.Mount("/auth/login", identityports.LoginHandler(login))

			identityHTTP := identityports.NewHttpServer(identityapp.Application{
				Commands: identityapp.Commands{
					SaveAccount: identitycommand.NewSaveAccountHandler(userRepo, slog.Default()),
				},
				Queries: identityapp.Queries{
					GetAccountByEmail: identityquery.NewGetAccountByEmailHandler(userRepo, slog.Default()),
					Login:             login,
				},
			})
			router.Post("/accounts", identityHTTP.SaveAccount)

			secured := chi.NewRouter()
			secured.Use(auth.Middleware(jwtSecret))

			ports.HandlerFromMux(ports.NewHttpServer(application), secured)
			router.Mount("/", secured)
			return router
		})
	}
}

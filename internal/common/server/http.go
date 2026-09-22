package server

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/wreckitral/pinjamean/internal/common/logs"
)

func RunHTTPServer(createHandler func(router chi.Router) http.Handler) {
	RunHTTPServerOnAddr(":"+os.Getenv("PORT"), createHandler)
}

func RunHTTPServerOnAddr(addr string, createHandler func(router chi.Router) http.Handler) {
	apiRouter := chi.NewRouter()

	setMiddlewares(apiRouter)

	rootRouter := chi.NewRouter()
	rootRouter.Mount("/api", createHandler(apiRouter))

	slog.Info("Starting HTTP server")

	server := &http.Server{
		Addr: addr,
		Handler: rootRouter,
		ReadHeaderTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		slog.Error("Unable to start HTTP server", "error", err)
		panic(err)
	}
}

func setMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(logs.NewStructuredLogger(slog.Default()))
	router.Use(middleware.Recoverer)

	addCorsMiddleware(router)

	router.Use(middleware.Recoverer)

	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.SetHeader("X-Frame-Options", "deny"),
	)

	router.Use(middleware.NoCache)
}

func addCorsMiddleware(router *chi.Mux) {
	allowedOrigins := strings.Split(
		os.Getenv("CORS_ALLOWED_ORIGINS"),
		";",
	)

	if len(allowedOrigins) == 1 && allowedOrigins[0] == "" {
		return
	}

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: false,
		MaxAge:           300,
	}))
}

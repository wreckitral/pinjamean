package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/wreckitral/pinjamean/internal/common/logs"
	"github.com/wreckitral/pinjamean/internal/common/server"
	"github.com/wreckitral/pinjamean/internal/loans/ports"
	"github.com/wreckitral/pinjamean/internal/loans/service"
)

func main() {
	logs.Init()

	ctx := context.Background()

	application := service.NewApplication(ctx)

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))

	switch serverType {
	case "http":
		server.RunHTTPServer(func(router chi.Router) http.Handler{
			return ports.HandlerFromMux(
				ports.NewHttpServer(application),
				router,
			)
		})
	}
}

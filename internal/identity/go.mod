module github.com/wreckitral/pinjamean/internal/identity

go 1.27.1

replace github.com/wreckitral/pinjamean/internal/common => ../common

require (
	github.com/go-chi/chi/v5 v5.3.2
	github.com/go-chi/render v1.0.3
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.11.0
	github.com/oapi-codegen/runtime v1.7.0
	github.com/stretchr/testify v1.12.1
	github.com/wreckitral/pinjamean/internal/common v0.0.0-20260922143626-e0655ff4dd89
	golang.org/x/crypto v0.57.0
)

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/text v0.42.0 // indirect
)

module github.com/wreckitral/pinjamean/internal/origination

go 1.27.0

replace github.com/wreckitral/pinjamean/internal/common => ../common

require (
	github.com/go-chi/chi/v5 v5.3.2
	github.com/stretchr/testify v1.12.1
	github.com/wreckitral/pinjamean/internal/common v0.0.0-00010101000000-000000000000
)

require go.yaml.in/yaml/v3 v3.0.5 // indirect

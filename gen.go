package gen

//go:generate go run cmd/tools/terndotenv/main.go
//go:generate sqlc generate -f ./internal/data/db/pgstore/sqlc.yaml

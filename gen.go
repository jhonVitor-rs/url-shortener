package gen

//go:generate go run cmd/tools/terndotenv/main.go
//go:generate sqlc generate -f ./internal/adapters/secondary/persistence/pgstore/sqlc.yaml

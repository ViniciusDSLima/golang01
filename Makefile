build:
	@go build -o bin/golang01 ./cmd/main.go

run: build
	@./bin/golang01

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations -seq

migrate-up:
	@go run ./cmd/migrate/migrations/main.go up

migrate-down:
	@go run ./cmd/migrate/migrations/main.go down
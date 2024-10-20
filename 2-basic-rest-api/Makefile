.PHONY: all build test clean run stop

all: build test

build:
	@go build -o bin/myapp cmd/main.go

run:
	@./bin/myapp

test:
	@go test -v ./...

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

tidy:
	@go mod tidy
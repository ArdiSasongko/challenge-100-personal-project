.PHONY: all build test clean run stop

all: build test

tidy:
	@go mod tidy

build:
	@go build -o cmd/bin/myapp cmd/myapp/main.go

run:
	@./cmd/bin/myapp

migration:
	@migrate create -ext sql -dir migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run migrate/main.go up

migrate-down:
	@go run migrate/main.go down

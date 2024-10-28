.PHONY: all build test clean run stop

all: build test

tidy:
	@go mod tidy

build:
	@go build -o cmd/bin/myapp cmd/myapp/main.go

run:
	@./cmd/bin/myapp
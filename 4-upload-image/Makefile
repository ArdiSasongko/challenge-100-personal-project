export PG_URL='postgres://root:mypassword@localhost:5432/upload_image?sslmode=disable'

migrate-create:
	@migrate create -ext sql -dir scripts/migrations -seq $(name)

migrate-up:
	@migrate -database ${PG_URL} -path scripts/migrations up

migrate-down:
	@migrate -database ${PG_URL} -path scripts/migrations down

tidy:
	@go mod tidy

build:
	@go build -o cmd/bin/main cmd/main.go

run:
	@./cmd/bin/main

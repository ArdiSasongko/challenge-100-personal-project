export POSTGREST_SQL='postgres://root:mypassword@localhost:5432/forum_local?sslmode=disable'

migrate-create:
	@migrate create -ext sql -dir scripts/migrations -seq $(name)

migrate-up:
	@migrate -database ${POSTGREST_SQL} -path scripts/migrations up

migrate-down:
	@migrate -database ${POSTGREST_SQL} -path scripts/migrations down

migrate-down-one:
	@migrate -database ${POSTGREST_SQL} -path scripts/migrations down 1
	
tidy:
	@go mod tidy

build:
	@go build -o cmd/bin/main cmd/main.go

run:
	@./cmd/bin/main
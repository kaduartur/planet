# Go parameters
GO_CMD=go
GO_RUN=$(GO_CMD) run
GO_TEST=$(GO_CMD) test
GO_TOOL_COVER=$(GO_CMD) tool cover

BIN=bin
ENVS=DATABASE_HOST=127.0.0.1 DATABASE_PORT=27017 DATABASE_USER=planets DATABASE_PASS=planets DATABASE_NAME=planets_db SWAPI_URL=https://swapi.dev

run-command-api:
	${ENVS} $(GO_RUN) ./cmd/command/main.go

run-query-api:
	${ENVS} $(GO_RUN) ./cmd/query/main.go

run-processor:
	${ENVS} $(GO_RUN) ./cmd/processor/main.go

run-test:
	${ENVS} $(GO_TEST) -race -coverprofile=coverage.txt -covermode=atomic `go list ./... | grep -v vendor/`
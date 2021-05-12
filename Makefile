# Go parameters
GO_CMD=go
GO_RUN=$(GO_CMD) run
GO_TEST=$(GO_CMD) test
GO_TOOL_COVER=$(GO_CMD) tool cover

DB_ENVS=DATABASE_HOST=127.0.0.1 DATABASE_PORT=27017 DATABASE_USER=planets DATABASE_PASS=planets DATABASE_NAME=planets_db
KAFKA_ENVS=KAFKA_SERVER=localhost KAFKA_PLANET_TOPIC=planet-processor

run-command-api:
	${DB_ENVS} ${KAFKA_ENVS} $(GO_RUN) ./cmd/command/main.go

run-query-api:
	${DB_ENVS} $(GO_RUN) ./cmd/query/main.go

run-processor:
	${DB_ENVS} ${KAFKA_ENVS} SWAPI_URL=https://swapi.dev $(GO_RUN) ./cmd/processor/main.go

run-test:
	${DB_ENVS} ${KAFKA_ENVS} SWAPI_URL=http://localhost:8882 $(GO_TEST) -race -coverprofile=coverage.txt -covermode=atomic `go list ./... | grep -v vendor/`
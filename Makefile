ENVS=DATABASE_HOST=127.0.0.1 DATABASE_PORT=27017 DATABASE_USER=planets DATABASE_PASS=planets DATABASE_NAME=planets_db SWAPI_URL=https://swapi.dev

run-command-api:
	${ENVS} go run ./cmd/command/main.go

run-query-api:
	${ENVS} go run ./cmd/query/main.go

run-processor:
	${ENVS} go run ./cmd/processor/main.go
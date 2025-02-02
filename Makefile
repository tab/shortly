DB_USER=postgres
DB_DEV_NAME=shortly-development
DB_TEST_NAME=shortly-test
DB_HOST=localhost
DB_PORT=5432

DB_SERVICE_NAME=database

GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=db/migrate

GO_ENV ?= development

SERVER_ADDRESS=localhost:8080
PROFILE_ADDRESS=localhost:2080

#PROFILE_NAME=base | result
PROFILE_NAME=result

ifeq ($(GO_ENV),test)
	ENV_FILE=.env.test
	LOCAL_ENV_FILE=.env.test.local
	DB_NAME=$(DB_TEST_NAME)
	GOOSE_DBSTRING="host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_USER) dbname=$(DB_TEST_NAME) sslmode=disable"
else
	ENV_FILE=.env.development
	LOCAL_ENV_FILE=.env.development.local
	DB_NAME=$(DB_DEV_NAME)
	GOOSE_DBSTRING="host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_USER) dbname=$(DB_DEV_NAME) sslmode=disable"
endif

ifneq (,$(wildcard $(ENV_FILE)))
	include $(ENV_FILE)
	export $(shell sed 's/=.*//' $(ENV_FILE))
endif

ifneq (,$(wildcard $(LOCAL_ENV_FILE)))
	include $(LOCAL_ENV_FILE)
	export $(shell sed 's/=.*//' $(LOCAL_ENV_FILE))
endif

.PHONY: db\:create
db\:create:
	@echo "Creating $(DB_NAME) database..."
	docker-compose exec -T $(DB_SERVICE_NAME) createdb -U $(DB_USER) $(DB_NAME)

.PHONY: db\:drop
db\:drop:
	@echo "Dropping $(DB_NAME) database..."
	docker-compose exec -T $(DB_SERVICE_NAME) dropdb -U $(DB_USER) --if-exists $(DB_NAME)

.PHONY: db\:migrate
db\:migrate:
	@echo "Applying migrations to $(DB_NAME) database..."
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) up

.PHONY: db\:migrate\:status
db\:migrate\:status:
	@echo "Migration status in $(DB_NAME) database..."
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) status

.PHONY: db\:rollback
db\:rollback:
	@echo "Rolling back last migration in $(DB_NAME) database..."
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir $(GOOSE_MIGRATION_DIR) down

.PHONY: db\:schema\:dump
db\:schema\:dump:
	@echo "Dumping database schema from $(DB_NAME) to db/schema.sql..."
	docker-compose exec -T $(DB_SERVICE_NAME) pg_dump -U $(DB_USER) -d $(DB_NAME) --schema-only --exclude-table=goose_db_version > db/schema.sql

.PHONY: db\:schema\:load
db\:schema\:load:
	@echo "Loading schema from db/schema.sql into $(DB_NAME) database..."
	cat db/schema.sql | docker-compose exec -T $(DB_SERVICE_NAME) psql -U $(DB_USER) -d $(DB_NAME)

.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	golangci-lint run

.PHONY: lint\:fix
lint\:fix:
	@echo "Running golangci-lint --fix ..."
	golangci-lint run --fix

.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

.PHONY: staticcheck
staticcheck:
	@echo "Running staticcheck..."
	staticcheck ./...

.PHONY: staticlint\:build
staticlint\:build:
	@echo "Building staticcheck binary..."
	go build -ldflags="-s -w" -o cmd/staticlint/staticlint cmd/staticlint/main.go
	chmod +x cmd/staticlint/staticlint

.PHONY: staticlint
staticlint:
	@echo "Running staticlint..."
	./cmd/staticlint/staticlint ./...

.PHONY: test
test:
	@echo "Running tests..."
	go test -cover ./...

.PHONY: coverage
coverage:
	@echo "Generating test coverage report..."
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

.PHONY: benchmark\:payload
benchmark\:payload:
	@echo "Running payload benchmark..."
	wrk -t16 -c100 -d30s -s scripts/benchmark.lua http://$(SERVER_ADDRESS)/api/shorten

.PHONY: pprof\:cpu
pprof\:cpu:
	@echo "Running pprof for CPU profiling..."
	curl --location "http://$(PROFILE_ADDRESS)/debug/pprof/profile?seconds=30" > profiles/$(PROFILE_NAME).cpu.pprof
	go tool pprof profiles/$(PROFILE_NAME).cpu.pprof

.PHONY: pprof\:cpu\:diff
pprof\:cpu\:diff:
	@echo "Running pprof CPU diff..."
	go tool pprof -top -diff_base=profiles/base.cpu.pprof profiles/result.cpu.pprof

.PHONY: pprof\:mem
pprof\:mem:
	@echo "Running pprof for memory profiling..."
	curl --location "http://$(PROFILE_ADDRESS)/debug/pprof/heap" > profiles/$(PROFILE_NAME).mem.pprof
	go tool pprof profiles/$(PROFILE_NAME).mem.pprof

.PHONY: pprof\:mem\:diff
pprof\:mem\:diff:
	@echo "Running pprof memory diff..."
	go tool pprof -top -diff_base=profiles/base.mem.pprof profiles/result.mem.pprof

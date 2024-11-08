DB_USER=postgres
DB_DEV_NAME=shortly-development
DB_TEST_NAME=shortly-test
DB_HOST=localhost
DB_PORT=5432

DB_SERVICE_NAME=database

GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=db/migrate

GO_ENV ?= development

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

.PHONY: test
test:
	@echo "Running tests..."
	go test -cover ./...

.PHONY: coverage
coverage:
	@echo "Generating test coverage report..."
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

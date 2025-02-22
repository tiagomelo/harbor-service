include .env
export

# ==============================================================================
# Help

.PHONY: help
## help: shows this help message
help:
	@ echo "Usage: make [target]\n"
	@ sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==============================================================================
# Database migrations

.PHONY: migrate-setup
## migrate-setup: installs golang-migrate
migrate-setup:
	@ go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: create-migrations
## create-migration: creates up and down migration files for a given name (make create-migrations NAME=<desired_name>)
create-migration: migrate-setup
	@ if [ -z "$(NAME)" ]; then echo >&2 please set the name of the migration via the variable NAME; exit 2; fi
	@ migrate create -ext sql -dir storage/sqlite/migrations -seq -digits 4 $(NAME)

.PHONY: migrate-up
## migrate-up: runs up N migrations, N is optional (make migrate-up N=<desired_migration_number>)
migrate-up: migrate-setup
	@ migrate -database "sqlite3://$(HARBOR_SVC_DB_FILE)?query" -path storage/sqlite/migrations up $(N)

.PHONY: migrate-down
## migrate-down: runs down N migrations, N is optional (make migrate-down N=<desired_migration_number>)
migrate-down: migrate-setup
	@ migrate -database "sqlite3://$(HARBOR_SVC_DB_FILE)?query" -path storage/sqlite/migrations down $(N)

.PHONY: migrate-to-version
## migrate-to-version: migrates to version V (make migrate-to-version V=<desired_version>)
migrate-to-version: migrate-setup
	@ if [ -z "$(V)" ]; then echo >&2 please set the desired version via the variable V; exit 2; fi
	@ migrate -database "sqlite3://$(HARBOR_SVC_DB_FILE)?query" -path storage/sqlite/migrations goto $(V)

.PHONY: migrate-force-version
## migrate-force-version: forces version V (make migrate-force-version V=<desired_version>)
migrate-force-version: migrate-setup
	@ if [ -z "$(V)" ]; then echo >&2 please set the desired version via the variable V; exit 2; fi
	@ migrate -database "sqlite3://$(HARBOR_SVC_DB_FILE)?query" -path storage/sqlite/migrations force $(V)

.PHONY: migrate-version
## migrate-version: checks current database migrations version
migrate-version: migrate-setup
	@ migrate -database "sqlite3://$(HARBOR_SVC_DB_FILE)?query" -path storage/sqlite/migrations version

.PHONY: migrate-test-up
## migrate-test-up: runs up N migrations on test db, N is optional (make migrate-up N=<desired_migration_number>)
migrate-test-up: migrate-setup
	@ migrate -database 'sqlite3://$(HARBOR_SVC_TEST_DB_FILE)?query' -path storage/sqlite/migrations up $(N)

.PHONY: migrate-test-down
## migrate-test-down: runs down N migrations on test db, N is optional (make migrate-down N=<desired_migration_number>)
migrate-test-down: migrate-setup
	@ migrate -database 'sqlite3://$(HARBOR_SVC_TEST_DB_FILE)' -path storage/sqlite/migrations down $(N)

# ==============================================================================
# DB

.PHONY: sqilite-console
## sqlite-console: opens the sqlite3 console
sqlite-console:
	@ sqlite3 $(HARBOR_SVC_DB_FILE)

# ==============================================================================
# Swagger

.PHONY: swagger
## swagger: generates api's documentation
swagger: 
	@ swagger generate spec -o http/doc/swagger.json --scan-models

.PHONY: swagger-ui
## swagger-ui: launches Swagger UI
swagger-ui: swagger
	@docker run --rm --name harbor-swagger-ui -p 8080:8080 \
		-e SWAGGER_JSON=/docs/swagger.json \
		-v $(shell pwd)/http/doc:/docs swaggerapi/swagger-ui

# ==============================================================================
# Tests

.PHONY: test
## test: run tests
test: migrate-test-up
	@ go test -v -race ./... -count=1

.PHONY: coverage
## coverage: run tests and generate coverage report in html format
coverage: migrate-test-up
coverage:
	@ packages=$$(go list ./... | grep -v "cmd" | grep -v "validate"); \
	if [ -z "$$packages" ]; then \
		echo "No valid Go packages found"; \
		exit 1; \
	fi; \
	go test -race -coverpkg=$$(echo $$packages | tr ' ' ',') -coverprofile=coverage.out $$packages && go tool cover -html=coverage.out

# ==============================================================================
# Quality Checks

.PHONY: fmt
## fmt: formats Go code
fmt:
	@ echo "Running go fmt..."
	@ go fmt ./...

.PHONY: vet
## vet: runs Go vet to analyze code for potential issues
vet:
	@ echo "Running go vet..."
	@ go vet ./...

.PHONY: govulncheck
## govulncheck: runs Go vulnerability check
govulncheck:
	@ go install golang.org/x/vuln/cmd/govulncheck@latest
	@ echo "Running go vuln check..."
	@ govulncheck ./...

.PHONY: lint
## lint: runs linter for all packages
lint: 
	@ echo "Running golangci-lint..."
	@ docker run  --rm -v "`pwd`:/workspace:cached" -w "/workspace/." golangci/golangci-lint:latest golangci-lint run

.PHONY: build
## build: checks if the code compiles correctly
build:
	@ echo "Running go build..."
	@ go build ./...

.PHONY: tidy
## tidy: ensures go.mod and go.sum are clean
tidy:
	@ echo "Running go mod tidy..."
	@ go mod tidy

# ==============================================================================
# App's execution

.PHONY: run pre-run-checks
## pre-run-checks: runs Go quality checks before starting the app
pre-run-checks: fmt vet govulncheck lint test build tidy

.PHONY: run
## run: runs the API
run: pre-run-checks migrate-up
	@ if [ -z "$(PORT)" ]; then echo >&2 please set the desired port via the variable PORT; exit 2; fi
	@ go run cmd/main.go -p $(PORT)

# ==============================================================================
# Docker

.PHONY: docker-build
## docker-build: builds the Docker image
docker-build:
	@ docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
## docker-run: runs the Docker container
docker-run: docker-build
	@ if [ -z "$(PORT)" ]; then echo >&2 please set the desired port via the variable PORT; exit 2; fi
	@ docker run --rm --name harbor-service -p $(PORT):$(PORT) -e PORT=$(PORT) -v harbor-db:/root/storage/sqlite harbor-service

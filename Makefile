include .env

# =========================================================================================== #
# HELPERS
# =========================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# =========================================================================================== #
# DEVELOPMENT
# =========================================================================================== #

## run/api: run the API application
.PHONY: run/api
run/api:
	go run ./cmd/api

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DB_DSN_LOCAL}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${DB_DSN_LOCAL} up

## run/docker/api: run the docker image
.PHONY: run/docker/api
run/docker/api:
	docker run -p 8080:8080 --env_file .env pr-service:latest -config=config/config.toml

## run/docker-compose/up: run the docker-compose(docker-compose.yml) stack in detached mode
.PHONY: run/docker-compose/up
run/docker-compose/up:
	docker-compose up -d

## stop/docker-compose/down: stop the docker-compose(docker-compose.yml) stack
.PHONY: run/docker-compose/down
run/docker-compose/down:
	docker-compose down

# =========================================================================================== #
# QUALITY CONTROL
# =========================================================================================== #
 
## audit: tidy and vendor dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	golangci-lint run --config=.golangci.yml
	@echo 'Running tests...'
	go test -race -vet=off -short ./...
 
## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# =========================================================================================== #
# BUILD
# =========================================================================================== #
 
## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api

## build/docker: build the docker image
.PHONY: build/docker
build/docker:
	@echo 'Building docker image...'
	docker build --tag pr-service .

## build-and-push/docker: build the docker image and push it to docker hub
.PHONY: build-and-push/docker
build-and-push/docker:
	@echo 'Building docker image...'
	docker build --tag vladgrskkh/pr-service .
	@echo 'Pushing docker image...'
	docker push vladgrskkh/pr-service:latest

# =========================================================================================== #
# E2E
# =========================================================================================== #

API_URL=http://localhost:8080

DC=docker-compose.yml

.PHONY: e2e/up
e2e/up:
	docker-compose -f $(DC) up -d --build

.PHONY:
e2e/test: e2e/test
	go test -v ./internal/e2e

.PHONY:
e2e/down: e2e/down
	docker-compose -f $(DC) down -v

.PHONY: e2e
e2e: e2e/down e2e/up e2e/migrate e2e/test e2e/down
	@echo "E2E tests complete successfuly!"

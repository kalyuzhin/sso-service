include .env
export $(shell sed 's/=.*//' .env)


APP_NAME=$(basename -s .git "$(git config --get remote.origin.url)")


GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

DB_SETUP := user=$(DB_USERNAME) password=$(DB_PASSWORD) dbname=$(DB_NAME) host=$(DB_HOST) port=$(DB_PORT) sslmode=disable

CUR_DIR=$(shell pwd)
MIGRATIONS_FOLDER=$(CUR_DIR)/scripts/migrations

.PHONY: help
## prints help about all targets
help:
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@awk '                                \
		BEGIN { comment=""; }             \
		/^\s*##/ {                         \
		    comment = substr($$0, index($$0,$$2)); next; \
		}                                  \
		/^[a-zA-Z0-9_-]+:/ {               \
		    target = $$1;                  \
		    sub(":", "", target);          \
		    if (comment != "") {           \
		        printf "  %-17s %s\n", target, comment; \
		        comment="";                \
		    }                              \
		}' $(MAKEFILE_LIST)
	@echo ""

.PHONY: tidy
## runs go mod tidy
tidy:
	go mod tidy

.PHONY: fmt
## runs go fmt
fmt:
	go fmt ./...

.PHONY: lint
## runs linter
lint: fmt
	golangci-lint run -c .golangci.yaml

.PHONY: vet
## runs go vet
vet: fmt
	go vet ./...

.PHONY: gen-protoc
## runs protoc in order to generate proto files
gen-protoc:
	protoc -I api api/sso.proto --go_out=./internal/pkg/pb --go_opt=paths=source_relative --go-grpc_out=./internal/pkg/pb --go-grpc_opt=paths=source_relative

.PHONY: run
run: fmt
	go run cmd/sso/main.go

.PHONY: migration-create
## creates migration with first param as name
migration-create:
	goose -dir "$(MIGRATIONS_FOLDER)" create $(name) sql

.PHONY: migration-up
## applies latest migration
migration-up:
	goose -dir "$(MIGRATIONS_FOLDER)" postgres "$(DB_SETUP)" up

.PHONY: migration-down
## rolls back latest migration
migration-down:
	goose -dir "$(MIGRATIONS_FOLDER)" postgres "$(DB_SETUP)" down

.PHONY: migration-status
## checks migration status
migration-status:
	goose -dir "$(MIGRATIONS_FOLDER)" postgres "$(DB_SETUP)" status
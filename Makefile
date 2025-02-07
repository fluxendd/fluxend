# Variables
include .env
export
DATABASE_CONNECTION="user=${DATABASE_USER} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} host=${DATABASE_HOST} sslmode=${DATABASE_SSL_MODE}"

# Include other files
-include Makefiles/*.mk

# Commands
.PHONY: help run build migrate

help: ## Shows this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\-\.]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the project
	@go run cmd/main.go

build: ## Build the project
	@go build -o bin/main cmd/main.go

lint: ## Run linter
	@golangci-lint run

lint-fix: ## Run linter and fix
	@golangci-lint run --fix

login-fluxton: ## Login to fluxton container
	@docker exec -it fluxton /bin/sh

login-db: ## Login to database container
	@docker exec -it fluxton_db /bin/bash

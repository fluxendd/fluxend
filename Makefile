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
	@go run cmd/*.go

list_routes: ## Show the routes
	@go run cmd/*.go -cmd=routes

build: ## Build the project
	@go build -o bin/fluxton cmd/*.go

lint: ## Run linter
	@golangci-lint run

lint-fix: ## Run linter and fix
	@golangci-lint run --fix

up: ## Start the project
	@docker-compose up

down: ## Stop the project
	@docker-compose down

login-fluxton: ## Login to fluxton container
	@docker exec -it fluxton /bin/sh

login-db: ## Login to database container
	@docker exec -it fluxton_db /bin/bash

postgrest-list: ## List all postgrest containers
	@docker ps --filter "name=postgrest_" --format "table {{.ID}}\t{{.Names}}\t{{.Ports}}\t{{.CreatedAt}}\t{{.Status}}"

postgrest-destroy: ## Destroy all postgrest containers
	@docker rm -f $(shell docker ps -a -q --filter "name=postgrest_")

docs: ## Generate docs
	@swag init --dir cmd,controllers,requests,resources,responses,types --output cmd/docs
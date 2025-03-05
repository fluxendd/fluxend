# Variables
-include .env
export
DATABASE_CONNECTION="user=${DATABASE_USER} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} host=${DATABASE_HOST} sslmode=${DATABASE_SSL_MODE}"

# Include other files
include makefiles/*.mk

# Commands
.PHONY: help serve routes build setup up down login-app login-db postgrest-list postgrest-destroy docs

help: ## Shows this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\-\.]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Setup the project
	cp .env.example .env
	make up

serve: ## Run the project in development mode
	@go run main.go server

routes: ## Show all the available routes
	@go run main.go routes

build: ## Build the project
	@go build -o bin/fluxton main.go

up: ## Start the project
	@docker-compose up -d

down: ## Stop the project
	@docker-compose down

login.app: ## Login to fluxton container
	@docker exec -it fluxton_app /bin/sh

login.db: ## Login to database container
	@docker exec -it fluxton_db /bin/bash

pgr.list: ## List all postgrest containers
	@docker ps --filter "name=postgrest_" --format "table {{.ID}}\t{{.Names}}\t{{.Ports}}\t{{.CreatedAt}}\t{{.Status}}"

pgr.destroy: ## Destroy all postgrest containers
	@docker rm -f $(shell docker ps -a -q --filter "name=postgrest_")

docs: ## Generate docs
	@swag init --dir cmd,controllers,requests,resources,responses,types --output cmd/docs
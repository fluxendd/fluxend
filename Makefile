# Variables
-include .env
export
DATABASE_CONNECTION="user=${DATABASE_USER} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} host=${DATABASE_HOST} sslmode=${DATABASE_SSL_MODE}"

# Include other files
include scripts/makefiles/*.mk

help: ## Shows this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\-\.]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Setup the project
	cp .env.example .env
	make up

build: ## Build the project with all containers
	@make down
	@docker-compose up -d --build

build.app: ## Rebuild the app container only
	@docker-compose stop $${APP_CONTAINER_NAME}
	@docker-compose rm -f $${APP_CONTAINER_NAME}
	@docker-compose build $${APP_CONTAINER_NAME}
	@docker-compose up -d $${APP_CONTAINER_NAME}

build.binary: ## Build the binary for the app
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o bin/fluxend cmd/main.go

up: ## Start the project
	@make down
	@docker-compose up -d

down: ## Stop the project
	@docker-compose down

login.app: ## Login to fluxend container
	@docker exec -it $${APP_CONTAINER_NAME} /bin/sh

login.frontend: ## Login to fluxend container
	@docker exec -it $${FRONTEND_CONTAINER_NAME} /bin/sh

login.db: ## Login to database container
	@docker exec -it $${DATABASE_CONTAINER_NAME} /bin/bash

inspect.labels:
	@docker ps --format '{{.Names}}' | \
	while read -r name; do \
		echo "Container: $$name"; \
		docker inspect "$$name" --format '{{json .Config.Labels}}' | jq .; \
		echo ""; \
	done

pgr.list: ## List all postgrest containers
	@docker ps --filter "name=postgrest_" --format "table {{.ID}}\t{{.Names}}\t{{.Ports}}\t{{.CreatedAt}}\t{{.Status}}"

pgr.destroy: ## Destroy all postgrest containers
	@docker rm -f $(shell docker ps -a -q --filter "name=postgrest_")

docs.generate: ## Generate docs
	swag init --dir cmd,internal --output docs
# Variables
-include .env
export
DATABASE_CONNECTION="user=${DATABASE_USER} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} host=${DATABASE_HOST} sslmode=${DATABASE_SSL_MODE}"

# Include other files
include scripts/makefiles/*.mk

ifeq ($(URL_SCHEME),https)
    COMPOSE_FILES = -f docker-compose.yml -f docker-compose.ssl.yml
else
    COMPOSE_FILES = -f docker-compose.yml
endif

DOCKER_COMPOSE = docker-compose $(COMPOSE_FILES)

help: ## Shows this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\-\.]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Interactive setup for new users
	@echo "🚀 Setting up Fluxend..."
	@make check-deps
	@make setup-env
	@make build
	@make migration.up
	@make seed.fresh
	@make verify-setup

check-deps: ## Check if required dependencies are installed
	@echo "Checking dependencies..."
	@command -v docker >/dev/null 2>&1 || { echo "❌ Docker is required but not installed."; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "❌ docker-compose is required but not installed."; exit 1; }
	@echo "✅ Dependencies check passed"

setup-env:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		cp web/.env.example web/.env; \
		echo "📝 Created .env file from template"; \
		echo "⚠️  Please edit .env with your configuration before continuing"; \
		read -p "Press enter when you've configured .env..."; \
	fi

verify-setup:
	@scripts/verify.sh


build: ## Build the project with all containers
	@make down
	@$(DOCKER_COMPOSE) up -d --build

build.api: ## Rebuild the api container only
	@$(DOCKER_COMPOSE) stop $${API_CONTAINER_NAME}
	@$(DOCKER_COMPOSE) rm -f $${API_CONTAINER_NAME}
	@$(DOCKER_COMPOSE) build $${API_CONTAINER_NAME}
	@$(DOCKER_COMPOSE) up -d $${API_CONTAINER_NAME}

build.binary: ## Build the binary for the app
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o bin/fluxend cmd/main.go

restart.all: ## Restart all containers
	@make down
	@$(DOCKER_COMPOSE) up -d

restart.api: ## Restart API container only
	@$(DOCKER_COMPOSE) stop $${API_CONTAINER_NAME}
	@$(DOCKER_COMPOSE) rm -f $${API_CONTAINER_NAME}
	@$(DOCKER_COMPOSE) up -d $${API_CONTAINER_NAME}

restart.frontend: ## Restart frontend container only
	@$(DOCKER_COMPOSE) stop $${FRONTEND_CONTAINER_NAME}
	@$(DOCKER_COMPOSE) rm -f $${FRONTEND_CONTAINER_NAME}
	@$(DOCKER_COMPOSE) up -d $${FRONTEND_CONTAINER_NAME}

restart.app: ## Restart both API and frontend containers only
	@make restart.api
	@make restart.frontend

down: ## Stop the project
	@$(DOCKER_COMPOSE) down

login.api: ## Login to API container
	@docker exec -it $${API_CONTAINER_NAME} /bin/sh

login.frontend: ## Login to frontend container
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
	@docker ps -a -q --filter "name=postgrest_" | xargs -r docker rm -f

pgr.restart: ## Restart all postgrest containers
	@make udb.restart

docs.generate: ## Generate docs
	swag init --dir cmd,internal --output docs

docs.toOpenAPI: ## Convert docs to OpenAPI format
	swagger2openapi docs/swagger.json -o docs/openapi.json

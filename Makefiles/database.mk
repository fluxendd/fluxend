.PHONY: migrate-create migrate-up migrate-down migrate-status migrate-reset migrate-redo migrate-fresh seed seed-fresh

migrate-create: ## Create a new database migration
	@read -p "Enter migration name: " name; \
	goose -dir migrations create $$name sql

migrate-up: ## Run database migrations
	goose -dir migrations postgres ${DATABASE_CONNECTION} up

migrate-down: ## Rollback database migrations
	goose -dir migrations postgres ${DATABASE_CONNECTION} down

migrate-status: ## Show the status of the database migrations
	goose -dir migrations postgres ${DATABASE_CONNECTION} status

migrate-reset: ## Rollback all migrations and run them again
	goose -dir migrations postgres ${DATABASE_CONNECTION} reset

migrate-redo: ## Rollback the last migration and run it again
	goose -dir migrations postgres ${DATABASE_CONNECTION} redo

migrate-fresh: ## Rollback all migrations and run them again
	make migrate-reset
	make migrate-up

seed: ## Seed the database
	@go run seed/main.go

seed-fresh: ## Seed the database with fresh data
	make migrate-fresh
	make seed
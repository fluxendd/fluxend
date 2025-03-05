.PHONY: migrate-create migrate-up migrate-down migrate-status migrate-reset migrate-redo migrate-fresh seed seed-fresh

drop-user-databases: ## Drop all user-created databases
	@docker exec -i fluxton_db psql -U ${DATABASE_USER} -d ${DATABASE_NAME} -t -c "SELECT datname FROM pg_database WHERE datname LIKE 'udb_%';" | sed 's/^[ \t]*//' | while read dbname; do \
		if [ ! -z "$$dbname" ]; then \
			echo "Dropping database: $$dbname"; \
			docker exec -i fluxton_db psql -U ${DATABASE_USER} -d ${DATABASE_NAME} -c "DROP DATABASE IF EXISTS \"$$dbname\""; \
		fi \
	done

migrate-create: ## Create a new database migration
	@read -p "Enter migration name: " name; \
	goose -dir database/migrations create $$name sql

migrate-up: ## Run database migrations
	goose -dir database/migrations postgres ${DATABASE_CONNECTION} up

migrate-down: ## Rollback database migrations
	goose -dir database/migrations postgres ${DATABASE_CONNECTION} down

migrate-status: ## Show the status of the database migrations
	goose -dir database/migrations postgres ${DATABASE_CONNECTION} status

migrate-reset: ## Rollback all migrations and run them again
	goose -dir database/migrations postgres ${DATABASE_CONNECTION} reset

migrate-redo: ## Rollback the last migration and run it again
	goose -dir database/migrations postgres ${DATABASE_CONNECTION} redo

migrate-fresh: ## Rollback all migrations and run them again
	make drop-user-databases
	make migrate-reset
	make migrate-up

seed: ## Seed the database
	@go run cmd/*.go -cmd=seed

seed-fresh: ## Seed the database with fresh data
	make migrate-fresh
	make seed




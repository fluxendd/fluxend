drop.user.dbs: ## Drop all user-created databases
	@docker exec -i fluxend_db psql -U ${DATABASE_USER} -d ${DATABASE_NAME} -t -c "SELECT datname FROM pg_database WHERE datname LIKE 'udb%';" | sed 's/^[ \t]*//' | while read dbname; do \
		if [ ! -z "$$dbname" ]; then \
			echo "Dropping database: $$dbname"; \
			docker exec -i fluxend_db psql -U ${DATABASE_USER} -d ${DATABASE_NAME} -c "DROP DATABASE IF EXISTS \"$$dbname\" with (force)"; \
		fi \
	done

migration.create: ## Create a new database migration
	@read -p "Enter migration name: " name; \
	goose -dir internal/database/migrations create $$name sql

migration.up: ## Run database migrations
	goose -dir internal/database/migrations postgres ${DATABASE_CONNECTION} up

migration.down: ## Rollback database migrations
	goose -dir internal/database/migrations postgres ${DATABASE_CONNECTION} down

migration.status: ## Show the status of the database migrations
	goose -dir internal/database/migrations postgres ${DATABASE_CONNECTION} status

migration.reset: ## Rollback all migrations and run them again
	goose -dir internal/database/migrations postgres ${DATABASE_CONNECTION} reset

migration.redo: ## Rollback the last migration and run it again
	goose -dir internal/database/migrations postgres ${DATABASE_CONNECTION} redo

migration.fresh: ## Rollback all migrations and run them again
	make drop.user.dbs
	make migration.reset
	make migration.up

seed.fresh: ## Seed the database with fresh data
	make migration.fresh
	make seed




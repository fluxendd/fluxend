# These are core command also available in binary CLI

serve: ## Run the project in development mode
	@go run main.go server

routes.list: ## Show all the available routes
	@go run main.go routes

seed: ## Seed the database
	@go run main.go seed

about: ## Show the project information
	@go run main.go about

optimize: ## Optimize the project
	@go run main.go optimize

udb.stats: ## Show the database stats
	@go run main.go udb.stats

udb.restart: ## Restart the database
	@go run main.go udb.restart
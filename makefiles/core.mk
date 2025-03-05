# These are core command also available in binary CLI

serve: ## Run the project in development mode
	@go run main.go server

routes: ## Show all the available routes
	@go run main.go routes

seed: ## Seed the database
	@go run main.go seed
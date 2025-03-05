.PHONY: lint lint.fix

lint: ## Run linter
	@golangci-lint run

lint.fix: ## Run linter and fix
	@golangci-lint run --fix


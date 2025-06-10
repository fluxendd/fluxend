lint: ## Run linter
	@golangci-lint run

lint.fix: ## Run linter and fix
	@golangci-lint run --fix

test: ## Run tests
	@go test -v ./... -coverprofile=coverage.out

cloc: ## Count lines of code
	cloc . \
      --exclude-dir=node_modules,vendor,pkg,dist,build,out,.next,.turbo,.cache,.git \
      --exclude-ext=xml,XML,json,md,SVG,svg


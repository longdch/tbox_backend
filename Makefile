
test: ## Run go test for whole project
	@go test -v ./...

build: ## Build containers
	@docker-compose build

run: ## Run containers
	@docker-compose up

stop: ## Stop containers
	@docker-compose down

lint: ## Run linter
	@golangci-lint run ./...

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

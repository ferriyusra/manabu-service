GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

## Live reload:
watch-prepare: ## Install the tools required for the watch command
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh

watch: ## Run the service with hot reload
	@echo "$(CYAN)Generating Swagger documentation...$(RESET)"
	@go generate
	@echo "$(GREEN)Swagger docs generated!$(RESET)"
	bin/air

run: ## Run the service (generate swagger + run)
	@echo "$(CYAN)Generating Swagger documentation...$(RESET)"
	@go generate
	@echo "$(GREEN)Swagger docs generated!$(RESET)"
	@echo "$(CYAN)Starting application...$(RESET)"
	go run main.go serve

## Build:
build: ## Build the service
	go generate
	go build -o manabu-service

## Swagger:
swagger: ## Generate Swagger documentation
	swag init -g cmd/main.go -o docs
	@echo "$(GREEN)Swagger documentation generated successfully!$(RESET)"
	@echo "$(CYAN)Access Swagger UI at: http://localhost:8001/swagger/index.html$(RESET)"

## Docker:
docker-compose: ## Start the service in docker
	docker-compose up -d --build --force-recreate

docker-build: ## Build the Docker image with a specified tag
	@echo "$(CYAN)Building Docker image...$(RESET)"
	@if [ -z "$(tag)" ]; then \
		echo "$(YELLOW)Error: Please specify the 'tag' parameter, e.g., make docker-build tag=1.0.0$(RESET)"; \
		exit 1; \
	fi
	docker build --platform linux/amd64 -t sikoding20/manabu-service:$(tag) .
	@echo "$(GREEN)Docker image built with tag '$(tag)'$(RESET)"

docker-push: ## Build the Docker image with a specified tag
	@echo "$(CYAN)Building Docker image...$(RESET)"
	@if [ -z "$(tag)" ]; then \
		echo "$(YELLOW)Error: Please specify the 'tag' parameter, e.g., make docker-push tag=1.0.0$(RESET)"; \
		exit 1; \
	fi
	docker push sikoding20/manabu-service:$(tag)
	@echo "$(GREEN)Docker image built with tag '$(tag)'$(RESET)"

## Testing:
test: ## Run all tests
	@echo "$(CYAN)Running all tests...$(RESET)"
	go test ./... -v
	@echo "$(GREEN)All tests passed!$(RESET)"

test-coverage: ## Run tests with coverage report
	@echo "$(CYAN)Running tests with coverage...$(RESET)"
	go test ./controllers/user/... ./services/user/... ./repositories/user/... -coverprofile=coverage.out
	@echo "$(GREEN)Coverage report generated!$(RESET)"
	@echo "$(CYAN)Coverage summary:$(RESET)"
	go tool cover -func=coverage.out
	@echo "$(YELLOW)To view HTML report: make test-coverage-html$(RESET)"

test-coverage-html: ## Generate HTML coverage report
	@echo "$(CYAN)Generating HTML coverage report...$(RESET)"
	go test ./controllers/user/... ./services/user/... ./repositories/user/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)HTML coverage report generated: coverage.html$(RESET)"

test-user: ## Run User API tests only
	@echo "$(CYAN)Running User API tests...$(RESET)"
	go test ./controllers/user/... ./services/user/... ./repositories/user/... -v -cover

test-controller: ## Run controller tests only
	@echo "$(CYAN)Running controller tests...$(RESET)"
	go test ./controllers/... -v -cover

test-service: ## Run service tests only
	@echo "$(CYAN)Running service tests...$(RESET)"
	go test ./services/... -v -cover

test-repository: ## Run repository tests only
	@echo "$(CYAN)Running repository tests...$(RESET)"
	go test ./repositories/... -v -cover

mock-generate: ## Generate mocks using mockery
	@echo "$(CYAN)Generating mocks...$(RESET)"
	mockery
	@echo "$(GREEN)Mocks generated successfully!$(RESET)"

mock-install: ## Install mockery tool
	@echo "$(CYAN)Installing mockery...$(RESET)"
	go install github.com/vektra/mockery/v2@latest
	@echo "$(GREEN)Mockery installed successfully!$(RESET)"
	@echo "$(YELLOW)Run 'make mock-generate' to generate mocks$(RESET)"

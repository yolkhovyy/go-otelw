PROJECT_NAME ?= $(notdir $(shell pwd))

DOCO = docker compose

export COMPOSE_DOCKER_CLI_BUILD = 1
export DOCKER_BUILDKIT = 1
export COMPOSE_BAKE = true

export DEPENDENCIES = otel-collector jaeger prometheus
export SERVICES = $(filter-out ${DEPENDENCIES}, $(shell ${DOCO} config --services))

## Docker compose:

.PHONY: doco-up-dependencies
doco-up-dependencies: ## Start project dependency containers
	@echo "🐳 Starting project dependencies"
	@${DOCO} up --detach --quiet-pull ${DEPENDENCIES}

.PHONY: doco-build
doco-build: lint ## Build Docker images for the project
	@echo "🐳 Building Docker images"
	@${DOCO} build ${SERVICES}

.PHONY: doco-build-up
doco-build-up: lint doco-up-dependencies ## Build Docker images and start the services
	@echo "🐳 Building and starting services"
	@${DOCO} up --build --detach --force-recreate --remove-orphans ${SERVICES}

.PHONY: doco-up
doco-up: doco-up-dependencies ## Start project services in Docker containers
	@echo "🐳 Starting project services"
	@${DOCO} up --detach --remove-orphans ${SERVICES}

.PHONY: doco-stop
doco-stop: ## Stop running Docker containers
	@echo "🐳 Stopping running Docker containers"
	@${DOCO} stop

.PHONY: doco-down
doco-down: ## Stop and remove Docker containers and associated network
	@echo "🗑 Stopping and removing Docker containers and associated network"
	@${DOCO} down


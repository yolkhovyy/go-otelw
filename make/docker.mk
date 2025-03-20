PROJECT_NAME ?= $(notdir $(shell pwd))

export COMPOSE_DOCKER_CLI_BUILD = 1
export DOCKER_BUILDKIT = 1
export COMPOSE_BAKE = true

DOCO = docker compose -f docker-compose.yml
export DEPENDENCIES = otel-collector

ifdef NR		#----- newrelic
	export OTEL_COLLECTOR_CONFIG = config-newrelic.yml
else ifdef UPT	#----- uptrace
	DOCO := $(DOCO) -f docker-compose.uptrace.yml
	DEPENDENCIES := $(DEPENDENCIES) postgres clickhouse uptrace
	export OTEL_COLLECTOR_CONFIG = config-uptrace.yml
else			#----- jaeger, prometheus, grafana loki
	DOCO := $(DOCO) -f docker-compose.jplpg.yml
	DEPENDENCIES := $(DEPENDENCIES) jaeger prometheus loki promtail grafana
	export OTEL_COLLECTOR_CONFIG = config-jplpg.yml
endif

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

.PHONY: doco-watch
doco-watch: ## Watch running docker containers
	@watch ${DOCO} ps

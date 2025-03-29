## Docker compose:	

.PHONY: doco-up-dependencies
doco-up-dependencies: ## Start project dependency containers
	@echo "🐳 Starting project dependencies"
	@${DOCO} up --detach --quiet-pull ${DEPENDENCIES}

.PHONY: doco-build
doco-build: lint ## Build Docker images for the project
	@echo "🐳 Building Docker images"
	@echo "APP_VERSION: ${APP_VERSION}"
	@${DOCO} build ${SERVICES}

.PHONY: doco-build-up
doco-build-up: lint doco-up-dependencies ## Build Docker images and start the services
	@echo "🐳 Building and starting services"
	@echo "APP_VERSION: ${APP_VERSION}"
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

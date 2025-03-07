## Lint:

.PHONY: lint-go-check
lint-go-check:
ifeq (, $(shell command -v golangci-lint 2> /dev/null))
	@echo "❌ golangci-lint not installed, run 'make install'"
	@exit 1
endif

.PHONY: lint-go
lint-go: lint-go-check ## lint go files
	@echo "🔍 Linting Go files"
	@golangci-lint run ./...

.PHONY: lint-yaml-check
lint-yaml-check:
ifeq (, $(shell command -v yamllint 2> /dev/null))
	@echo "❌ yamllint not installed, run 'sudo apt install yamllint'"
	@exit 1
endif

.PHONY: lint-yaml
lint-yaml: lint-yaml-check  ## lint yaml files
	@echo "🔍 Linting Yaml files"
	@yamllint -f parsable .

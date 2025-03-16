## Tools:

.PHONY: install-tools
install-tools: ## install tools
	@echo "🛠️  Installing tools"
	@go install github.com/boumenot/gocover-cobertura@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/jstemmer/go-junit-report@latest
	@go install github.com/powerman/gotmpl@latest
	@go install github.com/thomaspoignant/yamllint-checkstyle@latest
	@go install github.com/wadey/gocovmerge@latest
	@go install mvdan.cc/gofumpt@latest

.PHONY: install-git-hooks
install-git-hooks: ## install git hooks
	@echo "🛠️  Installing git hooks"
	@if [ -d ".git/hooks" ]; then \
		cp make/scripts/pre-commit.sh .git/hooks/pre-commit; \
		cp make/scripts/pre-push.sh .git/hooks/pre-push; \
	fi

.PHONY: install-env ## Install .env
install-env:
	@echo "🛠️  Installing .env"
	@cp .env.local .env
	@if [ -f .env.newrelic ]; then \
		cat .env.newrelic >> .env; \
	fi

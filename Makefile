
.PHONY: all
all: clean install lint test coverage ## Clean, install, lint, test, coverage

## Help:
.PHONY: help
help: display-help ## display this help

## Install:
.PHONY: install
install: install-tools install-git-hooks install-env ## install this project

## Lint:
.PHONY: lint
lint: lint-yaml lint-go  ## run all linters

## Test:
.PHONY: test
test: unit-test ## run all tests

## Clean:
.PHONY: clean
clean: remove-generated ## clean project

-include make/*.mk

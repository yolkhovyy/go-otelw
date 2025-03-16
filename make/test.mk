## Tests:

.PHONY: unit-test
unit-test: ## run unit tests
	@echo "⚙ Running unit tests"
	@mkdir -p test_results coverage
	@go test -count=1 -v -coverpkg=./... -coverprofile=coverage/unit-test.cov ./... | tee test_results/unit-test.out
	@go-junit-report -set-exit-code < test_results/unit-test.out > test_results/unit-test.xml

.PHONY: remove-generated
remove-generated: ## Remove generated folders and files
	@echo "🗑 Removing generated folders and files"
	@rm -rf coverage/ test_results/

## Coverage:

.PHONY: gocovmerge-check
gocovmerge-check:
ifeq (, $(shell command -v gocovmerge 2> /dev/null))
	@echo "❌ gocovmerge not installed, run 'make install'"
	@exit 1
endif

.PHONY: gocover-cobertura-check
gocover-cobertura-check:
ifeq (, $(shell command -v gocover-cobertura 2> /dev/null))
	@echo "❌ gocover-cobertura not installed, run 'make install'"
	@exit 1
endif

.PHONY: coverage
coverage: gocovmerge-check gocover-cobertura-check ## Make coverage report
	@if [ ! -d "coverage/" ]; then \
		${MAKE} test; \
	fi
	@for file in coverage/*.cov; do \
		cp $$file $$file.tmp; \
		cat $$file.tmp | grep -v -e "test" > $$file; \
	done
	@gocovmerge coverage/*.cov > coverage/total.cov
	@go tool cover -html=coverage/total.cov -o coverage/total.html
	@go tool cover -func coverage/total.cov > coverage/total.txt
	@gocover-cobertura < coverage/total.cov > coverage/total.xml

.PHONY: help
help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: new
new: ## create memo
	go run . new

.PHONY: new-truncate
new-truncate: ## Create memo after truncating the file if it exists
	go run . new -t

.PHONY: weekly
weekly: ## Create weekly report
	go run . weekly

.PHONY: tips
tips: ## Generate tips index
	go run . tips

.PHONY: test
test: ## Run tests
	go test ./... -cover

.PHONY: cyclo
cyclo: ## Run gocyclo
	@command -v gocyclo >/dev/null 2>&1 || { echo >&2 "gocyclo is required but it's not installed.  Aborting."; exit 1; }
	@echo "Running gocyclo"
	gocyclo -top 5 .
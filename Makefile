.PHONY: help
help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: create
create: ## Create today's memo
	go run . create

.PHONY: create-truncate
create-truncate: ## Create today's memo after truncating the file if it exists
	go run . create --truncate

.PHONY: weekly
weekly: ## Create weekly report
	go run . weekly

.PHONY: tips
tips: ## Generate tips index
	go run . tips

.PHONY: test
test: ## Run tests
	go test ./... -cover
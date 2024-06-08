.PHONY: help
help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: create
create: ## Create today's memo
	go run .

.PHONY: create-truncate
create-truncate: ## Create today's memo while truncating it if the file already exists
	go run . --truncate

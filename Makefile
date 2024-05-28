foo := 'foo'


.PHONY: help
help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: create
create: ## truncate and create today's memo
	go run . -truncate

aaa:  ## aaa
bbb:  ## bbb
ccc:  ## ccc

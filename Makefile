.PHONY: build
build: ## Do a caching build of the website
	cd builder; go run builder.go --src ../src --out ../out

.PHONY: build-watch
build-watch: ## Do a caching build of the website every time a file is changed
	./tool/serve-watch.sh

.PHONY: help
help: ## Print this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: prod-build
prod-build: ## Do a full build of the website
	cd builder; go run builder.go --src ../src --out ../out --nocache

.PHONY: serve
serve: ## Serve the website
	cd out; python -m http.server 8080
.PHONY: build
build: CMD:="make internal-cache-build"
build: internal-run-in-build-image ## Do a caching build of the website

.PHONY: build-prod
build-prod: CMD:="make internal-prod-build"
build-prod: internal-run-in-build-image ## Do a full non-caching build of the website

.PHONY: build-watch
build-watch: CMD:="./tool/serve-watch.sh"
build-watch: internal-run-in-build-image ## Do a caching build of the website every time a file is changed

.PHONY: help
help: ## Print this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: serve
serve: internal-build-image ## Serves the website to port 8080
	docker run --rm -p 8080:8080 -v ./out:/workdir naberhausj.com-builder sh -c "python -m http.server 8080"

# --------------------------------------
# Targets intended for internal use only
# --------------------------------------

internal-build-image: Dockerfile
	docker build -t naberhausj.com-builder .

.PHONY: internal-cache-build
internal-cache-build:
	cd builder; go run builder.go --src ../src --out ../out

.PHONY: internal-prod-build
internal-prod-build:m
	cd builder; go run builder.go --src ../src --out ../out --nocache

.PHONY: internal-run-in-build-image
internal-run-in-build-image: internal-build-image
	docker run --rm -v .:/workdir naberhausj.com-builder sh -c $(CMD)

.PHONY: serve
serve:
	cd out; python -m http.server 8080

.PHONY: build
build:
	cd builder; go run builder.go --src ../src --out ../out --nocache

.PHONY: dev-build
dev-build:
	cd builder; go run builder.go --src ../src --out ../out

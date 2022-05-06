.PHONY: serve
serve:
	cd out; python3 -m http.server 8080

.PHONY: build
build:
	cd builder; go run builder.go

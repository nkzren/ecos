.PHONY: setup
setup:
	@cp config.yaml.sample config.yaml

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build: fmt
	go build

.PHONY: run
run: build
	@./ecoscheduler

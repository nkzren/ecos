IMAGE_TAG_BASE ?= nkzren/ecoscheduler
VERSION ?= latest

.PHONY: setup
setup:
	@cp config.yaml.sample config.yaml

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build: fmt
	go build

.PHONY: docker-build
docker-build:
	docker build -t ${IMAGE_TAG_BASE}:${VERSION} .

.PHONY: docker-push
docker-push:
	docker push ${IMAGE_TAG_BASE}:${VERSION}

.PHONY: run
run: build
	@./ecoscheduler

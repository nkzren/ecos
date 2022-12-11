IMAGE_TAG_BASE ?= nkzren/ecos
VERSION ?= latest

.PHONY: release
release: docker-build docker-push

.PHONY: setup
setup:
	@echo "Copy config file"
	@cp config.yaml.sample config.yaml

.PHONY: cluster-setup
cluster-setup:
	@echo "Setting up cluster permissions"
	@kubectl apply -f kube/samples/rbac.yaml
	@kubectl create clusterrolebinding ecos --clusterrole=ecos --serviceaccount=default:default 2> /dev/null; true

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build: fmt
	go build

.PHONY: docker-build
docker-build: build
	docker build -t ${IMAGE_TAG_BASE}:${VERSION} .

.PHONY: docker-push
docker-push:
	docker push ${IMAGE_TAG_BASE}:${VERSION}

.PHONY: clean
clean:
	go clean

.PHONY: run
run: build
	@./ecos

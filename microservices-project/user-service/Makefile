IMAGE ?= ghcr.io/example/user-service
TAG ?= latest

build:
	docker build -f deployments/docker/Dockerfile -t $(IMAGE):$(TAG) .

push:
	docker push $(IMAGE):$(TAG)

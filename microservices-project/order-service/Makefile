IMAGE ?= ghcr.io/example/order-service
TAG ?= latest

build:
	docker build -f deployments/docker/Dockerfile -t $(IMAGE):$(TAG) .

push:
	docker push $(IMAGE):$(TAG)

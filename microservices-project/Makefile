REGISTRY ?= ghcr.io/example
TAG ?= latest

.PHONY: docker-build docker-push k8s-deploy k8s-delete

docker-build:
	docker build -t $(REGISTRY)/user-service:$(TAG) -f user-service/deployments/docker/Dockerfile user-service
	docker build -t $(REGISTRY)/order-service:$(TAG) -f order-service/deployments/docker/Dockerfile order-service

docker-push:
	docker push $(REGISTRY)/user-service:$(TAG)
	docker push $(REGISTRY)/order-service:$(TAG)

k8s-deploy:
	kubectl apply -f k8s/

k8s-delete:
	kubectl delete -f k8s/ --ignore-not-found=true

# Microservices Project (Go + Kubernetes)

Production-oriented sample with two Go microservices using Clean Architecture:

- `user-service`
- `order-service` (validates users via HTTP call to `user-service` before creating an order)

## Prerequisites

Only these tools are required:

- Docker
- Kubernetes cluster (e.g. kind, k3d, minikube, EKS, GKE, AKS)
- kubectl

No local Go toolchain and no local PostgreSQL installation are required.

## Project layout

```text
microservices-project/
├── user-service/
├── order-service/
├── k8s/
│   ├── namespace.yaml
│   ├── postgres.yaml
│   ├── user-service.yaml
│   ├── order-service.yaml
│   └── ingress.yaml
├── Makefile
└── README.md
```

## Build and push images

Update image names in Kubernetes manifests (`ghcr.io/example/*`) or set your own registry in Makefile.

```bash
make docker-build REGISTRY=ghcr.io/<your-org> TAG=v1
make docker-push  REGISTRY=ghcr.io/<your-org> TAG=v1
```

Then update `k8s/user-service.yaml` and `k8s/order-service.yaml` image tags accordingly.

## Deploy to Kubernetes

```bash
kubectl apply -f k8s/
```

Or:

```bash
make k8s-deploy
```

Delete all resources:

```bash
make k8s-delete
```

## Port-forward and test with curl

### 1) user-service

```bash
kubectl -n microservices port-forward svc/user-service 8080:8080
```

Create user:

```bash
curl -i -X POST http://localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com"}'
```

Get user:

```bash
curl -i http://localhost:8080/users/1
```

Health and metrics:

```bash
curl -i http://localhost:8080/healthz
curl -i http://localhost:8080/readyz
curl -i http://localhost:8080/metrics
```

### 2) order-service

```bash
kubectl -n microservices port-forward svc/order-service 8081:8081
```

Create order for existing user:

```bash
curl -i -X POST http://localhost:8081/orders \
  -H 'Content-Type: application/json' \
  -d '{"user_id":1,"amount":59.99}'
```

Get order:

```bash
curl -i http://localhost:8081/orders/1
```

Health and metrics:

```bash
curl -i http://localhost:8081/healthz
curl -i http://localhost:8081/readyz
curl -i http://localhost:8081/metrics
```

## Notes

- PostgreSQL is deployed as a StatefulSet.
- Separate databases are created: `userdb` and `orderdb`.
- Migrations are included under each service `migrations/` directory.
- Dockerfiles are multi-stage and use distroless non-root runtime images.

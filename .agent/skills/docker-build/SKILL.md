---
name: docker-build
description: Build the multi-stage Docker image for all services.
---

# Docker Build Skill

Use this skill to build the unified Docker image containing the API, Publisher, Worker, and Scaling Dashboard.

## Usage

### 1. Build the image locally
// turbo
```powershell
docker build -t email-system:latest .
```

### 2. Build for Minikube (Directly into Minikube's Docker registry)
// turbo
```powershell
minikube docker-env | Invoke-Expression; docker build -t email-system:latest .
```

## When to use
- After changing any Go code in `cmd/`, `internal/`, or `pkg/`.
- After updating the `Dockerfile`.
- Before deploying to Kubernetes.

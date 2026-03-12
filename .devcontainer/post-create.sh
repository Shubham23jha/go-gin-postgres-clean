#!/bin/bash
set -e

echo "🚀 Starting environment initialization..."

# 1. Install k3d (Lightweight K8s in Docker)
curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | TAG=v5.6.0 bash

# 2. Create Cluster
k3d cluster create email-cluster --port "30080:30080@loadbalancer" --wait

# 3. Install KEDA
echo "🍀 Installing KEDA..."
kubectl apply --server-side -f https://github.com/kedacore/keda/releases/download/v2.16.1/keda-2.16.1.yaml

# 4. Pull/Build Local Infrastructure (Optional but speeds things up)
docker-compose up -d db rabbitmq

# 5. Build and Load local image to k3d
echo "📦 Building project image..."
docker build -t email-system:latest .
k3d image import email-system:latest -c email-cluster

# 6. Apply K8s manifests
echo "☸️ Deploying to K8s..."
kubectl apply -f k8s/

echo "✅ Environment Ready!"
echo "-------------------------------------------------------"
echo "📊 Scaling Dashboard will be available at port 30080"
echo "🚀 API Server will be available at port 8080"
echo "-------------------------------------------------------"

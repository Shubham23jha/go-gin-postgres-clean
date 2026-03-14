#!/bin/bash

# Distributed Email Delivery System - Local Setup Script (BASH)
# This script automates the infrastructure setup, building, and deployment.

echo "🚀 Starting Automated Setup..."

# 1. Install KEDA (If missing)
echo "📦 Checking KEDA..."
kubectl apply --server-side -f https://github.com/kedacore/keda/releases/download/v2.16.1/keda-2.16.1.yaml

# 2. Deploy Infrastructure (Postgres & RabbitMQ)
echo "🏗️  Deploying Infrastructure (Postgres & RabbitMQ)..."
kubectl apply -f k8s/infrastructure.yaml

# 3. Build & Load Application
echo "🛠️  Building Application Image (v4)..."
docker build -t email-system:v4 .

echo "🚚 Loading image into Minikube..."
minikube image load email-system:v4

# 4. Deploy Application
echo "🚀 Deploying Application Services..."
kubectl apply -f k8s/

echo "---------------------------------------------------"
echo "✅ Setup Complete!"
echo "---------------------------------------------------"
echo "📊 Dashboard URL:"
minikube service scaling-dashboard-service --url &

echo "📧 API URL (Dynamic):"
minikube service email-api --url &

echo ""
echo "✨ BETTER SOLUTION (Static Port):"
echo "To use a permanent URL that never changes, run this in a new terminal:"
echo "👉 kubectl port-forward svc/email-api 8080:80"
echo "Then use: http://127.0.0.1:8080"
echo ""

echo "💡 Follow the instructions in README.md to trigger a scaling demo!"

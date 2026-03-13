# Distributed Email Delivery System - Local Setup Script (PowerShell)
# This script automates the infrastructure setup, building, and deployment.

Write-Host "🚀 Starting Automated Setup..." -ForegroundColor Cyan

# 1. Install KEDA (If missing)
Write-Host "📦 Checking KEDA..." -ForegroundColor Yellow
kubectl apply --server-side -f https://github.com/kedacore/keda/releases/download/v2.16.1/keda-2.16.1.yaml

# 2. Deploy Infrastructure (Postgres & RabbitMQ)
Write-Host "🏗️  Deploying Infrastructure (Postgres & RabbitMQ)..." -ForegroundColor Yellow
kubectl apply -f k8s/infrastructure.yaml

# 3. Build & Load Application
Write-Host "🛠️  Building Application Image (v3)..." -ForegroundColor Yellow
docker build -t email-system:v3 .

Write-Host "🚚 Loading image into Minikube..." -ForegroundColor Yellow
minikube image load email-system:v3

# 4. Deploy Application
Write-Host "🚀 Deploying Application Services..." -ForegroundColor Yellow
kubectl apply -f k8s/

Write-Host "---------------------------------------------------" -ForegroundColor Green
Write-Host "✅ Setup Complete!" -ForegroundColor Green
Write-Host "---------------------------------------------------" -ForegroundColor Green

# 5. Show Services
Write-Host "📊 Opening Dashboard..." -ForegroundColor Cyan
Start-Process powershell -ArgumentList "minikube service scaling-dashboard-service --url"

Write-Host "📧 Fetching API URL..." -ForegroundColor Cyan
Start-Process powershell -ArgumentList "minikube service email-api --url"

Write-Host "💡 Follow the instructions in README.md to trigger a scaling demo!" -ForegroundColor Gray

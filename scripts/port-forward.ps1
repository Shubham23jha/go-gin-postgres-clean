# Port Forwarding Script for PowerShell
Write-Host "🔗 Starting Static Port Forwarding for Email API..." -ForegroundColor Cyan
Write-Host "👉 Local URL: http://127.0.0.1:8080" -ForegroundColor Green
Write-Host "Keep this terminal open!" -ForegroundColor Yellow
kubectl port-forward svc/email-api 8080:80

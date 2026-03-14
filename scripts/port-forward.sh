#!/bin/bash
echo "🔗 Starting Static Port Forwarding for Email API..."
echo "👉 Local URL: http://127.0.0.1:8080"
echo "Keep this terminal open!"
kubectl port-forward svc/email-api 8080:80

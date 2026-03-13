---
name: local-setup
description: Run the automated project setup and deployment.
---

# Local Setup Skill

Use this skill to automate the entire setup, including KEDA installation, infrastructure deployment, building images, and rolling out the application.

## Usage

### 1. Run Automated Setup (Bash)
// turbo
```bash
chmod +x setup-local.sh
./setup-local.sh
```

### 2. Run Automated Setup (PowerShell)
// turbo
```powershell
.\setup-local.ps1
```

## When to use
- To start the project from scratch in a local cluster.
- To perform a "clean" redeployment of all services.
- To ensure all dependencies (Postgres, RabbitMQ, KEDA) are properly configured.

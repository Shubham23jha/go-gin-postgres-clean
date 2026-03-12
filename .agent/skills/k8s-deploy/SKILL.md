---
name: k8s-deploy
description: Deploy services to Kubernetes.
---

# Kubernetes Deployment Skill

Use this skill to apply manifests and check the status of the rollout.

## Usage

### 1. Apply all manifests
// turbo
```powershell
kubectl apply -f k8s/
```

### 2. Check deployment status
// turbo
```powershell
kubectl get pods -w
```

### 3. Check KEDA scaling
// turbo
```powershell
kubectl get scaledobject
```

### 4. Install KEDA (If not present)
// turbo
```powershell
kubectl apply --server-side -f https://github.com/kedacore/keda/releases/download/v2.16.1/keda-2.16.1.yaml
```

## When to use
- To deploy the application for the first time.
- After updating files in the `k8s/` directory.
- To troubleshoot scaling or pod issues.

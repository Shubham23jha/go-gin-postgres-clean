---
name: email-resilience
description: Understand and manage the resilience features of the email delivery system.
---

# Digital Post Office Resilience Skill

Use this skill to understand how the system handles failures and how to monitor its reliability.

## Key Features

### 1. Automatic Reconnection
Both the **Outbox Publisher** and **Email Worker Pool** include a robust reconnection mechanism.
- **How it works**: Monitors `NotifyClose` from RabbitMQ.
- **Benefit**: System recovers gracefully after RabbitMQ restarts or network blips without manual intervention.

### 2. Stalled Item Reclaiming
Prevents messages from being stuck in "Picking Up" state if a publisher process dies.
- **Frequency**: Runs every 1 minute.
- **Timeout**: Reclaims items that haven't been updated in 5 minutes.

### 3. Publishing Retries
The Outbox Publisher attempts 3 retries for Each item during a single processing loop.
- **Benefit**: Handles transient RabbitMQ errors immediately.

### 4. Static Port Solution (Connectivity)
Ensures external projects can access the API with a permanent, unchanging URL.
- **Usage**: Run `kubectl port-forward svc/email-api 8080:80`.
- **URL**: Always use `http://127.0.0.1:8080` in external config (e.g., .env of other projects).
- **Benefit**: Prevents `ECONNREFUSED` errors caused by Minikube's dynamic port assignment.

## Verification & Monitoring

### Check Outbox Status
Use `/db-migrate` or manual SQL to check for `FAILED` or `PICKED_UP` items:
```sql
SELECT status, count(*) FROM outbox GROUP BY status;
```

### Observe Logs
Look for these emojis in the service logs to identify resilience actions:
- ♻️ : Reclaiming stalled items.
- ⚠️ : Connection loss or retry attempt.
- ❌ : Critical failure or max retries reached.
- 📥 : Successful reconnection and start.

## 🗄️ Database Inspection
The resilient system uses PostgreSQL inside Kubernetes (Database name: `goDb`).
Use these commands to check your data:

### Check Campaigns
```bash
kubectl exec deployment/postgres -- psql -U postgres -d goDb -c "SELECT id, subject, status, created_at FROM campaigns ORDER BY id DESC;"
```

### Check Outbox Items
```bash
kubectl exec deployment/postgres -- psql -U postgres -d goDb -c "SELECT id, campaign_id, status, recipient, retry_count FROM outbox ORDER BY id DESC;"
```

### Check Email Logs (Final Result)
```bash
kubectl exec deployment/postgres -- psql -U postgres -d goDb -c "SELECT * FROM email_logs ORDER BY id DESC;"
```

## 🔍 Troubleshooting

### "ECONNREFUSED" or Wrong ID in Response
If you see success in the response but the email is not sent, or if you get a connection error:
1.  **Check for Local Process**: You might have a local `go run` or `main.exe` running on port 8080.
2.  **Verify Kubernetes Logs**: Use `kubectl logs deployment/email-publisher` to see if the cluster is actually receiving the messages.
3.  **Kill Local Server**: Stop any local servers and restart the port-forward:
    ```powershell
    kubectl port-forward svc/email-api 8080:80
    ```

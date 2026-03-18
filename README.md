# Digital Post Office 📧🚀

A robust, high-throughput email delivery system built with **Go**, **Postgres**, **RabbitMQ**, and orchestrated with **Kubernetes** + **KEDA**. Featuring a real-time **Scaling Visualization Dashboard**.

---

## 🏗️ Architecture
The system uses the **Transactional Outbox Pattern** to ensure "at-least-once" delivery and decoupled horizontal scaling.

1.  **Campaign API**: Atomic transaction that saves the campaign and its individual email tasks to an `outbox` table.
2.  **Outbox Publisher**: Polling service that reads `PENDING` outbox items and publishes them to RabbitMQ.
3.  **RabbitMQ**: Reliable message broker with built-in **Dead Letter Queues (DLQ)**.
4.  **Email Worker Pool**: A cluster of consumers that deliver emails via SMTP with exponential backoff and idempotency.
5.  **KEDA**: Automatically scales Worker pods from **0 to 10** based on queue length.
6.  **📊 Scaling Dashboard**: A real-time monitoring UI that visualizes queue growth and worker pod scaling.

## 🛡️ Resilience & Reliability
The system is designed to handle infrastructure failures and ensure zero message loss:
- **RabbitMQ Reconnection**: Services (Publisher & Worker) automatically detect connection loss and reconnect with exponential backoff.
- **Transactional Outbox**: Guaranteed "at-least-once" delivery by marking items as `PUBLISHED` only after successful RabbitMQ ack.
- **Stalled Item Reclaimer**: A background task automatically finds and re-processes outbox items stuck in `PICKED_UP` status (e.g., due to service crashes).
- **Worker Retries**: Email workers use exponential backoff for SMTP failures and safely handle message retries via DLX.

## 🚀 Quick Start (Automated)

Run the project with a single command based on your terminal:

**For BASH (Git Bash, WSL, Linux):**
```bash
chmod +x setup-local.sh
./setup-local.sh
```

**For POWERSHELL:**
```powershell
.\setup-local.ps1
```

---

## 🎨 Campaign Trigger UI (New!)
You can now trigger campaigns directly from a beautiful web interface.
1.  **Get the URL**: Run `minikube service email-api --url`
2.  **Open in Browser**: Navigate to that URL (e.g., `http://127.0.0.1:55321/`).
3.  **Trigger**: Enter the number of users and click **"Start Campaign"**.

### 🔗 Accessing from External Projects (Static Port)
If you are calling this API from another project (like a Node.js server), use a static port to avoid URL changes:
1.  Run: `kubectl port-forward svc/email-api 8080:80`
2.  Use URL: `http://127.0.0.1:8080`

---

### 1. Pre-requisites
- **Docker Desktop** (Kubernetes enabled) or **Minikube**.
- **KEDA** installed:
  ```bash
  kubectl apply --server-side -f https://github.com/kedacore/keda/releases/download/v2.16.1/keda-2.16.1.yaml
  ```

### 2. Start Infrastructure
Deploy the database and RabbitMQ directly to the cluster:
```bash
kubectl apply -f k8s/infrastructure.yaml
```

### 3. Build & Load Application
Build the unified Docker image (current version: **v1**) and load it:
```bash
# Build image
docker build -t digital-post-office:v1 .

# If using Minikube
minikube image load digital-post-office:v1
```

### 4. Deploy Application
Apply all application manifests:
```bash
kubectl apply -f k8s/
```

### 5. Verify & Showcase

#### A. Open the Scaling Dashboard
Copy the URL and open it in your browser:
```bash
minikube service scaling-dashboard-service --url
```

#### B. Verify API Status
Ensure the API is reachable (returns all campaigns):
```bash
# Get API URL
minikube service email-api --url

# Verify (replace <URL>)
curl <URL>/api/campaigns/
```

#### C. Trigger Scaling Demo
Send a campaign with 10 recipients. **Choose the command for your terminal type:**

**For BASH (Git Bash, WSL, Linux, macOS):**
```bash
curl -X POST <API_URL>/api/campaigns/ \
-H 'Content-Type: application/json' \
-d '{"subject": "Scaling Test", "body": "Hello!", "recipients": ["u1@ex.com", "u2@ex.com", "u3@ex.com", "u4@ex.com", "u5@ex.com", "u6@ex.com", "u7@ex.com", "u8@ex.com", "u9@ex.com", "u10@ex.com"]}'
```

**For POWERSHELL:**
```powershell
curl.exe -X POST <API_URL>/api/campaigns/ `
-H "Content-Type: application/json" `
-d "{\`"subject\`": \`"Scaling Test\`", \`"body\`": \`"Hello!\`", \`"recipients\`": [\`"u1@ex.com\`", \`"u2@ex.com\`", \`"u3@ex.com\`", \`"u4@ex.com\`", \`"u5@ex.com\`", \`"u6@ex.com\`", \`"u7@ex.com\`", \`"u8@ex.com\`", \`"u9@ex.com\`", \`"u10@ex.com\`"]}"
```

---

## 📁 Project Structure

- `cmd/server/`: API entrypoint and **Campaign Trigger UI** (`web/`).
- `cmd/dashboard/`: **Scaling Monitor** entrypoint and its UI (`web/`).
- `cmd/publisher/`: Transactional Outbox worker.
- `cmd/worker/`: Email delivery consumer.
- `k8s/`: Kubernetes manifests (Deployments, HPA, KEDA ScaledObjects).
- `internal/`: Core logic, models, and PostgreSQL repositories.
- `migrations/`: SQL migration files.
- `.agents/skills/`: Custom AI Agent Skills (One-click setup, deploys, etc.).

---

## 🛠️ AI Agent Tools (Skills)
This repository is optimized for AI assistance. You can ask an agent to:
- `/local-setup`: Run the entire deployment automatically.
- `/db-migrate`: Automatically apply new database changes.
- `/wire-gen`: Regenerate Go dependencies.
- `/go-test`: Run all unit and integration tests.
- `/email-resilience`: Understand and manage email delivery reliability.

---

## 👁️ Database Visibility
To see exactly what's happening inside the resilient database:
```bash
# List all campaigns
kubectl exec deployment/postgres -- psql -U postgres -d goDb -c "SELECT * FROM campaigns;"

# List the outbox (queuing status)
kubectl exec deployment/postgres -- psql -U postgres -d goDb -c "SELECT * FROM outbox;"

# List actual email logs (success/fail)
kubectl exec deployment/postgres -- psql -U postgres -d goDb -c "SELECT * FROM email_logs;"
```

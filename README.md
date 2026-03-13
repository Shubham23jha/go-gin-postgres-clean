# Distributed Email Delivery System (Mini-Mailchimp) 📧🚀

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
Build the unified Docker image (current version: **v2**) and load it:
```bash
# Build image
docker build -t email-system:v2 .

# If using Minikube
minikube image load email-system:v2
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

- `cmd/`: Entrypoints for API, Publisher, Worker, and Dashboard.
- `k8s/`: Kubernetes Manifests (Deployments, Config, KEDA).
- `internal/`: Core business logic, Handlers, Models, and Repositories.
- `migrations/`: SQL database schema migrations.
- `.agents/`: AI Agent Skills (Automated DB migrations, K8s deploys, etc.).

---

## 🛠️ Maintenance (AI Agent Skills)
This repository is optimized for AI assistance. Skills in `.agents/skills/` include:
- `db-migrate`: Manage SQL schema.
- `k8s-deploy`: Orchestrate cluster resources.
- `wire-gen`: Regenerate dependency injection code.

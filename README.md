# Distributed Email Delivery System (Mini-Mailchimp) 📧🚀

A robust, high-throughput email delivery system built with **Go**, **Postgres**, **RabbitMQ**, and orchestrated with **Kubernetes** + **KEDA**. Featuring a live **Scaling Visualization Dashboard**.

---

## 🏗️ Architecture
The system uses the **Transactional Outbox Pattern** to ensure "at-least-once" delivery and decoupled horizontal scaling.

1.  **Campaign API**: Atomic transaction that saves the campaign and its individual email tasks to an `outbox` table.
2.  **Outbox Publisher**: Polling service that reads `PENDING` outbox items and publishes them to RabbitMQ.
3.  **RabbitMQ**: Reliable message broker with built-in **Dead Letter Queues (DLQ)**.
4.  **Email Worker Pool**: A cluster of consumers that deliver emails via SMTP with exponential backoff and idempotency.
5.  **KEDA**: Automatically scales Worker pods from **0 to 10** based on queue length.
6.  **📊 Scaling Dashboard**: A real-time monitoring UI that visualizes queue growth and worker pod scaling.

---

## 🚀 Step-by-Step Local Setup (Kubernetes + KEDA)

To showcase the full system with autoscaling, follow these steps:

### 1. Pre-requisites
- **Docker Desktop** (with Kubernetes enabled) or **Minikube**.
- **Go 1.24**.
- **KEDA** installed in your cluster.

### 2. Install KEDA (If not present)
```powershell
kubectl apply --server-side -f https://github.com/kedacore/keda/releases/download/v2.16.1/keda-2.16.1.yaml
```

### 3. Start Infrastructure
We use Kubernetes for everything. Start by deploying the database and RabbitMQ:
```powershell
kubectl apply -f k8s/infrastructure.yaml
```

### 4. Build & Load Application
Build the unified Docker image and load it into your local cluster:
```powershell
# Build image
docker build -t email-system:v1 .

# If using Minikube
minikube image load email-system:v1
```

### 5. Deploy Application
Apply all application manifests:
```powershell
kubectl apply -f k8s/
```

### 6. Verify Scaling Activity
1.  **Open Dashboard**: Use Minikube to get the URL:
    ```powershell
    minikube service scaling-dashboard-service --url
    ```
2.  **Watch Pods**:
    ```powershell
    kubectl get pods -w
    ```
3.  **Trigger Campaign**:
    Send a campaign with multiple recipients (e.g., 10+) to see KEDA scale up workers:
    ```powershell
    # Get API URL
    minikube service email-api --url
    
    # Send Request (replace URL)
    curl.exe -X POST http://<URL>/api/campaigns/ `
    -H "Content-Type: application/json" `
    -d "{\`"subject\`": \`"Scaling Test\`", \`"body\`": \`"Hello!\`", \`"recipients\`": [\`"u1@ex.com\`", \`"u2@ex.com\`", ...]}"
    ```

---

## 📁 Project Structure

```bash
├── cmd/
│   ├── server/          # API Entrypoint
│   ├── publisher/       # Standalone Outbox Publisher
│   ├── worker/          # Standalone Email Worker
│   └── dashboard/       # Scaling Visualization Dashboard
├── k8s/                 # Kubernetes Manifests (Deployments, Config, KEDA)
├── internal/
│   ├── handlers/        # API Request Handlers
│   ├── models/          # DB Schema & GORM Models
│   ├── repository/      # Data Access Layer
│   ├── service/         # Business Logic & Worker Pool
│   └── bootstrap/       # Dependency Injection (Wire)
├── pkg/
│   ├── database/        # Postgres connection & Migrations
│   └── utils/           # Shared Utilities (JWT, etc.)
├── migrations/          # SQL Migration files
└── .agents/             # AI Agent Skills & Workflows
```

---

## 🛠️ Maintenance & AI Agent Skills
This repository is "Agent-Ready." It includes specialized **Skills** in the `.agents/skills/` directory:

- **`db-migrate`**: Manage SQL schema migrations.
- **`docker-build`**: Standardized multi-stage image building.
- **`go-lint`**: Enforce Go coding standards.
- **`go-test`**: Run unit and integration test suites.
- **`k8s-deploy`**: Orchestrate Kubernetes resource deployment.
- **`wire-gen`**: Regenerate dependency injection code.

---

## 🛡️ Key Features
- **Transactional Outbox**: Ensures "at-least-once" delivery.
- **Auto-Scaling**: KEDA handles worker resource allocation based on queue length.
- **Idempotency**: Message IDs prevent duplicate emails.
- **Real-time Monitoring**: Visual feedback on cluster state.

# Distributed Email Delivery System (Mini-Mailchimp) 📧🚀

A robust, high-throughput email delivery system built with **Go**, **Postgres**, **RabbitMQ**, and orchestrated with **Kubernetes** + **KEDA**. Featuring a live **Scaling Visualization Dashboard**.

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/Shubham23jha/go-gin-postgres-clean)

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
└── Dockerfile           # Multi-stage Build for all services
```

---

## 🚀 Quick Start (Local Docker Compose)

The easiest way to see the system in action:

1. **Pre-requisites**: Docker Desktop & Go 1.24.
2. **Setup Env**: Copy `.env.example` to `.env` and add your [Mailtrap](https://mailtrap.io/) credentials.
3. **Run Infrastructure**:
   ```bash
   docker-compose up -d
   ```
4. **Run Monolith Server**:
   ```bash
   SET RUN_BACKGROUND_SERVICES=true
   go run cmd/server/main.go
   ```

---

## 📊 Live Scaling Visualization (Kubernetes)

To witness KEDA dynamically scaling workers based on load:

1. **Enable KEDA**:
   ```bash
   kubectl apply --server-side -f https://github.com/kedacore/keda/releases/download/v2.16.1/keda-2.16.1.yaml
   ```
2. **Build & Deploy**:
   ```bash
   docker build -t email-system:latest .
   kubectl apply -f k8s/
   ```
3. **Open Dashboard**: Visit `http://localhost:30080` to see the live scaling UI.
4. **Trigger Scale Up**: Send ~20 emails via the API (see Testing section).

---

## 🧪 Testing the Flow

Create a campaign and watch it process automatically:

```powershell
curl -X POST http://localhost:8080/api/campaigns/ `
-H "Content-Type: application/json" `
-d '{
  "subject": "Distributed Scale Test",
  "body": "Hello from your autoscaling cluster!",
  "recipients": ["user1@example.com", "user2@example.com", "user3@example.com"]
}'
```

---

## 🛠️ Maintenance & Agent Skills
This repository is "Agentic Ready." It includes specialized **Skills** in the `.agent/skills/` directory that AI agents (like Antigravity) can use to manage the project:

- **`db-migrate`**: Automatic SQL schema migrations.
- **`docker-build`**: Standardized multi-stage image building.
- **`go-lint`**: Enforces Go coding standards.
- **`go-test`**: Runs unit and integration test suites.
- **`k8s-deploy`**: Orchestrates Kubernetes resource deployment.
- **`wire-gen`**: Regenerates dependency injection code.

---

## 🛡️ Key Features
- **Transactional Consistency**: Database-backed outbox ensures no lost emails.
- **Auto-Scaling**: KEDA handles worker resource allocation efficiently.
- **Idempotency**: Message IDs prevent duplicate emails during retries.
- **Real-time Monitoring**: Visual feedback on cluster state and queue depth.

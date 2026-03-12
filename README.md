# Distributed Email Delivery System (Mini-Mailchimp) 📧🚀

A robust, high-throughput email delivery system built with **Go**, **PostgreSQL**, **RabbitMQ**, and orchestrated with **Kubernetes** + **KEDA**.

## 🏗️ Architecture
The system uses the **Transactional Outbox Pattern** to ensure "at-least-once" delivery and decoupled horizontal scaling.

1.  **Campaign API**: Atomic transaction that saves the campaign and its individual email tasks to an `outbox` table.
2.  **Outbox Publisher**: Polling service that reads `PENDING` outbox items and publishes them to RabbitMQ.
3.  **RabbitMQ**: Serves as the reliable message broker with built-in **Dead Letter Queues (DLQ)**.
4.  **Email Worker Pool**: A cluster of consumers that deliver emails via SMTP with:
    - **Exponential Backoff**: Respects SMTP provider rate limits.
    - **Idempotency**: Prevents duplicate sends using a `message_id` registry in the database.
5.  **KEDA**: Automatically scales the number of Worker pods from **0 to 10** based on the real-time queue length.

---

## 📁 Project Structure

```bash
├── cmd/
│   ├── server/          # API Entrypoint (Monolith or Distributed)
│   ├── publisher/       # Standalone Outbox Publisher
│   └── worker/          # Standalone Email Worker
├── k8s/                 # Kubernetes Manifests (Deployments, Config, KEDA)
├── internal/
│   ├── models/          # DB Schema & GORM Models
│   ├── repository/      # Data Access Layer (Campaigns, Logs, Outbox)
│   ├── service/         # Business Logic & Worker Pool implementation
│   └── bootstrap/       # Dependency Injection (Wire)
├── pkg/
│   ├── database/        # Postgres connection & Migration logic
│   └── smtp/            # (Future) SMTP provider wrappers
├── migrations/          # SQL Migration files
└── Dockerfile           # High-performance Multi-stage Build
```

---

## 🚀 Getting Started

### 1. Pre-requisites
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Go 1.23+](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/) & [RabbitMQ](https://www.rabbitmq.com/) (or use Docker Compose)
- [Mailtrap Account](https://mailtrap.io/) (for testing)

### 2. Environment Setup
Copy `.env.example` to `.env` and fill in your Mailtrap credentials:
```bash
DB_PASSWORD=your_pass
RABBITMQ_PASS=your_pass
SMTP_USER=0cfc585...
SMTP_PASS=your_mailtrap_pass
```

### 3. Local Development (Monolith Mode)
To run everything in one process:
1. Set `RUN_BACKGROUND_SERVICES=true` in `.env`.
2. Run: `go run cmd/server/main.go`

### 4. Distributed Mode (Kubernetes)
1. **Enable KEDA**:
   ```powershell
   kubectl apply --server-side -f https://github.com/kedacore/keda/releases/download/v2.16.1/keda-2.16.1.yaml
   ```
2. **Build Image**:
   ```powershell
   docker build -t email-system:latest .
   ```
3. **Deploy Everything**:
   ```powershell
   kubectl apply -f k8s/
   ```

---

## 🧪 Testing the Flow
Create a campaign using cURL:
```powershell
curl -X POST http://localhost:8080/api/campaigns/ `
-H "Content-Type: application/json" `
-d '{
  "subject": "Distributed Scale Test",
  "body": "Hello from your autoscaling cluster!",
  "recipients": ["user1@example.com", "user2@example.com"]
}'
```

---

## 🛠️ Maintenance & Skills
The project includes AI-assisted **Skills** in `.agent/skills/` for:
- `db-migrate`: Schema updates.
- `docker-build`: Containerization.
- `k8s-deploy`: Kubernetes management.
- `wire-gen`: Dependency Injection.
- `go-lint/test`: Quality assurance.

---

## 🛡️ Reliability Features
- **Transactional Consistency**: If DB fails, no email is ever lost or double-accounted.
- **Circuit Breaking/Backoff**: Respects `550 Too many emails` errors from SMTP.
- **Visibility**: Every single attempt is logged in the `email_logs` table.
- **Recoverability**: Failed messages after 3 retries move to `email_dlq` for manual retry.
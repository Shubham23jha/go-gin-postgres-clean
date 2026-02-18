## 🏗 Architecture

This project follows **Clean Architecture** principles:

- Handler → Service → Repository → Database
- Dependency inversion via interfaces
- Scalable & testable structure



## 📁 Project Structure

```
go-gin-postgres-clean/
│
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
│
├── config/
│   ├── config.go               # Environment configuration loader
│   └── database.go             # Database config struct
│
├── internal/
│   ├── models/                 # GORM models (DB schema)
│   │   └── user.go
│   │
│   ├── repository/             # Data access layer
│   │   ├── interfaces.go
│   │   └── user_repository.go
│   │
│   ├── service/                # Business logic layer
│   │   └── user_service.go
│   │
│   ├── handler/                # HTTP handlers / controllers
│   │   └── user_handler.go
│   │
│   ├── routes/
│   │   └── routes.go           # Route registration
│   │
│   └── middleware/
│       └── logger.go           # Request logging middleware
│
├── pkg/
│   ├── database/
│   │   ├── postgres.go         # GORM DB connection setup
│   │   └── migrate.go          # Migration runner
│   │
│   └── utils/
│       └── response.go         # Standard API response helpers
│
├── migrations/                 # SQL migration files
│   ├── 000001_create_users.up.sql
│   └── 000001_create_users.down.sql
│
├── docker/
│   └── Dockerfile              # App container image
│
├── scripts/
│   └── migrate.sh              # Migration helper script
│
├── .env                        # Environment variables
├── .env.example                # Sample env file
├── docker-compose.yml          # Multi-container setup
├── Makefile                    # Dev commands automation
├── go.mod                      # Go module definition
└── go.sum                      # Dependency checksums
```


## create Migration
```
migrate create -ext sql -dir migrations -seq add_orders_table

```

## check Version

```
migrate -path migrations -database "<DB_URL>" version

```
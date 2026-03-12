---
name: db-migrate
description: Manage database migrations using golang-migrate.
---

# db-migrate Skill

This skill allows the agent to create and run database migrations.

## Commands

### Create a new migration
Replace `<name>` with a descriptive name for the migration.
```bash
migrate create -ext sql -dir migrations -seq <name>
```

### Run migrations (up)
Requires `DB_URL` environment variable.
```bash
migrate -path migrations -database "${DB_URL}" up
```

### Rollback migrations (down)
Requires `DB_URL` environment variable.
```bash
migrate -path migrations -database "${DB_URL}" down 1
```

### Check migration version
Requires `DB_URL` environment variable.
```bash
migrate -path migrations -database "${DB_URL}" version
```

### Use helper script (if applicable)
```bash
./scripts/migrate.sh
```

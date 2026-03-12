---
name: go-lint
description: Format and lint Go code.
---

# go-lint Skill

This skill allows the agent to format and lint the project's Go code.

## Commands

### Format all code
```bash
go fmt ./...
```

### Run static analysis (go vet)
```bash
go vet ./...
```

### Run golangci-lint (if installed)
```bash
golangci-lint run
```

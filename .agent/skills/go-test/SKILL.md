---
name: go-test
description: Run Go tests for the project.
---

# go-test Skill

This skill allows the agent to run tests for the Go project.

## Commands

### Run all tests
```bash
go test ./... -v
```

### Run tests for a specific package
Replace `<package_path>` with the path to the package.
```bash
go test ./<package_path> -v
```

### Run a specific test
Replace `<test_name>` with the name of the test and `<package_path>` with the path to the package.
```bash
go test -v -run <test_name> ./<package_path>
```

## Integration Testing

### Run all integration tests
Ensures the test environment (Docker) is up before running.
```bash
# 1. Start the test database
docker-compose -f docker-compose.test.yml up -d

# 2. Run the tests
go test ./tests/integration/... -v

# 3. (Optional) Stop the test database
docker-compose -f docker-compose.test.yml down
```

### Setup for Integration Tests
1. Ensure Docker is running.
2. The `TestMain` function in `tests/integration/test_helper_test.go` handles migrations and global DB initialization.
3. Test database runs on port `5433` by default to avoid conflicts.

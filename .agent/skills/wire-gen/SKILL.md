---
name: wire-gen
description: Generate dependency injection code using Wire.
---

# Wire Generation Skill

Use this skill to update the `wire_gen.go` files when dependencies change.

## Usage

### 1. Generate for initial bootstrap
// turbo
```powershell
cd internal/bootstrap; wire
```

### 2. Generate for services
// turbo
```powershell
cd internal/service; wire
```

### 3. Generate for repositories and handlers
// turbo
```powershell
cd internal/repository; wire
cd internal/handlers; wire
```

## When to use
- After adding a new dependency to a `New...` constructor.
- After adding a new provider to a `wire.Build` call.
- If you see compilation errors related to missing arguments in `wire_gen.go`.

# Quick Start Guide - Go Convert Service

## TL;DR - Launch in 5 seconds

```bash
# Option 1: Using scripts (simplest)
./scripts/start-service.sh

# Option 2: Using Make
make run

# Option 3: Using Docker
./scripts/start-docker.sh
```

## Choose Your Launch Method

### 1️⃣ Scripts (Recommended for Development)

**Start service:**
```bash
./scripts/start-service.sh
```

**With custom configuration:**
```bash
PORT=9000 LOG_LEVEL=info ./scripts/start-service.sh
```

**Start in Docker:**
```bash
./scripts/start-docker.sh
```

**Stop:**
```bash
./scripts/stop-service.sh
```

---

### 2️⃣ Makefile (Recommended for CI/CD)

```bash
# Build and run
make run

# Run with debug logging
make run-debug

# Run in Docker
make docker-run-foreground

# See all commands
make help
```

---

### 3️⃣ VS Code (Recommended for Debugging)

1. Press `F5`
2. Select "Launch go-convert Service"
3. Service starts with debugger attached ✨

**Keyboard shortcuts:**
- `F5` - Start debugging
- `⇧F5` - Stop
- `⌘⇧F5` - Restart

---

### 4️⃣ IntelliJ/GoLand

1. Open run configurations dropdown (top right)
2. Select "Go Convert Service"
3. Click ▶️ Run or 🐛 Debug

**Available configurations:**
- Go Convert Service (port 8090)
- Go Convert Service (Port 9000)
- Docker: Go Convert Service

---

## Test the Service

### Health Check
```bash
curl http://localhost:8090/healthz
```

**Expected:** `{"status":"ok"}`

### Example Request
```bash
curl -X POST http://localhost:8090/api/v1/convert/batch \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "id": "test-1",
        "entity_type": "pipeline",
        "yaml": "pipeline:\n  identifier: test\n  name: Test\n  stages: []"
      }
    ]
  }'
```

### With Template Reference Mapping
```bash
curl -X POST http://localhost:8090/api/v1/convert/batch \
  -H "Content-Type: application/json" \
  -d @test_batch_with_mapping.json
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8090 | HTTP port |
| `LOG_LEVEL` | debug | Log level (debug/info/warn/error) |
| `MAX_BATCH_SIZE` | 100 | Max items per batch |
| `MAX_YAML_BYTES` | 1048576 | Max request size (1MB) |

**Example:**
```bash
PORT=9000 LOG_LEVEL=info make run
```

---

## Common Issues

### Port Already in Use
```bash
# Change port
PORT=9000 ./scripts/start-service.sh

# Or kill existing process
./scripts/stop-service.sh
```

### Build Fails
```bash
# Download dependencies
make mod-download

# Rebuild
make clean build
```

### Docker Issues
```bash
# Clean and rebuild
make docker-stop
docker system prune -a
make docker-build
```

---

## Next Steps

- 📖 Read `SERVICE.md` for detailed documentation
- 🏗️ Read `TECH_SPEC.md` for architecture details
- 🧪 Run tests: `make test`
- 🐳 Deploy to Kubernetes (see SERVICE.md)

---

## Quick Reference

| Task | Command |
|------|---------|
| Start | `./scripts/start-service.sh` |
| Stop | `./scripts/stop-service.sh` |
| Build | `make build` |
| Test | `make test` |
| Docker | `./scripts/start-docker.sh` |
| Health | `curl localhost:8090/healthz` |
| Logs | `docker logs -f go-convert-service` |
| Help | `make help` |

# Go Convert

> 🔄 Convert third-party CI/CD pipelines and Harness v0 YAML to Harness v1 format

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

## Overview

Go Convert provides comprehensive tooling to convert CI/CD pipeline configurations to the Harness pipeline format. It supports multiple platforms and offers flexible deployment options for different use cases.

### 🎯 Key Features

- ✅ **Multi-Platform Support**: Bitbucket, Drone, GitLab, Jenkins → Harness
- ✅ **v0 to v1 Migration**: Convert Harness legacy format to modern v1
- ✅ **Batch Processing**: Convert up to 100 items in a single API call
- ✅ **Template Reference Mapping**: Update nested template references during conversion
- ✅ **Three Deployment Modes**: Library, CLI, or HTTP Microservice
- ✅ **Production Ready**: Docker support, graceful shutdown, health checks

## 🚀 Three Ways to Use

| Mode | Best For | Quick Start |
|------|----------|-------------|
| **📚 Go Library** | Embedded in Go applications | `import "github.com/drone/go-convert"` |
| **🖥️ Command Line** | Local development & debugging | `./go-convert bitbucket pipeline.yml` |
| **🌐 HTTP Microservice** ⭐ | Language-agnostic API access | `./scripts/start-service.sh` |

## 📑 Table of Contents

- [HTTP Microservice](#http-microservice) ⭐ NEW
- [Go Library Usage](#go-library-usage)
- [Command Line Tools](#command-line-tools)
- [Microservice Deployment](#microservice-deployment)
- [API Examples](#api-examples)
- [Development](#development)
- [IDE Integration](#ide-integration)
- [Contributing](#contributing)

---

## HTTP Microservice

> **New in this release!** Convert Harness v0 YAML (pipelines, templates, input sets) to v1 format via REST API with batch processing and template reference mapping support.

### Architecture

```
┌─────────────────┐
│   Client App    │
│  (Any Language) │
└────────┬────────┘
         │ HTTP/JSON
         ▼
┌─────────────────────────────────┐
│   Go Convert Microservice       │
│  ┌────────────────────────────┐ │
│  │  POST /api/v1/convert/batch│ │
│  │  - Pipelines               │ │
│  │  - Templates               │ │
│  │  - Input Sets              │ │
│  │  - Template Ref Mapping    │ │
│  └────────────────────────────┘ │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│   Conversion Engine             │
│  - v0 → v1 Pipeline Converter   │
│  - Template Converter           │
│  - Input Set Converter          │
│  - Reference Replacer           │
└─────────────────────────────────┘
```

### Quick Start

```bash
# Start the service
./scripts/start-service.sh

# Or using Make
make run

# Or using Docker
./scripts/start-docker.sh
```

Service starts on `http://localhost:8090`

### API Endpoints

- `GET /healthz` - Health check
- `POST /api/v1/convert/batch` - Batch convert pipelines, templates, and input sets

### Example Request

```bash
curl -X POST http://localhost:8090/api/v1/convert/batch \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "id": "pipeline-1",
        "entity_type": "pipeline",
        "yaml": "<v0 pipeline YAML>",
        "template_ref_mapping": {
          "oldTemplateRef": "newTemplateRef_v1"
        },
        "pipeline_ref_mapping": {
          "oldPipelineId": "newPipelineId_v1"
        }
      },
      {
        "id": "template-1",
        "entity_type": "template",
        "yaml": "<v0 template YAML>"
      }
    ]
  }'
```

### Features

✅ **Batch Processing**: Convert up to 100 items per request  
✅ **Multi-Entity Support**: Pipelines, templates (Pipeline/Stage/Step), and input sets  
✅ **Template & Pipeline Reference Mapping**: Rewrite template refs and pipeline identifiers independently during conversion  
✅ **Integrity Verification**: SHA-256 checksums for all conversions  
✅ **Production Ready**: Structured logging, graceful shutdown, health checks  
✅ **Container Ready**: Docker and Kubernetes support  
✅ **Language Agnostic**: REST API accessible from any language

### When to Use Which Mode?

| Use Case | Recommended Mode | Why? |
|----------|------------------|------|
| Integrate into Go application | 📚 Library | Direct function calls, type safety |
| Quick local conversion | 🖥️ CLI | Simple, no setup required |
| Production API service | 🌐 Microservice | Scalable, language-agnostic |
| Batch migrations | 🌐 Microservice | Process 100s of files efficiently |
| CI/CD pipeline integration | 🌐 Microservice | REST API, easy integration |
| Development & debugging | 🖥️ CLI | Fast iteration, syntax highlighting |

### Documentation

- **[QUICKSTART.md](QUICKSTART.md)** - Get the service running in 30 seconds
- **[SERVICE.md](SERVICE.md)** - Complete API documentation and deployment guide
- **[TECH_SPEC.md](TECH_SPEC.md)** - Architecture and design details

### Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `PORT` | 8090 | HTTP listen port |
| `LOG_LEVEL` | debug | Logging level (debug/info/warn/error) |
| `MAX_BATCH_SIZE` | 100 | Maximum items per batch request |
| `MAX_YAML_BYTES` | 1048576 | Maximum request body size (1MB) |

### Performance & Scalability

- 🚀 **Lightweight**: ~10MB binary, minimal memory footprint
- ⚡ **Fast**: Converts typical pipeline in <100ms
- 📦 **Batch Optimized**: Process up to 100 items per request
- 🔄 **Stateless**: Easily scale horizontally
- 🐳 **Container Ready**: Small distroless Docker image (~20MB)

---

## Supported Conversions

### Platform Migrations

```
Bitbucket Pipelines  ──┐
Drone CI             ──┤
GitLab CI            ──┼──➤  Harness v1 Pipeline
Jenkins              ──┤
GitHub Actions       ──┘
```

### Harness v0 → v1 Migration

```
Harness v0 Pipeline  ──➤  Harness v1 Pipeline
Harness v0 Template  ──➤  Harness v1 Template (Pipeline/Stage/Step)
Harness v0 InputSet  ──➤  Harness v1 InputSet
```

**Special Features:**
- 🔗 Template reference mapping during conversion
- ✅ Validates converted YAML structure
- 📝 Generates SHA-256 checksums

---

## Go Library Usage

__Sample Usage__

Sample code to convert a Bitbucket pipeline to a Harness pipeline:

```Go
import "github.com/drone/go-convert/convert/bitbucket"
```

```Go
converter := bitbucket.New(
	bitbucket.WithDockerhub(c.dockerConn),
	bitbucket.WithKubernetes(c.kubeConn, c.kubeName),
)
converted, err := converter.ConvertFile("bitbucket-pipelines.yml")
if err != nil {
	log.Fatalln(err)
}
```

---

## Command Line Tools

This package provides command line tools for local development and debugging purposes. These command line tools are intentionally simple. For more robust command line tooling please use the [harness-convert](https://github.com/harness/harness-convert) project.

### Installation

```bash
git clone https://github.com/drone/go-convert.git
cd go-convert
go build
```

__Bitbucket__

Convert a Bitbucket pipeline:

```
./go-convert bitbucket samples/bitbucket.yaml
```

Convert a Gitlab pipeline and print the before after:

```
./go-convert bitbucket --before-after samples/bitbucket.yaml
```

Convert a Bitbucket pipeline and downgrade to the Harness v0 format:

```
./go-convert bitbucket --downgrade samples/bitbucket.yaml
```

__Drone__

Convert a Drone pipeline:

```
./go-convert drone samples/drone.yaml
```

Convert a Drone pipeline and print the before after:

```
./go-convert drone --before-after samples/drone.yaml
```

Convert a Drone pipeline and downgrade to the Harness v0 format:

```
./go-convert drone --downgrade samples/drone.yaml
```

__Gitlab__

Convert a Gitlab pipeline:

```
./go-convert gitlab samples/gitlab.yaml
```

Convert a Gitlab pipeline and print the before after:

```
./go-convert gitlab --before-after samples/gitlab.yaml
```

Convert a Gitlab pipeline and downgrade to the Harness v0 format:

```
./go-convert gitlab --downgrade samples/gitlab.yaml
```

__Jenkins__

Convert a Jenkinsfile:

```
./go-convert jenkins --token=<chat-gpt-token> samples/Jenkinsfile
```

Convert a Jenkinsfile and downgrade to the Harness v0 format:

```
./go-convert jenkins --token=<chat-gpt-token> --downgrade samples/Jenkinsfile
```

__Syntax Highlighting__

The command line tools are compatible with [bat](https://github.com/sharkdp/bat) for syntax highlight.

```
./go-convert bitbucket --before-after samples/bitbucket.yaml | bat -l yaml
```

---

## Microservice Deployment

### Local Development

```bash
# Using scripts
./scripts/start-service.sh

# Using Makefile
make run-debug

# With custom configuration
PORT=9000 LOG_LEVEL=info make run
```

### Docker

```bash
# Build image
make docker-build

# Run container
docker run -p 8090:8090 go-convert-service:latest

# Or use the script
./scripts/start-docker.sh
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-convert-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-convert-service
  template:
    metadata:
      labels:
        app: go-convert-service
    spec:
      containers:
      - name: go-convert-service
        image: go-convert-service:latest
        ports:
        - containerPort: 8090
        env:
        - name: LOG_LEVEL
          value: "info"
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8090
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /healthz
            port: 8090
          initialDelaySeconds: 3
          periodSeconds: 5
```

See [SERVICE.md](SERVICE.md) for complete deployment documentation.

---

## Project Structure

```
go-convert/
├── cmd/
│   └── server/          # HTTP microservice entrypoint
├── convert/             # Conversion logic for various CI/CD platforms
│   ├── bitbucket/       # Bitbucket pipelines converter
│   ├── drone/           # Drone pipelines converter
│   ├── gitlab/          # Gitlab CI converter
│   ├── jenkins/         # Jenkins converter
│   └── v0tov1/          # Harness v0 to v1 converter
├── service/             # HTTP service implementation
│   ├── converter/       # v0→v1 conversion handlers
│   │   ├── pipeline.go  # Pipeline conversion
│   │   ├── template.go  # Template conversion (Pipeline/Stage/Step)
│   │   ├── inputset.go  # Input set conversion
│   │   └── template_refs.go  # Template reference replacement
│   ├── handler.go       # HTTP request handlers
│   ├── server.go        # HTTP server and middleware
│   └── request.go       # Request/response types
├── scripts/             # Launch scripts
│   ├── start-service.sh
│   ├── start-docker.sh
│   └── stop-service.sh
├── .vscode/             # VS Code launch configurations
├── .idea/               # IntelliJ run configurations
├── Dockerfile.service   # Docker build configuration
├── Makefile            # Build and run automation
├── QUICKSTART.md       # Quick start guide
├── SERVICE.md          # Service documentation
└── TECH_SPEC.md        # Technical specification
```

---

## Development

### Build

```bash
# Build CLI tool
go build -o go-convert

# Build microservice
go build -o go-convert-service ./cmd/server

# Or use Makefile
make build
```

### Test

```bash
# Run tests
go test ./...

# With coverage
make test-coverage
```

### Format

```bash
make fmt
```

### Run with Live Reload

```bash
# Install air: make install-tools
make dev
```

---

## IDE Integration

### VS Code

1. Open project in VS Code
2. Press `F5` to start debugging
3. Select "Launch go-convert Service"

Available configurations:
- Launch go-convert Service (port 8090)
- Launch go-convert Service (Custom Port 9000)
- Attach to running service

### IntelliJ IDEA / GoLand

1. Open run configurations dropdown
2. Select "Go Convert Service"
3. Click Run (▶️) or Debug (🐛)

Available configurations:
- Go Convert Service
- Go Convert Service (Port 9000)
- Docker: Go Convert Service

---

## API Examples

### Health Check

```bash
curl http://localhost:8090/healthz
# Response: {"status":"ok"}
```

### Convert Pipeline

```bash
curl -X POST http://localhost:8090/api/v1/convert/batch \
  -H "Content-Type: application/json" \
  -d '{
    "items": [{
      "id": "pipeline-1",
      "entity_type": "pipeline",
      "yaml": "pipeline:\n  identifier: test\n  name: Test Pipeline\n  stages: []"
    }]
  }'
```

### Convert Template with Reference Mapping

```bash
curl -X POST http://localhost:8090/api/v1/convert/batch \
  -H "Content-Type: application/json" \
  -d '{
    "items": [{
      "id": "template-1",
      "entity_type": "template",
      "yaml": "template:\n  type: Stage\n  spec: {...}",
      "template_ref_mapping": {
        "oldRef": "newRef_v1"
      }
    }]
  }'
```

### Batch Conversion

```bash
curl -X POST http://localhost:8090/api/v1/convert/batch \
  -H "Content-Type: application/json" \
  -d @test_batch_with_mapping.json
```

---

## Makefile Commands

Run `make help` to see all available commands:

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make build` | Build the service binary |
| `make run` | Build and run locally |
| `make run-debug` | Run with debug logging |
| `make test` | Run tests |
| `make test-coverage` | Run tests with coverage |
| `make docker-build` | Build Docker image |
| `make docker-run` | Run in Docker (background) |
| `make docker-stop` | Stop Docker container |
| `make health-check` | Check service health |
| `make example-request` | Send example request |
| `make clean` | Remove build artifacts |

---

## Troubleshooting

### Common Issues

<details>
<summary><b>Port Already in Use</b></summary>

```bash
# Find what's using the port
lsof -i :8090

# Kill the process
kill -9 <PID>

# Or use a different port
PORT=9000 ./scripts/start-service.sh
```
</details>

<details>
<summary><b>Service Not Responding</b></summary>

```bash
# Check if service is running
curl http://localhost:8090/healthz

# Check logs (if running in Docker)
docker logs go-convert-service

# Restart the service
./scripts/stop-service.sh
./scripts/start-service.sh
```
</details>

<details>
<summary><b>Docker Build Fails</b></summary>

```bash
# Clean Docker cache
docker system prune -a

# Rebuild
make docker-build
```
</details>

<details>
<summary><b>Module Dependencies Issues</b></summary>

```bash
# Download dependencies
make mod-download

# Tidy dependencies
make mod-tidy

# Rebuild
make clean build
```
</details>

---

## Quick Reference Card

### 🚀 Start Service
```bash
./scripts/start-service.sh          # Local
./scripts/start-docker.sh           # Docker
make run                            # Using Make
```

### 🧪 Test Service
```bash
curl localhost:8090/healthz         # Health check
make example-request                # Send test request
```

### 🛠️ Build & Test
```bash
make build                          # Build binary
make test                           # Run tests
make docker-build                   # Build Docker image
```

### 📚 Documentation
- [QUICKSTART.md](QUICKSTART.md) - Get started in 30 seconds
- [SERVICE.md](SERVICE.md) - Complete API documentation
- [TECH_SPEC.md](TECH_SPEC.md) - Architecture details

### 🔗 Useful Links
- API Endpoint: `http://localhost:8090/api/v1/convert/batch`
- Health Check: `http://localhost:8090/healthz`
- Test Files: `test_batch_request.json`, `test_batch_with_mapping.json`

---

## Contributing

We welcome contributions! Here's how:

1. **Fork** the repository
2. **Create** your feature branch
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Make** your changes
4. **Test** your changes
   ```bash
   make test
   make fmt
   ```
5. **Commit** your changes
   ```bash
   git commit -m 'Add amazing feature'
   ```
6. **Push** to the branch
   ```bash
   git push origin feature/amazing-feature
   ```
7. **Open** a Pull Request

### Development Guidelines

- Write tests for new features
- Follow existing code style (run `make fmt`)
- Update documentation as needed
- Ensure all tests pass (`make test`)

---

## Related Projects

- **[harness-convert](https://github.com/harness/harness-convert)** - Production-grade CLI tool
- **Harness Platform** - [https://harness.io](https://harness.io)

---

## License

Apache 2.0 - See [LICENSE](LICENSE) file for details.

---

## Support & Resources

### 📖 Documentation
- **Quick Start**: [QUICKSTART.md](QUICKSTART.md)
- **Service Guide**: [SERVICE.md](SERVICE.md)
- **Technical Spec**: [TECH_SPEC.md](TECH_SPEC.md)

### 🐛 Issues & Questions
- **Bug Reports**: [GitHub Issues](https://github.com/drone/go-convert/issues)
- **Feature Requests**: [GitHub Discussions](https://github.com/drone/go-convert/discussions)

### 🎯 Production Use
For production-grade command line tooling, use [harness-convert](https://github.com/harness/harness-convert).

---

<div align="center">

**Made with ❤️ by the Harness Team**

[Documentation](SERVICE.md) • [Quick Start](QUICKSTART.md) • [Technical Spec](TECH_SPEC.md)

</div>

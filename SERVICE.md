# Go Convert Service

HTTP microservice for converting Harness v0 YAML (pipelines, templates, input sets) to v1 format.

## Quick Start

### Using Scripts (Recommended)

```bash
# Start service locally (builds and runs)
./scripts/start-service.sh

# Start with custom configuration
PORT=9000 LOG_LEVEL=info ./scripts/start-service.sh

# Start in Docker
./scripts/start-docker.sh

# Stop service (local or Docker)
./scripts/stop-service.sh
```

### Using Makefile

```bash
# Show all available commands
make help

# Build and run locally
make run

# Run with debug logging
make run-debug

# Run in Docker (foreground)
make docker-run-foreground

# Run tests
make test

# Health check
make health-check

# Send example request
make example-request
```

### Using VS Code

1. Open project in VS Code
2. Press `F5` or go to Run and Debug
3. Select "Launch go-convert Service"
4. Service will build and start with debugger attached

**Available VS Code configurations:**
- `Launch go-convert Service` - Debug on port 8090
- `Launch go-convert Service (Custom Port)` - Debug on port 9000
- `Attach to running service` - Attach debugger to running process

### Using IntelliJ/GoLand

1. Open project in IntelliJ/GoLand
2. Select run configuration from dropdown (top right)
3. Choose:
   - `Go Convert Service` - Run on port 8090
   - `Go Convert Service (Port 9000)` - Run on port 9000
   - `Docker: Go Convert Service` - Run in Docker
4. Click Run (▶️) or Debug (🐛)

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8090 | HTTP listen port |
| `LOG_LEVEL` | info | Logging level (debug, info, warn, error) |
| `MAX_BATCH_SIZE` | 100 | Maximum items per batch request |
| `MAX_YAML_BYTES` | 1048576 | Maximum request body size (1MB) |

## API Documentation

### Endpoints

- `GET  /healthz` — Health check
- `POST /api/v1/convert/pipeline` — Convert v0 pipeline YAML → v1
- `POST /api/v1/convert/template` — Convert v0 template YAML → v1
- `POST /api/v1/convert/input-set` — Convert v0 input set YAML → v1
- `POST /api/v1/convert/trigger` — Convert v0 trigger YAML → v1 (inputYaml is converted in place)
- `POST /api/v1/convert/batch` — Convert multiple entities in one call
- `POST /api/v1/convert/expression` — Convert one or more Harness v0 expressions (see `EXPRESSION_CONVERSION_API.md`)
- `POST /api/v1/checksum` — Compute SHA-256 of a YAML payload

A matching gRPC surface (`pb.GoConvertService`) is exposed on a separate port (default `9090`).

Full endpoint, request, response, and error reference: see `TECH_SPEC.md` §3.

### Single-Entity Request Format

```json
{
  "yaml": "<v0 YAML string>",
  "entity_ref_mapping": { "oldTemplateRef": "newTemplateRef_v1" },
  "context_pipeline_yaml": "<optional v0 pipeline YAML for postprocess context>"
}
```

**Fields:**

- `yaml` (required) — v0 YAML payload to convert.
- `entity_ref_mapping` (optional) — string-level rewrite from old template/entity refs to v1 refs.
- `context_pipeline_yaml` (optional, **template / input-set / trigger only**) — raw v0 pipeline YAML used purely as expression-postprocess context. The server parses + structurally converts this pipeline (suppressing its diagnostic messages), harvests the resulting step-type map, and runs the entity's postprocess in **FQN mode** with that context. Empty/omitted → postprocess runs without FQN context. The pipeline endpoint ignores this field. See `TECH_SPEC.md` §3.10.

### Single-Entity Response Format

```json
{
  "yaml": "<v1 YAML string>",
  "checksum": "sha256:abc123...",
  "report": {
    "messages": [
      { "severity": "WARNING", "code": "UNKNOWN_STEP_TYPE", "message": "...", "context": { "type": "..." } }
    ],
    "unrecognized_fields": ["pipeline.foo"],
    "expressions": [
      { "original": "<+pipeline.variables.x>", "converted": "<+pipeline.inputs.x>", "status": "SUCCESS" }
    ]
  }
}
```

The `report` object (and its sub-arrays) are omitted when empty. See `TECH_SPEC.md` §3.6 for the full schema and the trigger-specific error codes.

### Batch Request Format

```json
{
  "items": [
    {
      "id": "unique-id-1",
      "entity_type": "pipeline",
      "yaml": "<v0 YAML string>",
      "entity_ref_mapping": {
        "oldTemplateRef": "newTemplateRef_v1"
      }
    }
  ]
}
```

**Fields:**
- `id` (required) — Your unique identifier, echoed in response
- `entity_type` (required) — `"pipeline"`, `"template"`, `"input-set"`, or `"trigger"`
- `yaml` (required) — Raw v0 YAML content
- `entity_ref_mapping` (optional) — Map old template refs to new v1 refs
- `context_pipeline_yaml` (optional) — v0 pipeline YAML used as postprocess context for non-pipeline entities (template / input-set / trigger). Same semantics as the single-entity endpoint.

### Batch Response Format

```json
{
  "results": [
    {
      "id": "unique-id-1",
      "entity_type": "pipeline",
      "yaml": "<v1 YAML string>",
      "checksum": "sha256:abc123...",
      "error": null,
      "report": { "messages": [], "unrecognized_fields": [], "expressions": [] }
    }
  ]
}
```

Per-item failures (parse/conversion errors) populate `error` and leave `yaml`/`checksum`/`report` as `null`. The outer call is always HTTP 200 unless the request itself is malformed.

## Testing

### Health Check

```bash
curl http://localhost:8090/healthz
```

Expected response:
```json
{"status":"ok"}
```

### Example Batch Request

```bash
curl -X POST http://localhost:8090/api/v1/convert/batch \
  -H "Content-Type: application/json" \
  -d @test_batch_request.json
```

### Example with Template Reference Mapping

```bash
curl -X POST http://localhost:8090/api/v1/convert/batch \
  -H "Content-Type: application/json" \
  -d @test_batch_with_mapping.json
```

## Development

### Build

```bash
# Using Make
make build

# Using Go directly
go build -o go-convert-service ./cmd/server
```

### Run Tests

```bash
# All tests
make test

# With coverage
make test-coverage
```

### Format Code

```bash
make fmt
```

### Run with Live Reload

```bash
# Install air first: make install-tools
make dev
```

## Docker

### Build Image

```bash
make docker-build
# or
docker build -f Dockerfile.service -t go-convert-service:latest .
```

### Run Container

```bash
# Background
make docker-run

# Foreground
make docker-run-foreground

# With custom config
docker run -p 9000:8090 -e LOG_LEVEL=info go-convert-service:latest
```

### View Logs

```bash
make docker-logs
# or
docker logs -f go-convert-service
```

### Stop Container

```bash
make docker-stop
# or
docker stop go-convert-service && docker rm go-convert-service
```

## Deployment

### Kubernetes

See `TECH_SPEC.md` for example Kubernetes deployment manifests.

### Production Settings

Recommended environment variables for production:

```bash
PORT=8090
LOG_LEVEL=info
MAX_BATCH_SIZE=100
MAX_YAML_BYTES=2097152  # 2MB
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using the port
lsof -i :8090

# Kill the process
kill -9 <PID>

# Or use a different port
PORT=9000 ./scripts/start-service.sh
```

### Docker Build Issues

```bash
# Clean Docker cache
docker system prune -a

# Rebuild
make docker-build
```

### Service Not Responding

```bash
# Check logs
docker logs go-convert-service

# Check health
curl http://localhost:8090/healthz

# Restart service
./scripts/stop-service.sh && ./scripts/start-service.sh
```

## Architecture

See `TECH_SPEC.md` for detailed architecture documentation.

**Key Components:**
- `cmd/server/main.go` — Service entrypoint
- `service/server.go` — HTTP server, gRPC server, and middleware
- `service/handler.go` — HTTP request handlers
- `service/grpc_handler.go` — gRPC handlers (mirrors HTTP)
- `service/report.go` — `ConversionReport` DTO and proto bridge
- `service/converter/` — Conversion logic
  - `pipeline.go` — Pipeline converter
  - `template.go` — Template converter (Pipeline/Stage/Step/StepGroup)
  - `inputset.go` — Input set converter
  - `trigger.go` — Trigger converter (inputYaml converted in place)
  - `expression.go` — Expression converter
  - `template_refs.go` — Template reference replacement

## Contributing

1. Make changes to code
2. Run tests: `make test`
3. Format code: `make fmt`
4. Build: `make build`
5. Test locally: `make run-debug`
6. Create PR

## Support

For issues or questions:
- Check `TECH_SPEC.md` for detailed documentation
- Review logs for error messages
- Test with example requests in `test_batch_request.json`

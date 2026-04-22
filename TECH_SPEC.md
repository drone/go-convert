# Tech Spec: go-convert Microservice

**Status:** Draft — for review before implementation  
**Author:** TBD  
**Date:** 2026-03-30

---

## 1. Overview

This document describes the design for wrapping the existing `go-convert` library into a standalone HTTP microservice. The service exposes REST APIs to convert Harness v0 YAML (legacy CI pipeline format) into Harness v1 YAML for three entity types: **pipeline**, **template**, and **input set**. Each conversion response includes a SHA-256 checksum of the resulting YAML.

### Goals

- Provide a language-agnostic HTTP interface to the existing conversion logic in `convert/v0tov1`
- Support pipeline, template, and input set YAML conversions
- Return a checksum alongside each converted artifact for integrity verification
- Keep the service thin — it delegates all conversion work to the existing library packages
- Be deployable as a standalone binary or container

### Non-Goals

- Authentication / authorization (deferred — a reverse proxy or API gateway handles that)
- Persistent storage of converted YAML
- Converting formats other than Harness v0 → v1

---

## 2. Architecture

### 2.1 High-level

```
  Client
    │
    ▼ HTTP/JSON (port 8090)
 ┌──────────────────────────┐
 │   service/server.go      │   HTTP router + middleware (logging, recovery)
 │   service/handler.go     │   HTTP handlers — one per entity type
 └────────┬─────────────────┘
          │  calls existing library packages
          ▼
 ┌──────────────────────────────────────────────────────┐
 │  convert/harness/yaml          (v0 parser)           │
 │  convert/v0tov1/pipeline_converter  (v0 → v1 logic)  │
 │  convert/v0tov1/yaml           (v1 marshal)          │
 │  service/converter/            (template + inputset  │
 │                                 converters — new)    │
 └──────────────────────────────────────────────────────┘
```

### 2.2 Directory layout (new files only)

```
go-convert/
├── cmd/
│   └── server/
│       └── main.go            # entrypoint — starts the HTTP server
├── service/
│   ├── server.go              # HTTP server setup, middleware chain, graceful shutdown
│   ├── handler.go             # HTTP handlers: ConvertPipeline, ConvertTemplate, ConvertInputSet
│   ├── request.go             # shared request / response structs
│   ├── checksum.go            # SHA-256 checksum helper
│   └── converter/
│       ├── pipeline.go        # thin wrapper around existing PipelineConverter
│       ├── template.go        # NEW — v0 template YAML → v1 template YAML
│       └── inputset.go        # NEW — v0 inputSet YAML → v1 inputSet YAML
└── TECH_SPEC.md               # this file
```

The rest of the repository (existing `convert/`, `command/`, etc.) is unchanged.

---

## 3. API Design

**Base URL:** `http://<host>:8090`  
**Content-Type:** all requests and responses use `application/json`  
**API version prefix:** `/api/v1`

### Endpoints summary

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/convert/pipeline` | Convert v0 pipeline YAML → v1 |
| `POST` | `/api/v1/convert/template` | Convert v0 template YAML → v1 |
| `POST` | `/api/v1/convert/input-set` | Convert v0 input set YAML → v1 |
| `POST` | `/api/v1/convert/batch` | Convert multiple entities in one call |
| `GET`  | `/healthz` | Health check |

---

### 3.1 Convert Pipeline

Converts a single Harness v0 pipeline YAML document into v1 format.

**Request**

```
POST /api/v1/convert/pipeline
Content-Type: application/json
```

```json
{
  "yaml": "<v0 pipeline YAML as a string>"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `yaml` | string | yes | Raw v0 pipeline YAML. Top-level key must be `pipeline:`. |

**Response — 200 OK**

```json
{
  "yaml": "<v1 pipeline YAML as a string>",
  "checksum": "sha256:<64-char hex digest of the v1 YAML bytes>"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `yaml` | string | Converted v1 pipeline YAML |
| `checksum` | string | `sha256:` prefix + lowercase hex SHA-256 of `yaml` as UTF-8 bytes |

**Example request body**

```json
{
  "yaml": "pipeline:\n  identifier: myPipeline\n  name: My Pipeline\n  orgIdentifier: default\n  projectIdentifier: MyProject\n  stages:\n    - stage:\n        identifier: build\n        name: Build\n        type: CI\n        spec:\n          execution:\n            steps:\n              - step:\n                  identifier: run\n                  name: Run Tests\n                  type: Run\n                  spec:\n                    command: go test ./...\n"
}
```

**Example response**

```json
{
  "yaml": "pipeline:\n  id: myPipeline\n  name: My Pipeline\n  stages:\n    - name: Build\n      type: ci\n      steps:\n        - name: Run Tests\n          type: run\n          spec:\n            run: go test ./...\n",
  "checksum": "sha256:a3f1c2d9e8b74056..."
}
```

---

### 3.2 Convert Template

Converts a Harness v0 template YAML document into v1 format.

> **Note:** The existing `go-convert` library has v1 `Template` types (`convert/v0tov1/yaml.Template`) and a v0 stage/step template reference model, but no end-to-end v0-template-entity → v1-template-entity converter. A new `service/converter/template.go` file will implement this conversion (see §6.2).

**Request**

```
POST /api/v1/convert/template
Content-Type: application/json
```

```json
{
  "yaml": "<v0 template YAML as a string>"
}
```

The v0 template YAML top-level key must be `template:` with the following shape:

```yaml
template:
  name: myTemplate
  identifier: myTemplate
  orgIdentifier: default
  projectIdentifier: MyProject
  type: Pipeline      # Pipeline | Stage | Step
  spec:
    # ... nested pipeline / stage / step spec matching the v0 schema
```

**Response — 200 OK**

```json
{
  "yaml": "<v1 template YAML as a string>",
  "checksum": "sha256:<hex digest>"
}
```

The v1 template YAML will have the shape:

```yaml
template:
  inputs:
    # declared runtime inputs (derived from <+input> expressions in the spec)
  pipeline: # or stage: / step: depending on type
    # ... converted v1 spec
```

---

### 3.3 Convert Input Set

Converts a Harness v0 input set YAML document into v1 format.

> **Note:** No end-to-end converter exists in the library today. A new `service/converter/inputset.go` will implement structural remapping (see §6.3).

**Request**

```
POST /api/v1/convert/input-set
Content-Type: application/json
```

```json
{
  "yaml": "<v0 input set YAML as a string>"
}
```

The v0 input set YAML top-level key must be `inputSet:`:

```yaml
inputSet:
  name: eventPR
  identifier: eventPR
  orgIdentifier: default
  projectIdentifier: MyProject
  pipeline:
    identifier: myPipeline
    properties:
      ci:
        codebase:
          build:
            type: PR
            spec:
              number: <+trigger.prNumber>
```

**Response — 200 OK**

```json
{
  "yaml": "<v1 input set YAML as a string>",
  "checksum": "sha256:<hex digest>"
}
```

The v1 input set YAML will have the shape:

```yaml
inputset:
  pipeline: myPipeline
  inputs:
    # key-value runtime inputs derived from the v0 pipeline fragment
```

---

### 3.4 Batch Convert

Converts multiple entities in a single request. Useful for migrating many files in one call without per-request overhead.

**Request**

```
POST /api/v1/convert/batch
Content-Type: application/json
```

```json
{
  "items": [
    {
      "id": "client-assigned-id-1",
      "entity_type": "pipeline",
      "yaml": "<v0 yaml>"
    },
    {
      "id": "client-assigned-id-2",
      "entity_type": "template",
      "yaml": "<v0 yaml>"
    },
    {
      "id": "client-assigned-id-3",
      "entity_type": "input-set",
      "yaml": "<v0 yaml>"
    }
  ]
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `items` | array | yes | List of items to convert. Max 100 per request. |
| `items[].id` | string | yes | Client-provided identifier echoed in the response to correlate results. |
| `items[].entity_type` | string | yes | `"pipeline"`, `"template"`, or `"input-set"` |
| `items[].yaml` | string | yes | Raw v0 YAML string |

**Response — 200 OK**

The response is always HTTP 200. Per-item errors are reported inline. The outer call only fails (non-200) if the request itself is malformed.

```json
{
  "results": [
    {
      "id": "client-assigned-id-1",
      "entity_type": "pipeline",
      "yaml": "<v1 yaml>",
      "checksum": "sha256:...",
      "error": null
    },
    {
      "id": "client-assigned-id-2",
      "entity_type": "template",
      "yaml": null,
      "checksum": null,
      "error": "failed to parse template: missing 'type' field"
    }
  ]
}
```

| Field | Type | Description |
|-------|------|-------------|
| `results[].id` | string | Echoed from request item |
| `results[].entity_type` | string | Echoed from request item |
| `results[].yaml` | string \| null | Converted v1 YAML, `null` on error |
| `results[].checksum` | string \| null | Checksum, `null` on error |
| `results[].error` | string \| null | Error message, `null` on success |

---

### 3.5 Health Check

```
GET /healthz
```

**Response — 200 OK**

```json
{
  "status": "ok"
}
```

---

## 4. Error Handling

### HTTP Status Codes

| Code | Meaning |
|------|---------|
| `200` | Success |
| `400` | Request is malformed or YAML is invalid / unparseable |
| `422` | YAML parsed successfully but conversion failed (e.g. unsupported stage type) |
| `500` | Unexpected internal error |

### Error Response Body

```json
{
  "code": "INVALID_YAML",
  "message": "failed to parse v0 pipeline: yaml: line 5: did not find expected key",
  "details": {}
}
```

| Field | Type | Description |
|-------|------|-------------|
| `code` | string | Machine-readable error code (see table below) |
| `message` | string | Human-readable explanation |
| `details` | object | Optional extra context (e.g. line numbers) |

### Error Codes

| Code | HTTP | Description |
|------|------|-------------|
| `MISSING_FIELD` | 400 | Required field absent in request JSON |
| `INVALID_JSON` | 400 | Request body is not valid JSON |
| `INVALID_YAML` | 400 | `yaml` field is not valid YAML |
| `WRONG_ENTITY_TYPE` | 400 | Top-level YAML key does not match requested entity type |
| `CONVERSION_FAILED` | 422 | Parsed successfully but converter returned nil or error |
| `UNSUPPORTED_STAGE_TYPE` | 422 | Stage type present in v0 has no v1 equivalent |
| `INTERNAL_ERROR` | 500 | Unexpected error |

---

## 5. Checksum Specification

- **Algorithm:** SHA-256
- **Input:** the UTF-8-encoded bytes of the v1 YAML string returned in the response
- **Encoding:** lowercase hex
- **Format in response:** `"sha256:" + <64-char hex string>`

Example in Go:

```go
import (
    "crypto/sha256"
    "fmt"
)

func Checksum(yamlBytes []byte) string {
    sum := sha256.Sum256(yamlBytes)
    return fmt.Sprintf("sha256:%x", sum)
}
```

The checksum is computed **after** marshalling the v1 YAML. The same bytes in `yaml` and the checksum always correspond — the service must not re-marshal between computing the checksum and writing the response.

---

## 6. Implementation Notes

### 6.1 Pipeline Converter (`service/converter/pipeline.go`)

This is a thin wrapper around the existing library — no new conversion logic required.

```
v0 YAML bytes
  → convert/harness/yaml.Parse(reader)   → *v0.Config
  → pipelineconverter.NewPipelineConverter().ConvertPipeline(&v0Config.Pipeline)  → *v1.Pipeline
  → v1.MarshalPipeline(pipeline)         → []byte (v1 YAML)
  → Checksum(bytes)                      → checksum string
```

Key packages:
- `github.com/drone/go-convert/convert/harness/yaml` — v0 parser
- `github.com/drone/go-convert/convert/v0tov1/pipeline_converter` — converter
- `github.com/drone/go-convert/convert/v0tov1/yaml` — v1 marshaller

### 6.2 Template Converter (`service/converter/template.go`) — New Work

No dedicated v0→v1 template converter exists today. The plan:

1. **Parse** the v0 template YAML. Define a minimal `v0Template` struct:

```go
type v0Template struct {
    Template struct {
        Name       string      `yaml:"name"`
        Identifier string      `yaml:"identifier"`
        Org        string      `yaml:"orgIdentifier"`
        Project    string      `yaml:"projectIdentifier"`
        Type       string      `yaml:"type"` // "Pipeline" | "Stage" | "Step"
        Spec       interface{} `yaml:"spec"`
    } `yaml:"template"`
}
```

2. **Route** based on `Type`:
   - `Pipeline` — re-serialize `Spec` as v0 pipeline YAML and run through the pipeline converter. Wrap the resulting `*v1.Pipeline` into a `v1.Schema{Template: &v1.Template{...}}`
   - `Stage` — extract the single stage from `Spec`, convert via `PipelineConverter.convertStage` (internal method — may need to export it), wrap in `v1.Template{Stage: ...}`
   - `Step` — similar, extract and convert a single step, wrap in `v1.Template{Step: ...}`

3. **Marshal** the `v1.Schema` (not just `v1.Pipeline`) using `gopkg.in/yaml.v3`.

> **Open question for reviewer:** Should Stage/Step template conversion be in scope for v1 of this service? It requires exporting some internal converter methods. A simpler initial scope could be **Pipeline templates only**.

### 6.3 Input Set Converter (`service/converter/inputset.go`) — New Work

The v0 input set YAML contains a partial pipeline snapshot used to supply runtime inputs. The v1 shape is a flat `inputset:` key with an `inputs:` map.

Conversion approach:

1. Parse the v0 `inputSet:` document.
2. Flatten `pipeline.properties.ci.codebase.build.spec.*` and any other overridden fields into a key-value `inputs` map (using a traversal or `gjson` paths).
3. Preserve `pipeline.identifier` as `pipeline:` reference field.
4. Marshal as:

```yaml
inputset:
  pipeline: <pipelineIdentifier>
  inputs:
    codebase.build.type: PR
    codebase.build.spec.number: <+trigger.prNumber>
```

> **Open question for reviewer:** The exact v1 input set schema is not fully defined in this codebase. Please confirm the target v1 shape before implementing this converter.

---

## 7. Server Setup

### 7.1 Server (`service/server.go`)

```go
type Server struct {
    httpServer *http.Server
    port       int
    logger     *slog.Logger
}

func NewServer(port int, logger *slog.Logger) *Server

func (s *Server) Start() error           // blocks until shutdown
func (s *Server) Stop(ctx context.Context) error  // graceful shutdown
```

### 7.2 Middleware chain (applied in order)

1. **Recovery** — catch panics, log, return `500 INTERNAL_ERROR`
2. **Request ID** — generate a UUID, set as `X-Request-ID` response header and in context
3. **Logger** — structured log at INFO level: method, path, status, latency, request ID

### 7.3 Router

Use `net/http` with `http.ServeMux` (Go 1.22+ pattern-matching mux). No external router dependency needed.

```
POST /api/v1/convert/pipeline   → handler.ConvertPipeline
POST /api/v1/convert/template   → handler.ConvertTemplate
POST /api/v1/convert/input-set  → handler.ConvertInputSet
POST /api/v1/convert/batch      → handler.ConvertBatch
GET  /healthz                   → handler.Healthz
```

### 7.4 Configuration (`cmd/server/main.go`)

Environment variables (with defaults):

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8090` | HTTP listen port |
| `LOG_LEVEL` | `info` | `debug` / `info` / `warn` / `error` |
| `MAX_BATCH_SIZE` | `100` | Maximum items per batch request |
| `MAX_YAML_BYTES` | `1048576` | Maximum YAML payload size (1 MB) |

### 7.5 Entrypoint (`cmd/server/main.go`)

```go
func main() {
    port, _ := strconv.Atoi(getEnv("PORT", "8090"))
    logger  := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel()}))

    srv := service.NewServer(port, logger)

    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
        <-sigCh
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        srv.Stop(ctx)
    }()

    logger.Info("starting server", "port", port)
    if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
        logger.Error("server exited", "err", err)
        os.Exit(1)
    }
}
```

---

## 8. Request / Response Types (`service/request.go`)

```go
// Single-entity request (shared by pipeline, template, input-set endpoints)
type ConvertRequest struct {
    YAML string `json:"yaml"`
}

// Single-entity response
type ConvertResponse struct {
    YAML     string `json:"yaml"`
    Checksum string `json:"checksum"`
}

// Batch request
type BatchConvertRequest struct {
    Items []BatchItem `json:"items"`
}

type BatchItem struct {
    ID         string `json:"id"`
    EntityType string `json:"entity_type"` // "pipeline" | "template" | "input-set"
    YAML       string `json:"yaml"`
}

// Batch response
type BatchConvertResponse struct {
    Results []BatchResult `json:"results"`
}

type BatchResult struct {
    ID         string  `json:"id"`
    EntityType string  `json:"entity_type"`
    YAML       *string `json:"yaml"`
    Checksum   *string `json:"checksum"`
    Error      *string `json:"error"`
}

// Error response
type ErrorResponse struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}
```

---

## 9. Dependencies

No new third-party dependencies are required beyond what is already in `go.mod`. The service uses only:

- `net/http` (stdlib) — HTTP server
- `encoding/json` (stdlib) — JSON marshal/unmarshal
- `crypto/sha256` (stdlib) — checksum
- `log/slog` (stdlib, Go 1.21+) — structured logging
- Existing `go-convert` library packages

> **Note:** Go 1.19 is declared in `go.mod`. `log/slog` was added in Go 1.21. Either bump the `go` directive to `1.21` or use `go.uber.org/zap` (already available transitively). This is a decision for the reviewer.

---

## 10. Deployment

### Dockerfile (`Dockerfile.service`)

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/go-convert-service ./cmd/server

FROM gcr.io/distroless/static-debian12
COPY --from=builder /bin/go-convert-service /go-convert-service
EXPOSE 8090
ENTRYPOINT ["/go-convert-service"]
```

### Health probe

```
GET /healthz → 200 {"status":"ok"}
```

---

## 11. Open Questions / Decisions Needed

| # | Question | Options | Owner |
|---|----------|---------|-------|
| 1 | gRPC vs REST? | REST (proposed here) is simpler for YAML payloads; gRPC (like SCM service) is more consistent with Harness platform. | Reviewer |
| 2 | Template scope for v1? | Pipeline templates only (simple), or include Stage + Step templates (requires exporting internal converter methods)? | Reviewer |
| 3 | Input set v1 schema? | Exact target YAML shape for v1 input set needs confirmation from product/spec. | Reviewer |
| 4 | Go version bump? | `go.mod` is on Go 1.19; `log/slog` needs 1.21. Bump or use zap? | Reviewer |
| 5 | Auth? | No-auth + API gateway, static token via header, or mTLS? | Reviewer |
| 6 | Max payload size? | 1 MB default — is this sufficient for large pipelines? | Reviewer |
| 7 | Checksum of v1 YAML vs checksum of input v0 YAML? | Spec proposes checksum of the **output** v1 YAML. Should input checksum also be returned for idempotency checks? | Reviewer |

---

## 12. Sequence Diagram — Single Pipeline Conversion

```
Client                    Server                     Library
  │                          │                           │
  │  POST /api/v1/convert/   │                           │
  │  pipeline {yaml: "..."}  │                           │
  │─────────────────────────>│                           │
  │                          │  json.Unmarshal(body)     │
  │                          │  validate(req.YAML)       │
  │                          │                           │
  │                          │  v0.Parse(yaml)──────────>│
  │                          │<─────────── *v0.Config ───│
  │                          │                           │
  │                          │  converter.ConvertPipeline│
  │                          │  (&v0Config.Pipeline)────>│
  │                          │<──────── *v1.Pipeline ────│
  │                          │                           │
  │                          │  v1.MarshalPipeline(p)───>│
  │                          │<────────── []byte ────────│
  │                          │                           │
  │                          │  Checksum(bytes)          │
  │                          │                           │
  │  200 {yaml, checksum}    │                           │
  │<─────────────────────────│                           │
```

---

*End of Tech Spec — please add comments and modify §11 decisions before implementation begins.*

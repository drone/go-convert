# Expression Conversion API

Convert Harness v0 expressions to v1 format. Useful for converting individual expressions outside of a full pipeline conversion, or for understanding how specific expressions transform.

## Endpoint

```
POST /api/v1/convert/expression
```

**Default HTTP port:** `8092`

## Request Format

There are two ways to provide context for expression conversion:

### Option 1: Pipeline YAML (recommended)

Pass the raw v0 pipeline YAML and the server automatically derives the step-type map,
v1 path map, and enables FQN mode — the same way pipeline, template, input-set, and
trigger conversions build context internally.

```json
{
  "expression": "<+pipeline.stages.build.spec.execution.steps.step1.output>",
  "context_pipeline_yaml": "pipeline:\n  name: my-pipeline\n  stages:\n    - stage:\n        identifier: build\n        type: CI\n        spec:\n          execution:\n            steps:\n              - step:\n                  identifier: step1\n                  type: Run\n                  spec:\n                    command: echo hello\n"
}
```

### Option 2: Manual context fields

Explicitly supply the step-type map and other context fields. This is useful when
you don't have the full pipeline YAML or need fine-grained control.

```json
{
  "expression": "<+pipeline.stages.build.spec.execution.steps.step1.output>",
  "context": {
    "current_step_id": "step1",
    "current_step_type": "Run",
    "current_step_v1_path": "pipeline.stages.build.steps.step1",
    "step_type_map": {
      "step1": "Run",
      "step2": "Action"
    },
    "step_v1_path_map": {
      "step1": "pipeline.stages.build.steps.step1",
      "step2": "pipeline.stages.build.steps.step2"
    },
    "use_fqn": true
  }
}
```

> **Note:** When `context_pipeline_yaml` is provided, the `context` field is ignored.

### Request Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `expression` | string | One of `expression`, `expressions`, or `remote_file` required | A single v0 expression to convert |
| `expressions` | string[] | One of `expression`, `expressions`, or `remote_file` required | Multiple v0 expressions to convert |
| `remote_file` | string | One of `expression`, `expressions`, or `remote_file` required | Raw contents of a remote file (manifest, values.yaml, config, etc.) with embedded `<+...>` expressions. All expressions are converted in place. |
| `context_pipeline_yaml` | string | Optional | Raw v0 pipeline YAML; server derives context automatically (recommended) |
| `context` | object | Optional | Manual context for conversion (ignored when `context_pipeline_yaml` is provided) |

### Context Fields

| Field | Type | Description |
|-------|------|-------------|
| `current_step_id` | string | ID of the current step (when converting expressions inside a step) |
| `current_step_type` | string | Type of the current step (e.g., "Run", "Action", "Plugin") |
| `current_step_v1_path` | string | V1 FQN base path to the current step |
| `step_type_map` | map[string]string | Maps step IDs to their types for all steps in the pipeline |
| `step_v1_path_map` | map[string]string | Maps step IDs to their v1 FQN base paths |
| `use_fqn` | boolean | Enable FQN mode for step expressions |

## Response Format

### Single Expression Response

```json
{
  "expression": "<+pipeline.stages.build.steps.step1.output>",
  "checksum": "sha256:abc123..."
}
```

### Multiple Expressions Response

```json
{
  "expressions": {
    "<+pipeline.stages.build.spec.execution.steps.step1.output>": "<+pipeline.stages.build.steps.step1.output>",
    "<+stage.spec.execution.steps.step2.output>": "<+stage.steps.step2.output>"
  },
  "checksum": "sha256:def456..."
}
```

### Remote File Response

```json
{
  "remote_file": "apiVersion: apps/v1\nmetadata:\n  name: <+pipeline.inputs.appName>\n  namespace: <+pipeline.stages.deploy.steps.deploy1.namespace>\n",
  "checksum": "sha256:ghi789..."
}
```

## Examples

### Basic Expression Conversion

Convert a simple expression without context:

```bash
curl -X POST http://localhost:8092/api/v1/convert/expression \
  -H "Content-Type: application/json" \
  -d '{
    "expression": "<+pipeline.stages.build.spec.execution.steps.step1.output>"
  }'
```

Response:
```json
{
  "expression": "<+pipeline.stages.build.steps.step1.output>",
  "checksum": "sha256:abc123..."
}
```

### Context-Aware Conversion (Relative)

Convert a step-relative expression with step type context:

```bash
curl -X POST http://localhost:8092/api/v1/convert/expression \
  -H "Content-Type: application/json" \
  -d '{
    "expression": "<+step.spec.command>",
    "context": {
      "current_step_type": "Run"
    }
  }'
```

Response:
```json
{
  "expression": "<+step.spec.script>",
  "checksum": "sha256:abc123"
}
```

### Context-Aware Conversion (FQN Mode)

Convert to fully qualified names when `use_fqn` is enabled:

```bash
curl -X POST http://localhost:8092/api/v1/convert/expression \
  -H "Content-Type: application/json" \
  -d '{
    "expression": "<+step.spec.command>",
    "context": {
      "current_step_type": "Run",
      "current_step_v1_path": "pipeline.stages.build.steps.runStep1",
      "use_fqn": true
    }
  }'
```

Response:
```json
{
  "expression": "<+pipeline.stages.build.steps.runStep1.spec.script>",
  "checksum": "sha256:abc123"
}
```

### Batch Expression Conversion

Convert multiple expressions at once:

```bash
curl -X POST http://localhost:8092/api/v1/convert/expression \
  -H "Content-Type: application/json" \
  -d '{
    "expressions": [
      "<+pipeline.stages.build.spec.execution.steps.step1.output>",
      "<+stage.spec.execution.steps.step2.output>",
      "<+pipeline.variables.myVar>"
    ]
  }'
```

Response:
```json
{
  "expressions": {
    "<+pipeline.stages.build.spec.execution.steps.step1.output>": "<+pipeline.stages.build.steps.step1.output>",
    "<+stage.spec.execution.steps.step2.output>": "<+stage.steps.step2.output>",
    "<+pipeline.variables.myVar>": "<+pipeline.variables.myVar>"
  },
  "checksum": "sha256:abc123"
}
```

### Remote File Conversion

Convert all expressions embedded in a remote file (manifest, values.yaml, config, etc.):

```bash
curl -X POST http://localhost:8092/api/v1/convert/expression \
  -H "Content-Type: application/json" \
  -d '{
    "remote_file": "apiVersion: apps/v1\nmetadata:\n  name: <+pipeline.variables.appName>\n  image: <+pipeline.stages.build.spec.execution.steps.build1.output>\n"
  }'
```

Response:
```json
{
  "remote_file": "apiVersion: apps/v1\nmetadata:\n  name: <+pipeline.variables.appName>\n  image: <+pipeline.stages.build.steps.build1.output>\n",
  "checksum": "sha256:abc123"
}
```

### Remote File with Pipeline YAML Context

For step-type-aware and FQN conversion inside a remote file, pass the pipeline YAML:

```bash
PIPELINE_YAML=$(cat my_v0_pipeline.yaml)

curl -X POST http://localhost:8092/api/v1/convert/expression \
  -H "Content-Type: application/json" \
  -d "$(jq -n \
    --arg file "$(cat manifest.yaml)" \
    --arg yaml "$PIPELINE_YAML" \
    '{remote_file: $file, pipeline_yaml: $yaml}')"
```

### Pipeline YAML Context Example

Instead of manually constructing the context, pass the full v0 pipeline YAML and
let the server derive the context automatically:

```bash
# Read the pipeline YAML from a file
PIPELINE_YAML=$(cat my_v0_pipeline.yaml)

curl -X POST http://localhost:8092/api/v1/convert/expression \
  -H "Content-Type: application/json" \
  -d "$(jq -n \
    --arg expr '<+pipeline.stages.build.spec.execution.steps.step1.output>' \
    --arg yaml "$PIPELINE_YAML" \
    '{expression: $expr, context_pipeline_yaml: $yaml}')"
```

Response:
```json
{
  "expression": "<+pipeline.stages.build.steps.step1.output>",
  "checksum": "sha256:abc123"
}
```

### Cross-Step Reference Example

When converting expressions that reference other steps (using `steps.STEPID`):

```bash
curl -X POST http://localhost:8092/api/v1/convert/expression \
  -H "Content-Type: application/json" \
  -d '{
    "expression": "<+steps.otherStep.spec.command>",
    "context": {
      "current_step_id": "currentStep",
      "current_step_type": "Run",
      "current_step_v1_path": "pipeline.stages.build.steps.currentStep",
      "step_type_map": {
        "currentStep": "Run",
        "otherStep": "Run"
      },
      "step_v1_path_map": {
        "currentStep": "pipeline.stages.build.steps.currentStep",
        "otherStep": "pipeline.stages.build.steps.otherStep"
      },
      "use_fqn": true
    }
  }'
```

## gRPC API

The expression conversion is also available via gRPC.

**Default gRPC port:** `8090`

### Service Definition

```protobuf
service GoConvertService {
  rpc ConvertExpression(ExpressionConvertRequest) returns (ExpressionConvertResponse);
}

message ExpressionConvertRequest {
  string expression = 1;
  repeated string expressions = 2;
  string context_pipeline_yaml = 3;
  ExpressionContext context = 4;
  string remote_file = 5;
}

message ExpressionContext {
  string current_step_id = 1;
  string current_step_type = 2;
  string current_step_v1_path = 3;
  map<string, string> step_type_map = 4;
  map<string, string> step_v1_path_map = 5;
  bool use_fqn = 6;
}

message ExpressionConvertResponse {
  string expression = 1;
  map<string, string> expressions = 2;
  string checksum = 3;
  string remote_file = 4;
}
```

### gRPC Example (grpcurl)

```bash
grpcurl -plaintext -d '{
  "expression": "<+pipeline.variables.myVar>"
}' localhost:8090 io.harness.pms.conversion.proto.GoConvertService/ConvertExpression
```

## Client Script Usage

The `convert_client.py` script supports expression conversion via `--expression`, `--expressions`, and `--remote-file`:

```bash
# Single expression (HTTP)
python convert_client.py --expression '<+pipeline.variables.foo>'

# Multiple expressions (HTTP)
python convert_client.py --expressions '<+pipeline.variables.foo>' '<+stage.spec.execution.steps.s1.output>'

# Remote file — convert all expressions in a manifest/config file
python convert_client.py --remote-file manifest.yaml

# Remote file with pipeline YAML context for FQN resolution
python convert_client.py --remote-file manifest.yaml --context-pipeline my_v0_pipeline.yaml

# With pipeline YAML context for FQN resolution
python convert_client.py --expression '<+pipeline.stages.build.spec.execution.steps.step1.output>' \
    --context-pipeline my_v0_pipeline.yaml

# Via gRPC
python convert_client.py --grpc --expression '<+pipeline.variables.foo>'

# Via gRPC with remote file
python convert_client.py --grpc --remote-file manifest.yaml --context-pipeline my_v0_pipeline.yaml
```

## Common Expression Conversions

### Path Structure Changes

| V0 Expression | V1 Expression | Notes |
|---------------|---------------|-------|
| `<+pipeline.stages.STAGE.spec.execution.steps.STEP.*>` | `<+pipeline.stages.STAGE.steps.STEP.*>` | Removes `spec.execution` |
| `<+stage.spec.execution.steps.STEP.*>` | `<+stage.steps.STEP.*>` | Removes `spec.execution` |
| `<+pipeline.stages.STAGE.spec.execution.rollbackSteps.STEP.*>` | `<+pipeline.stages.STAGE.rollback.STEP.*>` | Rollback steps |

### Step-Type Specific Conversions (requires `current_step_type` context)

| V0 Expression | V1 Expression | Step Type |
|---------------|---------------|-----------|
| `<+step.spec.command>` | `<+step.spec.script>` | Run |
| `<+step.spec.image>` | `<+step.spec.container.image>` | Run |
| `<+step.spec.shell>` | `<+step.spec.shell>` | Run |
| `<+step.spec.envVariables.X>` | `<+step.spec.env.X>` | Run |

## Programmatic Usage

Use the expression conversion directly in Go code:

```go
package main

import (
    "fmt"
    "github.com/drone/go-convert/service/converter"
)

func main() {
    // Simple conversion without context
    result := converter.ConvertExpression(
        "<+pipeline.stages.build.spec.execution.steps.step1.output>",
        nil,
    )
    fmt.Println(result) // <+pipeline.stages.build.steps.step1.output>

    // Automatic context from pipeline YAML (recommended)
    pipelineYAML := `pipeline:
  name: my-pipeline
  stages:
    - stage:
        identifier: build
        type: CI
        spec:
          execution:
            steps:
              - step:
                  identifier: step1
                  type: Run
                  spec:
                    command: echo hello`
    result = converter.ConvertExpressionWithPipeline(
        "<+pipeline.stages.build.spec.execution.steps.step1.output>",
        pipelineYAML,
    )
    fmt.Println(result) // <+pipeline.stages.build.steps.step1.output>

    // Context-aware conversion (relative, manual context)
    ctx := &converter.ExpressionContext{
        CurrentStepType: "Run",
    }
    result = converter.ConvertExpression("<+step.spec.command>", ctx)
    fmt.Println(result) // <+step.spec.script>

    // Context-aware conversion (FQN mode, manual context)
    ctx = &converter.ExpressionContext{
        CurrentStepType:   "Run",
        CurrentStepV1Path: "pipeline.stages.build.steps.runStep1",
        UseFQN:            true,
    }
    result = converter.ConvertExpression("<+step.spec.command>", ctx)
    fmt.Println(result) // <+pipeline.stages.build.steps.runStep1.spec.script>

    // Batch conversion with pipeline YAML context
    expressions := []string{
        "<+pipeline.variables.myVar>",
        "<+stage.spec.execution.steps.step1.output>",
    }
    results := converter.ConvertExpressionsWithPipeline(expressions, pipelineYAML)
    for orig, converted := range results {
        fmt.Printf("%s -> %s\n", orig, converted)
    }
}
```

## Error Responses

### Missing Field Error

```json
{
  "code": "MISSING_FIELD",
  "message": "either 'expression' or 'expressions' field is required"
}
```

### Invalid JSON Error

```json
{
  "code": "INVALID_JSON",
  "message": "unexpected EOF"
}
```

## Notes

1. **`pipeline_yaml` is the recommended approach**: Pass the raw v0 pipeline YAML and the server automatically derives all context (step types, v1 paths, FQN mode). This is the same mechanism used by pipeline, template, input-set, and trigger conversions.

2. **Context is optional**: Basic path conversions work without context (e.g., `spec.execution.steps` → `steps`). Context is only needed for step-type-specific field conversions.

3. **`pipeline_yaml` supersedes `context`**: When both are provided, `pipeline_yaml` takes precedence and `context` is ignored.

4. **Step type resolution**: For step-specific field conversions (like `spec.command` → `spec.script` for Run steps), provide:
   - `current_step_type` — for expressions starting with `step.`
   - `step_type_map` — for expressions referencing other steps via `steps.STEPID`

5. **FQN mode**: When `use_fqn: true`:
   - Relative expressions (`step.spec.X`) become fully qualified (`pipeline.stages.STAGE.steps.STEP.spec.X`)
   - Requires `current_step_v1_path` for the current step
   - Requires `step_v1_path_map` for cross-step references

6. **Non-expression strings**: Input without `<+` markers is returned unchanged.
# Expression Conversion API

Convert Harness v0 expressions to v1 format. Useful for converting individual expressions outside of a full pipeline conversion, or for understanding how specific expressions transform.

## Endpoint

```
POST /api/v1/convert/expression
```

**Default HTTP port:** `8092`

## Request Format

Expression conversion is **FQN-only**: context is supplied as a **v1** pipeline
YAML (`context_pipeline_yaml`) which the server walks into the FQN-keyed step
lookup, plus an optional `current_fqn` that anchors the call-site for
step-group chaining and step-type-aware resolution. There are no manual step maps.

```json
{
  "expression": "<+step.spec.command>",
  "context_pipeline_yaml": "pipeline:\n  identifier: myPipeline\n  stages:\n    - stage:\n        identifier: build\n        steps:\n          - step:\n              identifier: compile\n              type: Run\n              spec:\n                shell: Sh\n                run: go build ./...\n",
  "current_fqn": "pipeline.stages.build.steps.compile"
}
```

### Request Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `expression` | string | One of `expression`, `expressions`, or `remote_file` required | A single v0 expression to convert |
| `expressions` | string[] | One of `expression`, `expressions`, or `remote_file` required | Multiple v0 expressions to convert |
| `remote_file` | string | One of `expression`, `expressions`, or `remote_file` required | Raw contents of a remote file (manifest, values.yaml, config, etc.) with embedded `<+...>` expressions. All expressions are converted in place. |
| `context_pipeline_yaml` | string | Optional | A **v1** pipeline YAML. The server walks it into the FQN-keyed step lookup and resolves step references against it in FQN mode. |
| `current_fqn` | string | Optional | Full v1 FQN of the step the expression lives in, e.g. `pipeline.stages.prod.steps.G1.steps.deploy`. Used with `context_pipeline_yaml` to derive the call-site stage, enclosing step-group chain, and current step type (for `step.spec.*` resolution). |

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

### Context-Aware Conversion (FQN Mode)

Convert a step-relative `step.spec.*` expression to a fully qualified name. The
server looks up `current_fqn` in the v1 `context_pipeline_yaml` to recover the
step type, so the type-specific field mapping (`command` → `script` for Run) is
applied automatically:

```bash
curl -X POST http://localhost:8092/api/v1/convert/expression \
  -H "Content-Type: application/json" \
  -d "$(jq -n \
    --arg expr '<+step.spec.command>' \
    --arg yaml "$(cat my_v1_pipeline.yaml)" \
    --arg fqn 'pipeline.stages.build.steps.runStep1' \
    '{expression: $expr, context_pipeline_yaml: $yaml, current_fqn: $fqn}')"
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

For step-type-aware and FQN conversion inside a remote file, pass the **v1**
pipeline YAML (and optionally `current_fqn` to anchor `step.*` self-references):

```bash
PIPELINE_YAML=$(cat my_v1_pipeline.yaml)

curl -X POST http://localhost:8092/api/v1/convert/expression \
  -H "Content-Type: application/json" \
  -d "$(jq -n \
    --arg file "$(cat manifest.yaml)" \
    --arg yaml "$PIPELINE_YAML" \
    '{remote_file: $file, context_pipeline_yaml: $yaml}')"
```

### Pipeline YAML Context Example

Pass the full **v1** pipeline YAML and let the server derive the FQN context
automatically:

```bash
# Read the v1 pipeline YAML from a file
PIPELINE_YAML=$(cat my_v1_pipeline.yaml)

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

When converting an expression that references another step (`steps.STEPID`),
supply the v1 `context_pipeline_yaml` and the `current_fqn` of the call-site
step. The referenced step is resolved against the FQN-keyed lookup:

```bash
curl -X POST http://localhost:8092/api/v1/convert/expression \
  -H "Content-Type: application/json" \
  -d "$(jq -n \
    --arg expr '<+steps.otherStep.spec.command>' \
    --arg yaml "$(cat my_v1_pipeline.yaml)" \
    --arg fqn 'pipeline.stages.build.steps.currentStep' \
    '{expression: $expr, context_pipeline_yaml: $yaml, current_fqn: $fqn}')"
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
  string context_pipeline_yaml = 3;  // v1 pipeline YAML
  string current_fqn = 4;            // call-site v1 FQN
  string remote_file = 5;
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
| `<+...template.templateInputs.*>` | `<+...template.with.overlay.*>` | Template inputs at pipeline / stage / stepGroup / step level (e.g. `<+step.template.templateInputs>` → `<+step.template.with.overlay>`) |

### Step-Type Specific Conversions (requires `context_pipeline_yaml` + `current_fqn`)

These `<+step.*>` self-references are resolved by deriving the current step's
type from `current_fqn` looked up in the FQN-keyed step map built from
`context_pipeline_yaml`. The `step` prefix is rewritten to the step's full v1
FQN. For example, with `current_fqn = pipeline.stages.build.steps.compile`
(a Run step), `<+step.spec.command>` becomes
`<+pipeline.stages.build.steps.compile.spec.script>`.

| V0 Expression | V1 Expression (relative) | Step Type |
|---------------|--------------------------|-----------|
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

    // Automatic context from a v1 pipeline YAML
    pipelineYAML := `pipeline:
  identifier: myPipeline
  stages:
    - stage:
        identifier: build
        steps:
          - step:
              identifier: step1
              type: Run
              spec:
                shell: Sh
                run: echo hello`
    result = converter.ConvertExpressionWithPipeline(
        "<+pipeline.stages.build.spec.execution.steps.step1.output>",
        pipelineYAML,
    )
    fmt.Println(result) // <+pipeline.stages.build.steps.step1.output>

    // FQN-mode conversion: supply the v1 context pipeline + call-site FQN.
    // The step type is recovered from the FQN lookup, so step.spec.command
    // becomes step.spec.script for a Run step.
    ctx := &converter.ExpressionContext{
        ContextPipelineYAML: pipelineYAML,
        CurrentFQN:          "pipeline.stages.build.steps.step1",
    }
    result = converter.ConvertExpression("<+step.spec.command>", ctx)
    fmt.Println(result) // <+pipeline.stages.build.steps.step1.spec.script>

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

1. **`context_pipeline_yaml` must be a v1 pipeline**: The server parses it as a
   v1 pipeline and walks it into the FQN-keyed step lookup. A v0 pipeline (or an
   unparseable document) is ignored and conversion falls back to structural-only.

2. **Context is optional**: Basic path conversions work without context (e.g.,
   `spec.execution.steps` → `steps`). `context_pipeline_yaml` is only needed for
   step-type-specific field conversions and FQN resolution.

3. **Step type is derived, not passed**: For step-specific field conversions
   (like `spec.command` → `spec.script` for Run steps), supply `current_fqn`
   pointing at the call-site step within `context_pipeline_yaml`; the server
   recovers the step type from the FQN lookup. There are no manual step maps.

4. **FQN mode** is enabled automatically whenever `context_pipeline_yaml` yields
   a usable step lookup. Relative expressions (`step.spec.X`) then become fully
   qualified (`pipeline.stages.STAGE.steps.STEP.spec.X`); `current_fqn` anchors
   the call-site for `step.*` self-references and step-group chaining.

5. **Non-expression strings**: Input without `<+` markers is returned unchanged.
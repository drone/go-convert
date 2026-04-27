# Expression Conversion API

Convert Harness v0 expressions to v1 format. Useful for converting individual expressions outside of a full pipeline conversion, or for understanding how specific expressions transform.

## Endpoint

```
POST /api/v1/convert/expression
```

**Default HTTP port:** `8092`

## Request Format

```json
{
  "expression": "<+pipeline.stages.build.spec.execution.steps.step1.output>",
  "context": {
    "step_id": "step1",
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

### Request Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `expression` | string | One of `expression` or `expressions` required | A single v0 expression to convert |
| `expressions` | string[] | One of `expression` or `expressions` required | Multiple v0 expressions to convert |
| `context` | object | Optional | Context for context-aware conversion |

### Context Fields

| Field | Type | Description |
|-------|------|-------------|
| `step_id` | string | ID of the current step (when converting expressions inside a step) |
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

## Examples

### Basic Expression Conversion

Convert a simple expression without context:

```bash
curl -X POST http://localhost:8090/api/v1/convert/expression \
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
curl -X POST http://localhost:8090/api/v1/convert/expression \
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
curl -X POST http://localhost:8090/api/v1/convert/expression \
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
curl -X POST http://localhost:8090/api/v1/convert/expression \
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

### Cross-Step Reference Example

When converting expressions that reference other steps (using `steps.STEPID`):

```bash
curl -X POST http://localhost:8090/api/v1/convert/expression \
  -H "Content-Type: application/json" \
  -d '{
    "expression": "<+steps.otherStep.spec.command>",
    "context": {
      "step_id": "currentStep",
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

    // Context-aware conversion (relative)
    ctx := &converter.ExpressionContext{
        CurrentStepType: "Run",
    }
    result = converter.ConvertExpression("<+step.spec.command>", ctx)
    fmt.Println(result) // <+step.spec.script>

    // Context-aware conversion (FQN mode)
    ctx = &converter.ExpressionContext{
        CurrentStepType:   "Run",
        CurrentStepV1Path: "pipeline.stages.build.steps.runStep1",
        UseFQN:            true,
    }
    result = converter.ConvertExpression("<+step.spec.command>", ctx)
    fmt.Println(result) // <+pipeline.stages.build.steps.runStep1.spec.script>

    // Batch conversion
    expressions := []string{
        "<+pipeline.variables.myVar>",
        "<+stage.spec.execution.steps.step1.output>",
    }
    results := converter.ConvertExpressions(expressions, nil)
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

1. **Context is optional**: Basic path conversions work without context (e.g., `spec.execution.steps` → `steps`). Context is only needed for step-type-specific field conversions.

2. **Step type resolution**: For step-specific field conversions (like `spec.command` → `spec.script` for Run steps), provide:
   - `current_step_type` — for expressions starting with `step.`
   - `step_type_map` — for expressions referencing other steps via `steps.STEPID`

3. **FQN mode**: When `use_fqn: true`:
   - Relative expressions (`step.spec.X`) become fully qualified (`pipeline.stages.STAGE.steps.STEP.spec.X`)
   - Requires `current_step_v1_path` for the current step
   - Requires `step_v1_path_map` for cross-step references

4. **Non-expression strings**: Input without `<+` markers is returned unchanged.
# Skill: Making v0→v1 Converter Changes from YAML Examples

## Overview
This skill guides adding or modifying step/stage converters in the Harness v0→v1 pipeline converter. The user provides a **v0 YAML snippet** and the **desired v1 YAML output**; the assistant derives the struct and converter changes needed.

---

## 1. How the User Provides YAML

The user supplies two YAML blocks — v0 (input) and v1 (expected output).

**Example prompt:**
```
Convert this v0 step to v1:

v0:
- step:
    type: MyNewStep
    identifier: step1
    name: My Step
    timeout: 10m
    spec:
      connectorRef: myConnector
      region: us-east-1
      someBool: true

v1:
- id: step1
  name: My Step
  timeout: 10m
  template:
    uses: myNewStep@1.0.0
    with:
      connector: myConnector
      region: us-east-1
      some_bool: true
```

If only v0 is provided, the assistant should ask the user for the expected v1 output before proceeding.

---

## 2. Codebase Reference (Concise)

### 2.1 Package Layout

| Package | Import alias | Purpose |
|---|---|---|
| `convert/harness/yaml` | `v0` | v0 structs: `Pipeline`, `Stage`, `Step`, step specs (`StepRun`, `StepPlugin`, etc.) |
| `convert/v0tov1/yaml` | `v1` | v1 structs: `Pipeline`, `Stage`, `Step`, `StepRun`, `StepAction`, `StepTemplate`, etc. |
| `convert/v0tov1/convert_helpers` | `convert_helpers` | Per-step converter functions (`ConvertStepRun`, `ConvertStepK8sApply`, etc.) |
| `convert/v0tov1/pipeline_converter` | `pipelineconverter` | Orchestrator: `ConvertSingleStep` switch, `convertCommonStepSettings` |
| `internal/flexible` | `flexible` | `Field[T]` — holds either a typed value or a Harness expression string |

### 2.2 v0 Step Struct (input)

```
convert/harness/yaml/step.go
```

Every v0 step has a generic wrapper:
```go
type Step struct {
    ID       string      `json:"identifier,omitempty"`
    Name     string      `json:"name,omitempty"`
    Type     string      `json:"type,omitempty"`
    Spec     interface{} `json:"spec,omitempty"`       // concrete type set by UnmarshalJSON
    Timeout  string      `json:"timeout"`
    Env      *flexible.Field[map[string]string]  `json:"envVariables,omitempty"`
    When     *flexible.Field[StepWhen]           `json:"when,omitempty"`
    Strategy *Strategy                           `json:"strategy,omitempty"`
    FailureStrategies *flexible.Field[[]*FailureStrategy] `json:"failureStrategies,omitempty"`
    Template *StepTemplate                       `json:"template,omitempty"`
}
```

Step spec structs embed `CommonStepSpec` (provides `DelegateSelectors`, `IncludeInfraSelectors`):
```go
type CommonStepSpec struct {
    IncludeInfraSelectors *flexible.Field[bool]     `json:"includeInfraSelectors,omitempty"`
    DelegateSelectors     *flexible.Field[[]string] `json:"delegateSelectors,omitempty"`
}
```

**Step type constants:** `convert/harness/yaml/const.go`  
**Step UnmarshalJSON switch:** `convert/harness/yaml/step.go` (maps `Type` string → spec struct)

### 2.3 v1 Step Struct (output)

```
convert/v0tov1/yaml/step.go
```

```go
type Step struct {
    Action     *StepAction            `json:"action,omitempty"`
    Approval   *StepApproval          `json:"approval,omitempty"`
    Background *StepRun               `json:"background,omitempty"`
    Barrier    *StepBarrier           `json:"barrier,omitempty"`
    Delegate   *flexible.Field[*Delegate] `json:"delegate,omitempty"`
    Env        *flexible.Field[map[string]string] `json:"env,omitempty"`
    Group      *StepGroup             `json:"group,omitempty"`
    Id         string                 `json:"id,omitempty"`
    If         string                 `json:"if,omitempty"`
    Inputs     map[string]*Input      `json:"inputs,omitempty"`
    Name       string                 `json:"name,omitempty"`
    OnFailure  *flexible.Field[[]*FailureStrategy] `json:"on-failure,omitempty"`
    Parallel   *StepGroup             `json:"parallel,omitempty"`
    Queue      *StepQueue             `json:"queue,omitempty"`
    Run        *StepRun               `json:"run,omitempty"`
    RunTest    *StepTest              `json:"run-test,omitempty"`
    Strategy   *Strategy              `json:"strategy,omitempty"`
    Template   *StepTemplate          `json:"template,omitempty"`
    Timeout    string                 `json:"timeout,omitempty"`
    Wait       *StepWait              `json:"wait,omitempty"`
}
```

**Key v1 sub-structs:**
- `StepRun` — `run:` block (script, container, env, shell, report, outputs)
- `StepAction` — `action:` block (uses, with, env)
- `StepTemplate` — `template:` block (uses, with, env) — **most common for new steps**
- `StepApproval`, `StepBarrier`, `StepQueue`, `StepWait` — specialized

### 2.4 flexible.Field[T]

Defined in `internal/flexible/field.go`. A generic wrapper for fields that can hold either a **typed value** or a **Harness expression string** (e.g. `<+input>`, `<+pipeline.variables.x>`).

```go
type Field[T any] struct { Value interface{} }
```

**Key methods:**
- `AsStruct() (T, bool)` — get typed value
- `AsString() (string, bool)` — get string/expression value
- `IsExpression() bool` — true if Value is string containing `<+`
- `Set(value T)` / `SetString(value string)` — set value
- `IsNil() bool` — nil check

#### When to Use flexible.Field in Structs

Use `*flexible.Field[T]` for **all non-string fields** that could receive a Harness expression instead of their normal type. This includes bools, numbers, arrays, maps, and nested structs.

**Rule of thumb:** If a YAML field's Go type is not `string` and the field could be set to `<+input>` or any `<+...>` expression, wrap it in `*flexible.Field[T]`. String fields don't need wrapping because an expression string is already a valid string.

| Go type | Wrap? | Declaration |
|---|---|---|
| `string` | No | `Field string` — expressions are already strings |
| `bool` | Yes | `Field *flexible.Field[bool]` |
| `int` / `int64` | Yes | `Field *flexible.Field[int]` |
| `[]string` | Yes | `Field *flexible.Field[[]string]` |
| `map[string]string` | Yes | `Field *flexible.Field[map[string]string]` |
| custom struct `T` | Yes | `Field *flexible.Field[T]` or `Field flexible.Field[T]` |

#### Usage in v0 Structs

v0 step spec structs (`convert/harness/yaml/step.go`) use `*flexible.Field[T]` for fields that Harness UI allows expressions on:
```go
type StepBuildAndPushDockerRegistry struct {
    CommonStepSpec
    ConnectorRef string                             // string — no wrapper needed
    Tags         *flexible.Field[[]string]          // array — wrapped
    Caching      *flexible.Field[bool]              // bool — wrapped
    Labels       *flexible.Field[map[string]string] // map — wrapped
    Optimize     *flexible.Field[bool]              // bool — wrapped
    RunAsUser    *flexible.Field[int]               // int — wrapped
}
```

#### Usage in v1 Structs

v1 structs (`convert/v0tov1/yaml/`) use the same pattern for fields that support expressions:
```go
type Container struct {
    Privileged *flexible.Field[bool]       // bool — wrapped
    User       *flexible.Field[int]        // int — wrapped
    Args       *flexible.Field[[]string]   // array — wrapped
    Entrypoint *flexible.Field[[]string]   // array — wrapped
    Image      string                      // string — no wrapper
}
```

#### FlexibleField Handling in Converters

The handling depends on the **v0 field type** vs the **v1 target type**:

**Case 1: v0 `*flexible.Field[T]` → v1 `with` map value of same type T**
Don't unwrap — pass the FlexibleField directly. It marshals correctly as either the value or expression string.
```go
// EXAMPLE: spec.Optimize is *flexible.Field[bool] and template input "optimize" is bool
if spec.Optimize != nil {
    with["optimize"] = spec.Optimize
}
```

**Case 2: v0 `*flexible.Field[T]` → v1 `with` map value of different type (e.g., string)**
Unwrap and convert. Handle both expression and struct cases.
```go
// EXAMPLE: spec.RunAsUser is *flexible.Field[int] and template input "run_as_user" is string
if spec.RunAsUser != nil {
    if user, ok := spec.RunAsUser.AsString(); ok {
        with["run_as_user"] = user
    } else if user, ok := spec.RunAsUser.AsStruct(); ok {
        with["run_as_user"] = fmt.Sprintf("%d", user)
    }
}
```

**Case 3: v0 `*flexible.Field[T]` → v1 struct field of type `*flexible.Field[T]`**
Pass through directly (same type on both sides).
```go
// EXAMPLE: v0 Privileged *flexible.Field[bool] → v1 Container.Privileged *flexible.Field[bool]
container.Privileged = spec.Privileged
```

**General rule:** Match the field type/format/allowed values of the v1 target. If a template input is required and has no v0 equivalent:
- If a default value is defined in the template, set it
- Otherwise, note it in a comment in the converter

### 2.5 Converter Dispatcher

```
convert/v0tov1/pipeline_converter/convert_steps.go
```

`ConvertSingleStep` is the central switch:
```go
switch src.Type {
case v0.StepTypeRun:
    step.Run = convert_helpers.ConvertStepRun(src)
case v0.StepTypeK8sApply:
    step.Template = convert_helpers.ConvertStepK8sApply(src)
// ... etc
}
```

After the switch, `convertCommonStepSettings(src, step)` handles common step-level fields automatically (see [Section 5.1](#51-automatic-mappings-handled-by-convertcommonstepsettings)).

---

## 3. v1 Output Patterns

### Pattern A: `run` (script-based steps)
Steps that produce a `run:` block (Run, ShellScript, Action, Plugin, HTTP, Background, Container):
```yaml
- id: step1
  name: My Step
  run:
    script: "echo hello"
    container:
      image: alpine
    shell: bash
    env:
      KEY: value
```
Converter returns `*v1.StepRun`.

### Pattern B: `template` (most CI/CD steps)
Steps converted to template references (BuildAndPush*, K8s*, Helm*, Cache*, Upload*, Git*, IACM*, Email, ServiceNow*, Jira*):
```yaml
- id: step1
  name: My Step
  template:
    uses: buildAndPushToDocker
    with:
      connector: myConn
      repo: myrepo
      tags:
        - latest
```
Converter returns `*v1.StepTemplate`.

### Pattern C: Specialized v1 fields
`approval:`, `barrier:`, `queue:`, `wait:`, `run-test:` — used for specific step types.

---

## 4. Step-by-Step: Adding a New Step Converter

Given v0 and v1 YAML, follow these steps:

### Step 1: Define v0 Spec Struct (if not exists)

File: `convert/harness/yaml/step.go`

```go
type StepMyNewStep struct {
    CommonStepSpec                               // always embed this
    ConnectorRef string                          `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
    Region       string                          `json:"region,omitempty"       yaml:"region,omitempty"`
    SomeBool     *flexible.Field[bool]           `json:"someBool,omitempty"     yaml:"someBool,omitempty"`
}
```

**Rules:**
- JSON tags must match v0 YAML field names exactly (camelCase)
- Always embed `CommonStepSpec` for delegate selector support
- Use `*flexible.Field[T]` for non-string fields (see [Section 2.4](#24-flexiblefieldt))

### Step 2: Add Step Type Constant

File: `convert/harness/yaml/const.go`
```go
StepTypeMyNewStep = "MyNewStep"
```

### Step 3: Register in UnmarshalJSON

File: `convert/harness/yaml/step.go`, in `Step.UnmarshalJSON`:
```go
case StepTypeMyNewStep:
    s.Spec = new(StepMyNewStep)
```

### Step 4: Write Converter Function

File: `convert/v0tov1/convert_helpers/convert_step_my_new_step.go`

**For template-based output (Pattern B):**
```go
package converthelpers

import (
    "fmt"
    v0 "github.com/drone/go-convert/convert/harness/yaml"
    v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

func ConvertStepMyNewStep(src *v0.Step) *v1.StepTemplate {
    if src == nil { return nil }
    spec, ok := src.Spec.(*v0.StepMyNewStep)
    if !ok { return nil }

    with := make(map[string]interface{})

    if spec.ConnectorRef != "" {
        with["connector"] = spec.ConnectorRef
    }
    if spec.Region != "" {
        with["region"] = spec.Region
    }

    // Same type → pass through (Case 1, see Section 2.4)
    if spec.SomeBool != nil {
        with["some_bool"] = spec.SomeBool
    }

    // Different type → unwrap and convert (Case 2, see Section 2.4)
    if spec.RunAsUser != nil {
        if user, ok := spec.RunAsUser.AsString(); ok {
            with["run_as_user"] = user
        } else if user, ok := spec.RunAsUser.AsStruct(); ok {
            with["run_as_user"] = fmt.Sprintf("%d", user)
        }
    }

    return &v1.StepTemplate{
        Uses: "myNewStep",
        With: with,
    }
}
```

**For run-based output (Pattern A):**
```go
func ConvertStepMyNewStep(src *v0.Step) *v1.StepRun {
    if src == nil || src.Spec == nil { return nil }
    spec, ok := src.Spec.(*v0.StepMyNewStep)
    if !ok { return nil }

    return &v1.StepRun{
        Script:    v1.Stringorslice{spec.Command},
        Shell:     strings.ToLower(spec.Shell),
        Container: &v1.Container{Image: spec.Image, Connector: spec.ConnRef},
    }
}
```

### Step 5: Wire into Dispatcher

File: `convert/v0tov1/pipeline_converter/convert_steps.go`

Add to the switch in `ConvertSingleStep`:
```go
case v0.StepTypeMyNewStep:
    step.Template = convert_helpers.ConvertStepMyNewStep(src)
```

### Step 6: Add v1 Type Constant (if template-based)

File: `convert/v0tov1/yaml/const.go` — add constant if needed for the `uses` value.

### Step 7: Write Tests

File: `convert/v0tov1/convert_helpers/convert_step_my_new_step_test.go`

Test the converter function directly with constructed v0 structs and verify the v1 output fields.

---

## 5. Field Mapping Guidelines

When the user provides v0 and v1 YAML without explicit field mapping, use these rules:

### 5.1 Automatic Mappings (handled by `convertCommonStepSettings`)

These v0 step-level fields are auto-converted — **do NOT re-map them in the converter function**:

| v0 field | v1 field | Notes |
|---|---|---|
| `identifier` | `id` | Direct copy |
| `name` | `name` | Direct copy |
| `timeout` | `timeout` | Direct copy |
| `failureStrategies` | `on-failure` | Converted by `ConvertFailureStrategies` |
| `envVariables` | `env` | Direct copy (FlexibleField) |
| `when` | `if` | Converted by `ConvertStepWhen` |
| `strategy` | `strategy` | Converted by `ConvertStrategy` |
| `delegateSelectors` + `includeInfraSelectors` | `delegate` | Converted by `ConvertDelegate` |

### 5.2 Spec-Level Field Mapping Rules

For fields inside `spec:` (the step-specific part):

1. **Same name, same type** — map directly:
   - v0: `spec.region: us-east-1` → v1: `with.region: us-east-1`

2. **Renamed field** — match by position/meaning in the provided YAML:
   - v0: `spec.connectorRef` → v1: `with.connector` (common rename)
   - v0: `spec.repo` → v1: `with.repo` (same name)
   - v0: `spec.imageName` → v1: `with.image_name` (camelCase → snake_case)

3. **Common renames across existing converters** (use as defaults):

   | v0 spec field | v1 `with` key | Notes |
   |---|---|---|
   | `connectorRef` | `connector` | Always rename |
   | `baseImageConnectorRefs` | `base_image_connector` | Take first element if array |
   | `envVariables` | `env_vars` | For template steps |
   | `remoteCacheRepo` / `remoteCacheImage` | `cache_repo` / `remote_cache_image` | |
   | `imagePullPolicy` | `pull` (in container) | Convert: Always→always, Never→never, IfNotPresent→if-not-exists |
   | `command` | `script` (in run) | |
   | `runAsUser` | `user` (in container) | |

4. **FlexibleField handling** — see [Section 2.4 → FlexibleField Handling in Converters](#flexiblefield-handling-in-converters).

5. **Dropped fields** — if a v0 field has no v1 counterpart, omit it (don't convert).

### 5.3 Ambiguous Mappings

If the mapping cannot be confidently inferred:
- **Ask the user** for clarification on that specific field
- Note which v0 fields were dropped/unmapped in the response
- Flag any v1 fields that have no v0 source

---

## 6. Checklist for Converter Changes

- [ ] v0 spec struct exists in `convert/harness/yaml/step.go` (with `CommonStepSpec` embedded)
- [ ] Step type constant in `convert/harness/yaml/const.go`
- [ ] UnmarshalJSON case in `convert/harness/yaml/step.go`
- [ ] Converter function in `convert/v0tov1/convert_helpers/convert_step_<name>.go`
- [ ] Switch case in `convert/v0tov1/pipeline_converter/convert_steps.go`
- [ ] v1 type constant in `convert/v0tov1/yaml/const.go` (if needed)
- [ ] Test file in `convert/v0tov1/convert_helpers/convert_step_<name>_test.go`
- [ ] Common fields (timeout, env, when, strategy, delegate, failureStrategies) are NOT re-handled in the converter — they're automatic

---

## 7. Modifying an Existing Converter

When the user provides updated v0/v1 YAML for an existing step type:

1. **Diff the YAML** against the current converter output
2. **Identify new/changed/removed fields** in v0 spec struct
3. **Update the v0 struct** if new fields were added
4. **Update the converter function** to map new fields
5. **Update tests** to cover the new mappings
6. Do NOT touch the dispatcher or UnmarshalJSON unless the step type changed

# Harness v0 to v1 Step Converter Template Sync

## Overview
This skill enables syncing step converters in the go-convert repository with the latest template definitions from the Harness templates repository. Step converters transform v0 Harness pipeline steps to v1 format, where many steps are converted to template-based steps using `StepTemplate` with `uses` and `with` fields.

## Repository Paths

### Path Configuration (One-Time Setup)

Repository paths can be configured once and reused across all sync commands. The AI assistant resolves paths in this priority order:

#### Option 1: Configuration File (Recommended)

Create `.claude/skills/sync-with-templates/config.json` in the go-convert repo:

```json
{
  "go_convert_repo": "/Users/me/projects/go-convert",
  "template_library_repo": "/Users/me/projects/template-library"
}
```

The AI assistant will automatically read this file before any sync operation.

#### Option 2: Environment Variables

Set these environment variables in your shell profile (`.bashrc`, `.zshrc`, etc.):

```bash
export GO_CONVERT_REPO="/Users/me/projects/go-convert"
export TEMPLATE_LIBRARY_REPO="/Users/me/projects/template-library"
```

The AI assistant will check for these if no config file exists.

#### Option 3: Inline in Message

Provide paths directly in your sync request:

```
Sync GitClone step
- go-convert: /path/to/go-convert
- template-library: /path/to/template-library
```

Or inline format:
```
Sync GitClone step with go-convert=/Users/me/projects/go-convert and template-library=/Users/me/projects/template-library
```

### Path Resolution Order

The AI assistant resolves paths in this order:

1. **Inline paths in message** — highest priority, overrides all
2. **Config file** — `.claude/skills/sync-with-templates/config.json`
3. **Environment variables** — `GO_CONVERT_REPO` and `TEMPLATE_LIBRARY_REPO`
4. **Prompt user** — if none of the above are found

### AI Assistant Path Resolution Logic

```
1. Check if paths provided in user message → use those
2. Check if config.json exists at $GO_CONVERT_REPO/.claude/skills/sync-with-templates/config.json → read and use
3. Check environment variables via `echo $GO_CONVERT_REPO` and `echo $TEMPLATE_LIBRARY_REPO` → use if set
4. Prompt user for paths
```

**AI Assistant Prompt Template (only if no paths found):**
```
To sync step converters, I need the paths to both repositories.

You can either:
1. Create a config file at .claude/skills/sync-with-templates/config.json:
   {"go_convert_repo": "/path/to/go-convert", "template_library_repo": "/path/to/template-library"}

2. Set environment variables:
   export GO_CONVERT_REPO="/path/to/go-convert"
   export TEMPLATE_LIBRARY_REPO="/path/to/template-library"

3. Provide paths inline with your request

Please provide the paths or set up one of the above options.
```

### Path Variables

Once resolved, paths are referenced as:
- `$GO_CONVERT_REPO` - Path to go-convert repository
- `$TEMPLATE_LIBRARY_REPO` - Path to template-library repository

**go-convert Repository (Converter):**
- Path: `$GO_CONVERT_REPO`
- Contains: `convert/v0tov1/convert_helpers/`, `convert/harness/yaml/`

**template-library Repository (Templates):**
- Path: `$TEMPLATE_LIBRARY_REPO`
- Contains: `.harness/<templateName>/<version>/template.yaml`

**Sync State Tracking:**
- File: `convert/v0tov1/.template-sync-state.yaml` (in go-convert repo)
- Tracks the last template-library commit SHA that converters are synced to
- Format:
  ```yaml
  last_synced_commit: "<commit_sha>"
  last_synced_at: "<ISO8601_timestamp>"
  last_synced_pr: "<PR title or number>"
  ```

## Repository Structure

### go-convert Repository 

**Key Directories:**
- `convert/harness/yaml/` - v0 step struct definitions
  - `step.go` - Main step structs (StepGitClone, StepRun, StepAction, etc.)
  - `step_k8s.go` - K8s-specific step structs
  - `step_helm.go` - Helm specific steps, 
  - In general step related structs will be present in files with pattern `step_*.go`
  - `const.go` - Step type constants (StepTypeGitClone, StepTypeK8sApply, etc.)

- `convert/v0tov1/convert_helpers/` - Step conversion functions
  - `convert_step_git_clone.go` - GitClone step converter
  - `convert_step_k8s.go` - K8s step converters (K8sApply, K8sRollingDeploy, etc.)
  - `convert_step_build_push_docker.go` - Docker build/push converter
  - `convert_step_*.go` - Other step converters

- `convert/v0tov1/yaml/` - v1 output structures
  - `step_template.go` - StepTemplate struct definition
  - `const.go` - v1 step type constants

- `convert/v0tov1/pipeline_converter/convert_steps.go` - Main step conversion dispatcher, look in `ConvertSteps()` for the converter functions used for different step types.

### Templates Repository (template-library)

```
template-library/
├── .harness/                           # All templates live here
│   ├── CODEOWNERS                      # Code ownership definitions
│   ├── [templateName]/                 # camelCase directory name
│   │   ├── config.yaml                 # Version tracking (stable, prod1, prod0)
│   │   └── [version]/                  # e.g., 1.0.0/, 1.1.0/, 1.2.0/
│   │       └── template.yaml           # Template definition
│   │
│   ├── gitCloneStep/                   # Example: Git Clone step template
│   │   ├── config.yaml
│   │   ├── 1.1.1/
│   │   │   └── template.yaml
│   │   └── 1.2.0/
│   │       └── template.yaml
│   │
│   ├── buildAndPushToDocker/           # Docker build/push template
│   ├── buildAndPushToECR/              # ECR build/push template
│   │
│   ├── k8sApplyStep/                   # K8s Apply step
│   ├── k8sDeleteStep/                  # K8s Delete step
│   ├── k8sScaleStep/                   # K8s Scale step
│   │
│   ├── helmBasicDeployStep/            # Helm Basic Deploy step
│   ├── helmRollbackStep/               # Helm Rollback step
│   ├── helmDeleteStep/                 # Helm Delete step
│   ├── helmBlueGreenDeployStep/        # Helm Blue-Green Deploy step
│   ├── helmCanaryDeployStep/           # Helm Canary Deploy step
│   │
│   │
│   └── [security-scanning-steps]/      # Security scanning templates
│   etc.
│
├── scripts/                            # Utility scripts
│   ├── fix_optional_layout_indent.py
│   ├── fix_tooltip_audit.py
│   └── unhide_infra_defaults.py
│
├── manage.py                           # Template management utility
├── AGENTS.md                           # AI agent instructions
├── README.md                           # Repository documentation
└── sanity_pipeline.yaml                # Sanity test pipeline
```

### config.yaml Structure

Each template directory contains a `config.yaml` that tracks versions:

```yaml
stable: 1.2.0                    # Current stable version
versions:
  prod1: 1.2.0                   # Production environment 1 version
  prod0: 1.2.0                   # Production environment 0 version
icon-name: git                   # Icon identifier

metadata:
  category:
    - step                       # Template category (step, strategy, action)
  module:
    - builds                     # Module (builds, cd, ci)
```

## Template Definition Format

Templates define inputs that step converters must map to. Example template structure:

```yaml
template:
  inputs:
    <input_name>:
      type: <string|boolean|connector|number|etc.>
      label: <Display Label>
      required: <true|false>
      default: <default_value>
      options:  # For select inputs
        - Option1
        - Option2
      ui:
        component: <string|select|checkbox>
        tooltip: <Help text>
        visible: <expression>  # Conditional visibility
        placeholder: <placeholder text>

  layout:
    - <input_name>
    - variant: more
      items:
        - <input_name>

  id: <templateId>
  name: <Template Name>
  description: <Description>
  version: <version_number>
  author: harness
  module:
    - ci  # or cd
  alias: <short_alias>
  
  step:
    # Template implementation
```

## Step Converter Pattern

Each step converter follows this pattern:

```go
// ConvertStep<StepType> converts a v0 <StepType> step to v1 template format
func ConvertStep<StepType>(src *v0.Step) *v1.StepTemplate {
    if src == nil || src.Spec == nil {
        return nil
    }

    spec, ok := src.Spec.(*v0.Step<StepType>)
    if !ok {
        return nil
    }

    with := make(map[string]interface{})

    // Map v0 fields to template inputs
    // be sure to omit empty values
    if spec.<V0Field> != "" {
        with["<template_input_name>"] = spec.<V0Field>
    }

    // Handle flexible fields (can be expression or value)
    // If corresponding template input is of same type T, 
    // EXAMPLE: here spec.Optmize is *felxible.Field[bool] and optmize in template is a bool, don't unwrap in this case as it will be marshalled correctly anyways
    if spec.Optimize != nil {
      with["optimize"] = spec.Optimize
    }
    // If corresponding template input is of type string, 
    // EXAMPLE: here spec.RunAsUser is *felxible.Field[int] and run_as_user in template is a string
    if spec.RunAsUser != nil {
      if user, ok := spec.RunAsUser.AsString(); ok {
        with["run_as_user"] = user
      } else if user, ok := spec.RunAsUser.AsStruct(); ok {
        with["run_as_user"] = fmt.Sprintf("%d", user)
      }
    }

    // Handle accordingly if template input is of any other type,
    // on basis of allowed values of the template input,
    // if some template input is required and not exists in the v0 step yaml:
    // - if default value defined in the template set to that
    // otherwise write in comments of the step converter

    return &v1.StepTemplate{
        Uses: "<templateDirectoryName>",  // Must match the template directory name (e.g., "jiraCreate", "k8sApplyStep")
        With: with,
    }
}
```

## Sync Process

### Step 1: Identify Template Changes
1. Read the latest version template YAML file from `.harness/[templateName]/[version]/template.yaml`
2. Parse the `template.inputs` section to get all input definitions
3. Note the `template.id` value (used in `Uses` field)

### Step 2: Compare with Existing Converter
1. Locate the corresponding converter file in `convert/v0tov1/convert_helpers/` 
2. Check the `Uses` field matches `template.id` 
3. Compare `with` map keys against template input names

### Step 3: Update Converter
For each template input:
1. **Find v0 source field**: Check `convert/harness/yaml/` directory for the v0 struct
   - **Note**: v0 and v1 field names may differ — map based on functionality, not just name matching
   - Use the template input's `ui.tooltip` description as a guide to understand the field's purpose
   - Example: v0 `ImageName` might map to v1 template input `image_name` or `repo`
2. **Map field to input**: Add mapping in the converter's `with` map
3. **Handle field types**:
   - See [Step Converter Pattern](#step-converter-pattern)
   - In general match the field type/format/allowed values in the template definition.

### Step 4: Update v1 Constants (if needed)
If template directory name changed, update `convert/v0tov1/yaml/const.go`:
```go
const (
    StepType<Name> = "<templateDirectoryName>"  // e.g., "jiraCreate", "k8sApplyStep"
)
```
**DO NOT HARDCODE THE TEMPLATE VERSIONS ANYWHERE IN THE CONVERTER, JUST USE THE TEMPLATE DIRECTORY NAME**
## Field Mapping Reference

### Common v0 to Template Input Mappings

| v0 Field Pattern | Template Input Pattern | Notes |
|------------------|----------------------|-------|
| `ConnRef` / `ConnectorRef` | `connector` | Connector reference |
| `Repository` / `RepoName` | `repo_name` / `repo` | Repository name |
| `Image` | `image` | Container image |
| `Env` / `EnvVariables` | `env` / `envvars` | Environment variables |
| `Resources.Limits.CPU` | `cpu` | CPU limit |
| `Resources.Limits.Memory` | `memory` | Memory limit |
| `Timeout` | Handled at step level | Not in template `with` |

### FlexibleField Handling
See [Step Converter pattern](#step-converter-pattern)


## Sync Commands

**CRITICAL WORKFLOW RULE**: Always process steps ONE AT A TIME. Never batch-read multiple templates or converters simultaneously. Only load files into context when actively processing that specific step.

If not provided with a specific command type, ask the user for one of the below.

---

### Sync Single Step

When asked to sync a specific step (e.g., "sync GitClone step with latest template"):

#### Phase 1: Discovery (NO file reading yet)
1. **Identify the template directory name** from user request (e.g., "gitCloneStep", "buildAndPushToDocker")
2. **List version directories** in `.harness/<templateName>/` to find the latest semver version
   - Find the latest version by listing version directories and selecting the highest semver (e.g., `1.3.0` > `1.2.0` > `1.1.0`)
   - Do NOT use the `stable` version from `config.yaml` — always use the highest version directory
3. **Record paths** (do NOT read file contents yet):
   - Template path: `.harness/<templateName>/<latestVersion>/template.yaml`
   - Converter search pattern: `convert/v0tov1/convert_helpers/convert_step_*.go`

#### Phase 2: Load & Compare (NOW read files)
4. **Read the template file** from the path identified above
5. **Extract template inputs**: Parse `template.inputs` section, note each input's:
   - Name, type, required, default, tooltip/description
6. **Find the converter file**: Search for `Uses: "<templateName>"` in converter files
7. **Read the converter file** once found
8. **Extract converter mappings**: Parse all `with["<key>"]` assignments
9. **Find the v0 struct**: 
   - Look in `convert/harness/yaml/step*.go` for the struct used in the converter
   - Read only the relevant struct definition

#### Phase 3: Analysis
10. **Compare template inputs vs converter mappings**:
    - Which template inputs have mappings in converter?
    - Which template inputs are missing from converter?
    - Which converter `with` keys don't exist in template? (stale)
    - **Note**: v0 and v1 field names may differ — map based on functionality, not just name matching
    - Use the template input's `ui.tooltip` description to understand the field's purpose when names don't match
11. **Compare template inputs vs v0 struct fields**:
    - Which v0 fields could map to unmapped template inputs?
    - Which v0 fields have no template equivalent?

#### Phase 4: Update (if needed)
12. **Update converter** to add missing input mappings where v0 field exists
13. **Update tests** if new mappings added
14. **Run tests**: `go test ./convert/v0tov1/convert_helpers/... -run Test<StepType>`

#### Phase 5: Report
15. **Generate Mapping Report** (see [Mapping Report Format](#mapping-report-format))

### Sync Multiple Steps

When asked to sync multiple steps:

#### Phase 1: Discovery ONLY (NO file reading)
1. **Create a processing queue** listing each step to sync:
   ```
   Processing Queue:
   1. [ ] buildAndPushToDocker
   2. [ ] buildAndPushToECR  
   3. [ ] gitCloneStep
   ...
   ```
2. **Save the queue** (mentally or in a note) — do NOT read any template/converter files yet

#### Phase 2: Sequential Processing
3. **For each step in the queue, ONE AT A TIME**:
   
   ```
   ═══════════════════════════════════════════════════════
   PROCESSING STEP 1 of N: <stepName>
   ═══════════════════════════════════════════════════════
   ```
   
   a. **NOW load files for this step only**:
      - Read template file
      - Read converter file  
      - Read v0 struct
   
   b. **Analyze and compare** (as in Sync Single Step Phase 3)
   
   c. **Make changes if needed** (as in Sync Single Step Phase 4)
   
   d. **Generate Mapping Report for this step**
   
   e. **Mark step as complete**:
      ```
      Processing Queue:
      1. [✓] buildAndPushToDocker — COMPLETED (2 fields added)
      2. [ ] buildAndPushToECR — NEXT
      3. [ ] gitCloneStep
      ```
   
   f. **Clear context** — forget the files just read before moving to next step
   
   g. **Move to next step** — repeat from (a)

#### Phase 3: Summary
4. **After ALL steps processed**, provide final summary:
   ```
   ═══════════════════════════════════════════════════════
   SYNC COMPLETE: X of Y steps updated
   ═══════════════════════════════════════════════════════
   
   | Step | Status | Changes |
   |------|--------|--------|
   | buildAndPushToDocker | ✓ Updated | +2 fields |
   | buildAndPushToECR | ✓ No changes needed | — |
   | gitCloneStep | ✓ Updated | +1 field |
   ```

### Mapping Report Format

Generate this report **after processing each individual step**:

```markdown
═══════════════════════════════════════════════════════════════
## <StepName> Mapping Report
═══════════════════════════════════════════════════════════════

### Template Info
- **Directory**: <templateDirectoryName>
- **Version**: <X.Y.Z>
- **Path**: `.harness/<dir>/<version>/template.yaml`

### Converter Info  
- **File**: `convert/v0tov1/convert_helpers/<filename>.go`
- **Function**: `ConvertStep<Name>()`
- **Uses**: `"<templateDirectoryName>"`

### v0 Struct Info
- **File**: `convert/harness/yaml/<filename>.go`
- **Struct**: `Step<Name>`

---

### Field Mappings (v0 → template input)

| v0 Field | Template Input | Type | Status |
|----------|----------------|------|--------|
| `ConnectorRef` | `connector` | string | ✓ Mapped |
| `Repo` | `repo` | string | ✓ Mapped |
| `Tags` | `tags` | []string | ✓ Mapped |
| `Caching` | `caching` | bool | ✓ Mapped |

---

### Template Inputs Without v0 Mapping

| Template Input | Type | Required | Default | Action Taken |
|----------------|------|----------|---------|-------------|
| `build_mode` | string | yes | — | Set default: "build_and_push" |
| `platforms` | []string | no | — | No v0 equivalent, skipped |

---

### v0 Fields Without Template Mapping

| v0 Field | Type | Reason |
|----------|------|--------|
| `Resources` | *Resources | Handled at step level, not in template |
| `Reports` | []*Report | No template input exists |

---

### Converter Status Checklist
- [x] `Uses` field matches template directory name
- [x] All required template inputs mapped or have defaults
- [x] No stale mappings (all `with` keys exist in template)
- [ ] Optional inputs mapped where v0 equivalent exists

### Changes Made
- Added mapping: `spec.Dockerfile` → `with["dockerfile"]`
- Added mapping: `spec.Context` → `with["context"]`
- Set default: `with["build_mode"] = "build_and_push"`
```

### Sync All Steps

When asked to sync all steps in the repository:

#### Phase 1: Discovery ONLY
1. **Read `convert/v0tov1/pipeline_converter/convert_steps.go`** to find all template-based converters
2. **Extract step types** that assign to `step.Template = ...`
3. **For each, find the `Uses` value** by searching converter files
4. **Build processing queue**:
   ```
   Discovered N template-based steps to sync:
   1. [ ] buildAndPushToDocker (from StepTypeBuildAndPushDockerRegistry)
   2. [ ] buildAndPushToECR (from StepTypeBuildAndPushECR)
   ...
   ```
5. **STOP and confirm with user** before proceeding:
   > "Found N steps to sync. Proceed with sequential processing?"

#### Phase 2: Sequential Processing
6. **Follow "Sync Multiple Steps" Phase 2** exactly — one step at a time

### Sync from Template-Library PR Merge

This flow is triggered after a PR is merged in the template-library repo. It uses git diff to determine exactly which templates changed and only syncs those converters.

#### Prerequisites
- The file `convert/v0tov1/.template-sync-state.yaml` must exist with a valid `last_synced_commit` SHA
- If the file doesn't exist, create it by:
  1. Ask the user for the template-library commit SHA to use as baseline (or use current HEAD if doing a full sync first)
  2. Create the file with that SHA

#### Phase 1: Identify Changed Templates ONLY (NO file reading)
1. Navigate to the template-library repo at `$TEMPLATE_LIBRARY_REPO`
2. Read `last_synced_commit` from `convert/v0tov1/.template-sync-state.yaml` in go-convert
3. Run: `git log --oneline <last_synced_commit>..HEAD -- .harness/`
   - This shows all commits affecting templates since last sync
   - If `last_synced_commit` is not an ancestor of HEAD (e.g., after a rebase), fall back to: `git diff --name-only <last_synced_commit> HEAD -- .harness/`
4. Run: `git diff --name-only <last_synced_commit> HEAD -- .harness/`
   - **Extract unique template directory names** from changed paths
   - e.g., `.harness/k8sApplyStep/1.3.0/template.yaml` → `k8sApplyStep`
   - Ignore changes to `CODEOWNERS`, scripts, or non-template files
5. **Build processing queue** of changed templates:
   ```
   Templates changed since last sync:
   1. [ ] k8sApplyStep (files: template.yaml modified)
   2. [ ] helmDeployStep (files: new version 1.3.0 added)
   ...
   ```
6. **Map each to converter** (if exists) — just record the mapping, don't read files:
   ```
   Template → Converter mapping:
   1. k8sApplyStep → convert_step_k8s.go (ConvertStepK8sApply)
   2. helmDeployStep → convert_step_helm.go (ConvertStepHelmDeploy)
   3. newTemplate → NO CONVERTER FOUND (skip)
   ```

#### Phase 2: Sequential Processing  
7. **For each template with a converter, ONE AT A TIME**:
   - Follow "Sync Single Step" Phases 2-5
   - Only NOW read the template, converter, and v0 struct files

#### Phase 3: Update Sync State
After all converters are successfully synced:
1. Get the current HEAD SHA of the template-library repo: `git rev-parse HEAD`
2. Get the PR title/number from the user or from: `git log -1 --format="%s"`
3. Update `convert/v0tov1/.template-sync-state.yaml`:
   ```yaml
   last_synced_commit: "<new_HEAD_sha>"
   last_synced_at: "<current_ISO8601_timestamp>"
   last_synced_pr: "<PR title or description>"
   ```
4. Commit the tracking file update along with the converter changes

#### Phase 4: Report
Summarize:
- Which templates changed in the template-library since last sync
- Which converters were updated (and what changed)
- Which changed templates had no corresponding converter
- Which changed templates had only non-functional changes (e.g., tooltip text, layout order) — no converter update needed
- The new sync state

#### Edge Cases
- **First-time sync**: If no `.template-sync-state.yaml` exists, ask the user whether to do a full sync of all templates or to set the current HEAD as baseline
- **Template deleted**: If a template directory was deleted, report it but don't modify the converter
- **New template added**: If a new template directory was added with no existing converter, report it as needing a new converter
- **Only config.yaml changed**: If only `config.yaml` changed (version bump without new version directory), no sync needed unless a new version directory was added
- **Only layout/UI changes**: If template inputs didn't change (only `layout`, `ui.tooltip`, `ui.visible`, etc.), report "no converter update needed" and skip

## Testing

After updating a converter:
1. Run existing tests: `go test ./convert/v0tov1/convert_helpers/... -run Test<StepType>` and full test suite `go test ./...`
2. modify test cases for new input mappings
3. Verify YAML output matches expected template format

## Validation Commands

### Validate Sync (Dry-Run)
Verify that existing converters are in sync with the latest templates **without modifying any files**. Use this as a sanity check.

For each converter that uses a template:
1. Read the converter file and extract the `Uses` field (template directory name) and all `with` map keys
2. Find the corresponding template directory in template-library (the `Uses` value IS the directory name)
3. Read the **latest version** template and extract all `template.inputs` names
4. Report:
   - ✅ `Uses` field matches the template directory name
   - ✅ Inputs with correct mappings in converter
   - ⚠️ Template inputs missing from converter `with` map (may need mapping)
   - ⚠️ Converter `with` keys that don't exist in template inputs (stale mappings)
   - ℹ️ Template inputs that are optional with defaults (may not need mapping)
5. Summary: "X of Y converters fully in sync"

**Usage**: When asked to "validate sync" or "check sync status" or "dry-run sync", run this validation without making changes.

### Validate Single Step
Same as above but for a specific step. Report detailed field-by-field comparison.

## Testing Skill Accuracy

### Regression Test Procedure
To verify the skill correctly syncs a converter:

1. **Pick a known-good converter** (e.g., `buildAndPushToDocker`) that is currently in sync
2. **Record golden state**: Copy the current converter file content as expected output
3. **Introduce a regression**: Manually remove one or more `with` mappings from the converter
4. **Run "Sync Single Step"** for that converter
5. **Diff result against golden state**: 
   - The removed mappings should be restored
   - No unintended changes should be introduced
   - `Uses` field should remain correct

### Template Change Detection Test
To verify the PR-merge sync flow correctly identifies changed templates:

1. **Find two known commits** in template-library where a specific template changed
2. **Set `.template-sync-state.yaml`** to the older commit SHA
3. **Run "Validate Sync" (dry-run)** for the PR-merge flow
4. **Verify detection**:
   - The changed template(s) should be identified
   - The correct converter file(s) should be mapped
   - Templates that didn't change should NOT be flagged

### Field Mapping Accuracy Checklist
For each synced converter, verify:
- [ ] Every **required** template input without a default has a mapping in the converter
- [ ] Every `with[key]` in the converter matches an actual input name in the template
- [ ] The `Uses` field exactly matches the template directory name (e.g., `jiraCreate`, `k8sApplyStep`)
- [ ] No stale mappings exist for inputs that were removed from the template
- [ ] FlexibleField types are handled correctly (expression passthrough vs value extraction)
- [ ] Type conversions are correct (e.g., int to string when template expects string)
- [ ] Field name differences are handled (v0 and v1 names may differ — mapped by functionality)
- [ ] Full mapping report generated showing all template inputs and v0 fields

---

## Quick Reference

| Task | Command/Location |
|------|------------------|
| Find converter for template | `grep -r "Uses: \"<templateName>\"" convert/v0tov1/convert_helpers/` |
| Find v0 struct | `convert/harness/yaml/step*.go` |
| Find step type constant | `convert/harness/yaml/const.go` |
| Find step dispatcher | `convert/v0tov1/pipeline_converter/convert_steps.go` |
| Run converter tests | `go test ./convert/v0tov1/convert_helpers/... -run Test<StepType>` |
| Template location | `.harness/<templateName>/<version>/template.yaml` |
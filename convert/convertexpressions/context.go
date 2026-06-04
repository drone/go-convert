package convertexpressions

// StepInfoFQN is the FQN-keyed step info record. The full v1 FQN is the map key
// in ConversionContext.StepInfoByFQN, so V1Path is not stored here.
type StepInfoFQN struct {
	Type string `json:"type,omitempty"` // v0 step type (e.g. "Run", "StepGroup")
	// Types holds candidate v0 types when a single v1 field maps to several v0
	// types (e.g. v1 step.run -> Run/ShellScript/Plugin/...). When non-empty it
	// takes precedence over Type during step-type resolution.
	Types   []string `json:"types,omitempty"`
	StageID string   `json:"stage_id,omitempty"`
	Chain   []string `json:"chain,omitempty"`   // ancestor step-group IDs, stage→leaf, excluding the leaf
	StepID  string   `json:"step_id,omitempty"` // leaf step/group ID (final FQN segment)
}

// ConversionContext holds metadata for context-aware conversion
type ConversionContext struct {
	StepType string // step type (e.g. "Run", "Http"); resolved lazily by the trie
	CurrentStepID   string // ID of the current step (from postprocess)

	// CurrentStepType is the type of the step we're inside; the fallback for
	// "step.spec.*" expressions (no explicit step ID).
	CurrentStepType string

	// CurrentStepTypes holds candidate types for the step we're inside, used
	// like Types when the current step's v1 field maps to several v0 types.
	CurrentStepTypes []string

	// Warnings collects per-conversion diagnostics (e.g. unmapped template
	// uses, ambiguous step types resolved via best-match fallback). The caller
	// resets this between expressions to attribute warnings per expression.
	Warnings []string

	// UseFQN enables FQN mode: at step_node the v1Path becomes the step's v1
	// FQN base path.
	UseFQN bool

	// CurrentStageID is the call-site's stage ID (set during the postprocess walk).
	CurrentStageID string

	// --- FQN-keyed step info lookup (sole step-resolution model) ---

	// StepInfoByFQN maps each step/step-group's full v1 FQN to its metadata
	// (e.g. "pipeline.stages.prod.steps.PostDeploy.steps.smoke_test"). Callers
	// filter step groups by Type.
	StepInfoByFQN map[string]*StepInfoFQN

	// CurrentStepGroupChain is the call-site's enclosing step-group chain
	// (outermost first, leaf parent last; empty at stage level). Group-relative
	// aliases prepend it to the captured chain to form effectiveChain.
	CurrentStepGroupChain []string

	// CurrentFQN is the call-site step's full v1 FQN; short-circuits references
	// to "the current step" (e.g. "step.spec.*"). Empty when not inside a step.
	CurrentFQN string

	// Lazy indexes built by EnsureIndexes from StepInfoByFQN:
	indexesBuilt     bool
	stageIDs         []string                       // deduped stage IDs, for fuzzy matching
	stepsByStageStep map[string]map[string][]string // [stage][stepID] -> FQNs; candidate gather
	flatStepIDCount  map[string]int                 // leaf stepID -> count across stages (flat fallback)
	flatStepIDFQN    map[string]string              // leaf stepID -> canonical FQN (when count==1)
}

// EnsureIndexes builds the lazy lookup indexes from StepInfoByFQN once per
// context. Callers that mutate StepInfoByFQN afterwards must reset indexesBuilt.
func (c *ConversionContext) EnsureIndexes() {
	if c == nil || c.indexesBuilt {
		return
	}
	c.indexesBuilt = true
	if c.StepInfoByFQN == nil {
		return
	}
	stageSet := map[string]struct{}{}
	c.stepsByStageStep = make(map[string]map[string][]string, len(c.StepInfoByFQN))
	c.flatStepIDCount = make(map[string]int, len(c.StepInfoByFQN))
	c.flatStepIDFQN = make(map[string]string, len(c.StepInfoByFQN))
	for fqn, info := range c.StepInfoByFQN {
		if info == nil {
			continue
		}
		stageSet[info.StageID] = struct{}{}
		byStep := c.stepsByStageStep[info.StageID]
		if byStep == nil {
			byStep = make(map[string][]string)
			c.stepsByStageStep[info.StageID] = byStep
		}
		byStep[info.StepID] = append(byStep[info.StepID], fqn)
		c.flatStepIDCount[info.StepID]++
		c.flatStepIDFQN[info.StepID] = fqn
	}
	c.stageIDs = make([]string, 0, len(stageSet))
	for s := range stageSet {
		c.stageIDs = append(c.stageIDs, s)
	}
}

// addWarning appends a diagnostic message to the context (no-op on nil).
func (c *ConversionContext) addWarning(msg string) {
	if c == nil || msg == "" {
		return
	}
	c.Warnings = append(c.Warnings, msg)
}

package converter

import (
	"strings"

	convertexpressions "github.com/drone/go-convert/convert/convertexpressions"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// buildStepInfoByFQNFromV1 walks a parsed v1 pipeline and produces the
// FQN-keyed step-info map consumed by the new ResolveStepFQN lookup. 
func buildStepInfoByFQNFromV1(p *v1.Pipeline) map[string]*convertexpressions.StepInfoFQN {
	if p == nil {
		return nil
	}
	out := make(map[string]*convertexpressions.StepInfoFQN)
	for _, stage := range p.Stages {
		walkV1Stage(stage, out)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// walkV1Stage registers all steps under a single stage and recurses into
// container stages (group / parallel) which themselves hold nested stages.
func walkV1Stage(stage *v1.Stage, out map[string]*convertexpressions.StepInfoFQN) {
	if stage == nil {
		return
	}
	// Container stages hold nested stages rather than steps.
	if stage.Group != nil {
		for _, s := range stage.Group.Stages {
			walkV1Stage(s, out)
		}
	}
	if stage.Parallel != nil {
		for _, s := range stage.Parallel.Stages {
			walkV1Stage(s, out)
		}
	}
	if stage.Id == "" {
		return
	}
	base := "pipeline.stages." + stage.Id + ".steps."
	for _, step := range stage.Steps {
		walkV1Step(step, stage.Id, nil, base, out)
	}
}

// walkV1Step registers a single step (or step group) and recurses into nested
// group/parallel steps. prefix is the FQN up to and including the trailing
// ".steps." for this level; chain is the ancestor step-group ID list.
func walkV1Step(step *v1.Step, stageID string, chain []string, prefix string, out map[string]*convertexpressions.StepInfoFQN) {
	if step == nil || step.Id == "" {
		return
	}
	fqn := prefix + step.Id
	info := &convertexpressions.StepInfoFQN{
		StageID: stageID,
		Chain:   append([]string(nil), chain...),
		StepID:  step.Id,
	}
	// A single v0 type (Type) when unambiguous; a candidate list (Types) when
	// one v1 field maps to several v0 types (e.g. step.run).
	if types := v1StepTypes(step); len(types) == 1 {
		info.Type = types[0]
	} else if len(types) > 1 {
		info.Types = types
	}
	out[fqn] = info

	// Recurse into step groups (group / parallel both hold child steps).
	childChain := append(append([]string(nil), chain...), step.Id)
	childPrefix := fqn + ".steps."
	if step.Group != nil {
		for _, child := range step.Group.Steps {
			walkV1Step(child, stageID, childChain, childPrefix, out)
		}
	}
	if step.Parallel != nil {
		for _, child := range step.Parallel.Steps {
			walkV1Step(child, stageID, childChain, childPrefix, out)
		}
	}
}

// v1StepTypes derives the candidate v0 step type(s) from which v1 field is set.
func v1StepTypes(step *v1.Step) []string {
	switch {
	case step.Group != nil:
		return []string{"StepGroup"}
	case step.Run != nil:
		// HTTP converts to step.run but has a distinct, recoverable signature
		// (hardcoded plugin image / PLUGIN_URL); all other run-family types
		// share a clash-free rule set.
		if isHTTPRunStep(step.Run) {
			return []string{convertexpressions.StepTypeHTTP}
		}
		return convertexpressions.RunCandidateTypes
	case step.RunTest != nil:
		return convertexpressions.RunTestCandidateTypes
	case step.Background != nil:
		return []string{convertexpressions.StepTypeBackground}
	case step.Approval != nil:
		if t := convertexpressions.ApprovalUsesToV0Type[step.Approval.Uses]; t != "" {
			return []string{t}
		}
		return nil
	case step.Barrier != nil:
		return []string{convertexpressions.StepTypeBarrier}
	case step.Queue != nil:
		return []string{convertexpressions.StepTypeQueue}
	case step.Wait != nil:
		return []string{convertexpressions.StepTypeWait}
	case step.Template != nil:
		if t := convertexpressions.TemplateUsesToV0Type[step.Template.Uses]; t != "" {
			return []string{t}
		}
		return nil
	default:
		return nil
	}
}

// isHTTPRunStep reports whether a step.run was produced from a v0 HTTP step,
// detected by the hardcoded harness-http plugin image or its PLUGIN_URL env.
func isHTTPRunStep(run *v1.StepRun) bool {
	if run == nil {
		return false
	}
	if run.Container != nil && strings.Contains(run.Container.Image, "harnessdev/harness-http") {
		return true
	}
	return false
}

package converter

import (
	"strings"
	"testing"

	pipelineconverter "github.com/drone/go-convert/convert/v0tov1/pipeline_converter"
)

func TestApplyRefMappings_NoMaps_Passthrough(t *testing.T) {
	in := []byte("pipeline:\n  id: p1\n")
	out, err := ApplyRefMappings(in, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(in) {
		t.Fatalf("expected passthrough, got: %q", string(out))
	}
}

func TestApplyRefMappings_TemplateUses_RefOnly(t *testing.T) {
	in := []byte("pipeline:\n  stages:\n    - template:\n        uses: oldRef@v1\n")
	out, err := ApplyRefMappings(in, map[string]string{"oldRef": "newRef"}, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "uses: newRef@v1") {
		t.Fatalf("expected 'uses: newRef@v1' in:\n%s", string(out))
	}
}

func TestApplyRefMappings_ChainUses_PipelineSegment(t *testing.T) {
	// The map is keyed by the bare pipeline segment; the value is a
	// fully-qualified ref that replaces the whole chain uses verbatim
	// (org/project prefix included).
	in := []byte("pipeline:\n  stages:\n    - chain:\n        uses: org1/proj1/oldPipe\n")
	out, err := ApplyRefMappings(in, nil, map[string]string{"oldPipe": "org1/proj1/newPipe"}, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "uses: org1/proj1/newPipe") {
		t.Fatalf("expected chain uses rewritten; got:\n%s", string(out))
	}
}

func TestApplyRefMappings_ChainUses_FullValueKey_NotMatched(t *testing.T) {
	// Only the pipeline segment (3rd path element) is looked up; a
	// fully-qualified key must NOT match, so the value is left unchanged.
	in := []byte("pipeline:\n  stages:\n    - chain:\n        uses: org1/proj1/oldPipe\n")
	out, err := ApplyRefMappings(in, nil, map[string]string{"org1/proj1/oldPipe": "org2/proj2/newPipe"}, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "uses: org1/proj1/oldPipe") {
		t.Fatalf("expected chain uses unchanged for full-value key; got:\n%s", string(out))
	}
}

func TestApplyRefMappings_InputSetOverlayId(t *testing.T) {
	in := []byte("inputs:\n  overlay:\n    id: oldPipe\n    stages: []\n")
	out, err := ApplyRefMappings(in, nil, map[string]string{"oldPipe": "newPipe"}, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "id: newPipe") {
		t.Fatalf("expected overlay.id rewritten; got:\n%s", string(out))
	}
}

func TestApplyRefMappings_PipelineId(t *testing.T) {
	in := []byte("pipeline:\n  id: oldId\n  name: X\n")
	out, err := ApplyRefMappings(in, nil, map[string]string{"oldId": "newId"}, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "id: newId") {
		t.Fatalf("expected pipeline.id rewritten; got:\n%s", string(out))
	}
}

func TestApplyRefMappings_TriggerPipelineIdentifier(t *testing.T) {
	in := []byte("trigger:\n  name: t1\n  pipelineIdentifier: oldPipe\n")
	out, err := ApplyRefMappings(in, nil, map[string]string{"oldPipe": "newPipe"}, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "pipelineIdentifier: newPipe") {
		t.Fatalf("expected pipelineIdentifier rewritten; got:\n%s", string(out))
	}
}

func TestApplyRefMappings_TriggerInputYamlRecurses(t *testing.T) {
	inner := "pipeline:\n  id: oldPipe1\n  stages:\n    - chain:\n        uses: org/proj/oldPipe2\n"
	in := []byte("trigger:\n  pipelineIdentifier: oldPipe1\n  inputYaml: |\n    " + strings.ReplaceAll(inner, "\n", "\n    ") + "\n")
	out, err := ApplyRefMappings(in, nil, map[string]string{"oldPipe1": "newPipe1", "oldPipe2": "org/proj/newPipe2"}, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "pipelineIdentifier: newPipe1") {
		t.Fatalf("expected pipelineIdentifier rewritten; got:\n%s", s)
	}
	if !strings.Contains(s, "id: newPipe1") {
		t.Fatalf("expected embedded pipeline.id rewritten; got:\n%s", s)
	}
	if !strings.Contains(s, "uses: org/proj/newPipe2") {
		t.Fatalf("expected embedded chain uses rewritten; got:\n%s", s)
	}
}

func TestApplyRefMappings_MissingKey_NoChange(t *testing.T) {
	in := []byte("pipeline:\n  id: keep\n")
	out, err := ApplyRefMappings(in, map[string]string{"other": "x"}, map[string]string{"another": "y"}, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "id: keep") {
		t.Fatalf("expected id unchanged; got:\n%s", string(out))
	}
}

func TestApplyRefMappings_BothMapsAppliedIndependently(t *testing.T) {
	in := []byte("pipeline:\n  id: oldPipe\n  stages:\n    - template:\n        uses: oldT@v1\n")
	out, err := ApplyRefMappings(in,
		map[string]string{"oldT": "newT"},
		map[string]string{"oldPipe": "newPipe"},
		nil, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "id: newPipe") {
		t.Fatalf("expected pipeline.id rewritten; got:\n%s", s)
	}
	if !strings.Contains(s, "uses: newT@v1") {
		t.Fatalf("expected template.uses rewritten; got:\n%s", s)
	}
}

func TestApplyRefMappings_NonInputOverlayId_NotRewritten(t *testing.T) {
	in := []byte("other:\n  overlay:\n    id: oldPipe\n")
	out, err := ApplyRefMappings(in, nil, map[string]string{"oldPipe": "newPipe"}, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(string(out), "id: newPipe") {
		t.Fatalf("expected overlay.id NOT rewritten when grandparent != 'inputs'; got:\n%s", string(out))
	}
}

func TestApplyRefMappings_InputYaml_MalformedLogsWarning(t *testing.T) {
	pipelineconverter.ResetMessageLogger()
	pipelineconverter.GetMessageLogger().Enable("")
	pipelineconverter.GetMessageLogger().SetCurrentFile("test")

	in := []byte("trigger:\n  inputYaml: \"not: valid: yaml: [[\"\n")
	out, err := ApplyRefMappings(in, nil, map[string]string{"oldPipe": "newPipe"}, nil, nil)
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	if !strings.Contains(string(out), "not: valid: yaml") {
		t.Fatalf("expected malformed inputYaml unchanged; got:\n%s", string(out))
	}

	fl := pipelineconverter.GetMessageLogger().GetFileLog("test")
	if fl == nil || len(fl.Messages) == 0 {
		t.Fatal("expected a warning message to be logged for malformed inputYaml")
	}
	found := false
	for _, m := range fl.Messages {
		if m.Code == "INPUT_YAML_REF_MAPPING_FAILED" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected INPUT_YAML_REF_MAPPING_FAILED warning; not found in logged messages")
	}
}

// --- destination-branch remap (template_ref_branch_mapping / pipeline_ref_branch_mapping) ---

func TestApplyRefMappings_TemplateUsesAndBranch_BothRemapped(t *testing.T) {
	in := []byte("pipeline:\n  stages:\n    - template:\n        uses: oldRef@v1\n        ref: sourceBranch\n")
	out, err := ApplyRefMappings(in,
		map[string]string{"oldRef": "newRef"}, nil,
		map[string]string{"oldRef": "destBranch"}, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "uses: newRef@v1") {
		t.Fatalf("expected 'uses: newRef@v1' in:\n%s", s)
	}
	if !strings.Contains(s, "ref: destBranch") {
		t.Fatalf("expected 'ref: destBranch' in:\n%s", s)
	}
}

func TestApplyRefMappings_ChainUsesAndBranch_BothRemapped(t *testing.T) {
	in := []byte("pipeline:\n  stages:\n    - chain:\n        uses: org1/proj1/oldPipe\n        ref: sourceBranch\n")
	out, err := ApplyRefMappings(in,
		nil, map[string]string{"oldPipe": "org1/proj1/newPipe"},
		nil, map[string]string{"oldPipe": "destBranch"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "uses: org1/proj1/newPipe") {
		t.Fatalf("expected chain uses rewritten; got:\n%s", s)
	}
	if !strings.Contains(s, "ref: destBranch") {
		t.Fatalf("expected 'ref: destBranch' in:\n%s", s)
	}
}

func TestApplyRefMappings_TemplateUsesRemapped_NoBranchEntry_RefUntouched(t *testing.T) {
	in := []byte("pipeline:\n  stages:\n    - template:\n        uses: oldRef@v1\n        ref: sourceBranch\n")
	out, err := ApplyRefMappings(in,
		map[string]string{"oldRef": "newRef"}, nil,
		nil, nil, // no branch-mapping entry for "oldRef"
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "uses: newRef@v1") {
		t.Fatalf("expected 'uses: newRef@v1' in:\n%s", s)
	}
	if !strings.Contains(s, "ref: sourceBranch") {
		t.Fatalf("expected 'ref: sourceBranch' untouched; got:\n%s", s)
	}
}

func TestApplyRefMappings_ChainUsesRemapped_NoBranchEntry_RefUntouched(t *testing.T) {
	in := []byte("pipeline:\n  stages:\n    - chain:\n        uses: org1/proj1/oldPipe\n        ref: sourceBranch\n")
	out, err := ApplyRefMappings(in,
		nil, map[string]string{"oldPipe": "org1/proj1/newPipe"},
		nil, nil, // no branch-mapping entry for "oldPipe"
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "uses: org1/proj1/newPipe") {
		t.Fatalf("expected chain uses rewritten; got:\n%s", s)
	}
	if !strings.Contains(s, "ref: sourceBranch") {
		t.Fatalf("expected 'ref: sourceBranch' untouched; got:\n%s", s)
	}
}

func TestApplyRefMappings_NeitherBranchMapPresent_ByteIdenticalToLegacyBehavior(t *testing.T) {
	// Same input/maps as TestApplyRefMappings_BothMapsAppliedIndependently,
	// but explicitly exercising the 4-map signature with nil branch maps to
	// prove today's exact (pre-feature) output is unchanged.
	in := []byte("pipeline:\n  id: oldPipe\n  stages:\n    - template:\n        uses: oldT@v1\n")
	outNoBranchMaps, err := ApplyRefMappings(in,
		map[string]string{"oldT": "newT"},
		map[string]string{"oldPipe": "newPipe"},
		nil, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	outEmptyBranchMaps, err := ApplyRefMappings(in,
		map[string]string{"oldT": "newT"},
		map[string]string{"oldPipe": "newPipe"},
		map[string]string{},
		map[string]string{},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(outNoBranchMaps) != string(outEmptyBranchMaps) {
		t.Fatalf("expected nil and empty branch maps to produce identical output:\nnil:\n%s\nempty:\n%s",
			string(outNoBranchMaps), string(outEmptyBranchMaps))
	}
}

func TestApplyRefMappings_OrphanedBranchEntry_NoMatchingIdMapping_NoOp(t *testing.T) {
	// "oldRef" has a branch-mapping entry but is NOT present in the
	// id-mapping map at all, so no "uses" rewrite happens for it. Per spec,
	// orphaned branch-mapping entries are never looked up (not an error) —
	// the sibling ref field must stay untouched.
	in := []byte("pipeline:\n  stages:\n    - template:\n        uses: oldRef@v1\n        ref: sourceBranch\n")
	out, err := ApplyRefMappings(in,
		nil, nil, // no id-mapping entry for "oldRef"
		map[string]string{"oldRef": "destBranch"}, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "uses: oldRef@v1") {
		t.Fatalf("expected 'uses:' unchanged (no id-mapping entry); got:\n%s", s)
	}
	if !strings.Contains(s, "ref: sourceBranch") {
		t.Fatalf("expected 'ref: sourceBranch' untouched for orphaned branch entry; got:\n%s", s)
	}
}

func TestApplyRefMappings_OrphanedChainBranchEntry_NoMatchingIdMapping_NoOp(t *testing.T) {
	in := []byte("pipeline:\n  stages:\n    - chain:\n        uses: org1/proj1/oldPipe\n        ref: sourceBranch\n")
	out, err := ApplyRefMappings(in,
		nil, nil, // no id-mapping entry for "oldPipe"
		nil, map[string]string{"oldPipe": "destBranch"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "uses: org1/proj1/oldPipe") {
		t.Fatalf("expected chain 'uses:' unchanged (no id-mapping entry); got:\n%s", s)
	}
	if !strings.Contains(s, "ref: sourceBranch") {
		t.Fatalf("expected 'ref: sourceBranch' untouched for orphaned branch entry; got:\n%s", s)
	}
}

func TestApplyRefMappings_TemplateUsesAndBranch_NoPreexistingRefKey_Appended(t *testing.T) {
	// Real-world "no branch pin in V0" case: the source step has no ref
	// at all (ambient/same-branch reference), so go-convert's base converter
	// never emits the key. A destination-branch override must still land by
	// appending a new "ref" sibling, not by rewriting a nonexistent one.
	in := []byte("pipeline:\n  stages:\n    - template:\n        uses: oldRef@v1\n")
	out, err := ApplyRefMappings(in,
		map[string]string{"oldRef": "newRef"}, nil,
		map[string]string{"oldRef": "destBranch"}, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "uses: newRef@v1") {
		t.Fatalf("expected 'uses: newRef@v1' in:\n%s", s)
	}
	if !strings.Contains(s, "ref: destBranch") {
		t.Fatalf("expected appended 'ref: destBranch' in:\n%s", s)
	}
}

func TestApplyRefMappings_ChainUsesAndBranch_NoPreexistingRefKey_Appended(t *testing.T) {
	in := []byte("pipeline:\n  stages:\n    - chain:\n        uses: org1/proj1/oldPipe\n")
	out, err := ApplyRefMappings(in,
		nil, map[string]string{"oldPipe": "org1/proj1/newPipe"},
		nil, map[string]string{"oldPipe": "destBranch"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "uses: org1/proj1/newPipe") {
		t.Fatalf("expected chain uses rewritten; got:\n%s", s)
	}
	if !strings.Contains(s, "ref: destBranch") {
		t.Fatalf("expected appended 'ref: destBranch' in:\n%s", s)
	}
}

func TestApplyRefMappings_TemplateUsesRemapped_NoBranchEntry_NoRefAppended(t *testing.T) {
	// No pre-existing key AND no branch-mapping entry: nothing should be
	// appended (no stray "ref: ''").
	in := []byte("pipeline:\n  stages:\n    - template:\n        uses: oldRef@v1\n")
	out, err := ApplyRefMappings(in,
		map[string]string{"oldRef": "newRef"}, nil,
		nil, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if strings.Contains(s, "\n        ref:") {
		t.Fatalf("expected no ref key appended when there is no branch-mapping entry; got:\n%s", s)
	}
}

func TestApplyRefMappings_TemplateBranch_KeyOrderIndependent(t *testing.T) {
	// ref appears BEFORE uses in the source mapping (valid YAML; key
	// order is caller's choice). The rewrite must still apply because
	// resolvedBranch is computed via an up-front scan over the whole
	// mapping node's Content, not a sequential single-pass lookup.
	in := []byte("pipeline:\n  stages:\n    - template:\n        ref: sourceBranch\n        uses: oldRef@v1\n")
	out, err := ApplyRefMappings(in,
		map[string]string{"oldRef": "newRef"}, nil,
		map[string]string{"oldRef": "destBranch"}, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "uses: newRef@v1") {
		t.Fatalf("expected 'uses: newRef@v1' in:\n%s", s)
	}
	if !strings.Contains(s, "ref: destBranch") {
		t.Fatalf("expected 'ref: destBranch' (key order independent) in:\n%s", s)
	}
}

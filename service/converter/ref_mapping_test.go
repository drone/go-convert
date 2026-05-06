package converter

import (
	"strings"
	"testing"

	pipelineconverter "github.com/drone/go-convert/convert/v0tov1/pipeline_converter"
)

func TestApplyRefMappings_NoMaps_Passthrough(t *testing.T) {
	in := []byte("pipeline:\n  id: p1\n")
	out, err := ApplyRefMappings(in, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(in) {
		t.Fatalf("expected passthrough, got: %q", string(out))
	}
}

func TestApplyRefMappings_TemplateUses_RefOnly(t *testing.T) {
	in := []byte("pipeline:\n  stages:\n    - template:\n        uses: oldRef@v1\n")
	out, err := ApplyRefMappings(in, map[string]string{"oldRef": "newRef"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "uses: newRef@v1") {
		t.Fatalf("expected 'uses: newRef@v1' in:\n%s", string(out))
	}
}


func TestApplyRefMappings_ChainUses_PipelineSegment(t *testing.T) {
	in := []byte("pipeline:\n  stages:\n    - chain:\n        uses: org1/proj1/oldPipe\n")
	out, err := ApplyRefMappings(in, nil, map[string]string{"oldPipe": "newPipe"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "uses: org1/proj1/newPipe") {
		t.Fatalf("expected chain uses rewritten; got:\n%s", string(out))
	}
}

func TestApplyRefMappings_ChainUses_FullValueMatch(t *testing.T) {
	in := []byte("pipeline:\n  stages:\n    - chain:\n        uses: org1/proj1/oldPipe\n")
	out, err := ApplyRefMappings(in, nil, map[string]string{"org1/proj1/oldPipe": "org2/proj2/newPipe"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "uses: org2/proj2/newPipe") {
		t.Fatalf("expected full chain uses rewritten; got:\n%s", string(out))
	}
}

func TestApplyRefMappings_InputSetOverlayId(t *testing.T) {
	in := []byte("inputs:\n  overlay:\n    id: oldPipe\n    stages: []\n")
	out, err := ApplyRefMappings(in, nil, map[string]string{"oldPipe": "newPipe"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "id: newPipe") {
		t.Fatalf("expected overlay.id rewritten; got:\n%s", string(out))
	}
}

func TestApplyRefMappings_PipelineId(t *testing.T) {
	in := []byte("pipeline:\n  id: oldId\n  name: X\n")
	out, err := ApplyRefMappings(in, nil, map[string]string{"oldId": "newId"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "id: newId") {
		t.Fatalf("expected pipeline.id rewritten; got:\n%s", string(out))
	}
}

func TestApplyRefMappings_TriggerPipelineIdentifier(t *testing.T) {
	in := []byte("trigger:\n  name: t1\n  pipelineIdentifier: oldPipe\n")
	out, err := ApplyRefMappings(in, nil, map[string]string{"oldPipe": "newPipe"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "pipelineIdentifier: newPipe") {
		t.Fatalf("expected pipelineIdentifier rewritten; got:\n%s", string(out))
	}
}

func TestApplyRefMappings_TriggerInputYamlRecurses(t *testing.T) {
	inner := "pipeline:\n  id: oldPipe\n  stages:\n    - chain:\n        uses: org/proj/oldPipe\n"
	in := []byte("trigger:\n  pipelineIdentifier: oldPipe\n  inputYaml: |\n    " + strings.ReplaceAll(inner, "\n", "\n    ") + "\n")
	out, err := ApplyRefMappings(in, nil, map[string]string{"oldPipe": "newPipe"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "pipelineIdentifier: newPipe") {
		t.Fatalf("expected pipelineIdentifier rewritten; got:\n%s", s)
	}
	if !strings.Contains(s, "id: newPipe") {
		t.Fatalf("expected embedded pipeline.id rewritten; got:\n%s", s)
	}
	if !strings.Contains(s, "org/proj/newPipe") {
		t.Fatalf("expected embedded chain uses rewritten; got:\n%s", s)
	}
}

func TestApplyRefMappings_MissingKey_NoChange(t *testing.T) {
	in := []byte("pipeline:\n  id: keep\n")
	out, err := ApplyRefMappings(in, map[string]string{"other": "x"}, map[string]string{"another": "y"})
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
	out, err := ApplyRefMappings(in, nil, map[string]string{"oldPipe": "newPipe"})
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
	out, err := ApplyRefMappings(in, nil, map[string]string{"oldPipe": "newPipe"})
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

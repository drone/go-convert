package converter

import (
	"testing"
)

func TestConvertRemoteFilesWithWarnings_Basic(t *testing.T) {
	files := map[string]string{
		"values.yaml": "tag: <+pipeline.properties.ci.codebase.branch>",
		"config.yaml": "repo: <+pipeline.properties.ci.codebase.repoName>",
		"vars.yaml":   "env: <+serviceConfig.serviceDefinition.spec.variables.env>",
		"plain.txt":   "no expressions here",
		"empty.txt":   "",
	}

	got, _ := ConvertRemoteFilesWithWarnings(files, nil)

	want := map[string]string{
		"values.yaml": "tag: <+codebase.branch>",
		"config.yaml": "repo: <+codebase.repoName>",
		"vars.yaml":   "env: <+serviceVariables.env>",
		"plain.txt":   "no expressions here",
		"empty.txt":   "",
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d files, got %d", len(want), len(got))
	}
	for k, w := range want {
		if got[k] != w {
			t.Errorf("file %q:\n got:  %s\n want: %s", k, got[k], w)
		}
	}
}

func TestConvertRemoteFilesWithWarnings_EmptyInput(t *testing.T) {
	got, warnings := ConvertRemoteFilesWithWarnings(map[string]string{}, nil)
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(got))
	}
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %d", len(warnings))
	}
}

// TestConvertRemoteFilesWithWarnings_MatchesSingle ensures batch conversion of a
// single file produces the same output as the single-file API, so the batch
// path is behaviour-equivalent to ConvertExpressionWithWarnings.
func TestConvertRemoteFilesWithWarnings_MatchesSingle(t *testing.T) {
	content := "echo \"branch=<+pipeline.properties.ci.codebase.branch>\""

	batch, _ := ConvertRemoteFilesWithWarnings(map[string]string{"f": content}, nil)
	single, _ := ConvertExpressionWithWarnings(content, nil)

	if batch["f"] != single {
		t.Errorf("batch vs single mismatch\n batch:  %s\n single: %s", batch["f"], single)
	}
}

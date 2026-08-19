package converthelpers

import (
	"reflect"
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/convert/v0tov1/messagelog"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

func TestHelmCommandFlagsConversion(t *testing.T) {
	messagelog.ResetMessageLogger()
	logger := messagelog.GetMessageLogger()
	logger.Enable("")
	logger.SetCurrentFile("pipeline.yaml")
	defer messagelog.ResetMessageLogger()

	got := convertHelmCommandFlags([]v0.HelmCommandFlag{
		{CommandType: "Template", Flag: "--debug"},
		{CommandType: "Install", Flag: "--atomic"},
		{CommandType: "Delete", Flag: "--keep-history"},
		{CommandType: "Uninstall", Flag: "--wait"},
		{CommandType: "Fetch", Flag: "--devel"}, // no template option, lowercased passthrough
	})
	want := []HelmFlag{
		{Command: "template", Flag: "--debug"},
		{Command: "install", Flag: "--atomic"},
		{Command: "uninstall", Flag: "--keep-history"},
		{Command: "uninstall", Flag: "--wait"},
		{Command: "fetch", Flag: "--devel"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("convertHelmCommandFlags() = %v, want %v", got, want)
	}

	fileLog := logger.GetFileLog("pipeline.yaml")
	if fileLog == nil {
		t.Fatal("expected converter messages to be recorded")
	}
	found := false
	for _, m := range fileLog.Messages {
		if m.Code == "UNSUPPORTED_HELM_COMMAND_FLAG_TYPE" && m.Severity == messagelog.SeverityWarning {
			found = true
		}
	}
	if !found {
		t.Errorf("expected UNSUPPORTED_HELM_COMMAND_FLAG_TYPE warning, got %v", fileLog.Messages)
	}
}

func TestHelmDeleteCommandFlagsConversion(t *testing.T) {
	// v0 HelmDelete commandFlags is a plain []string; all flags belong to `helm uninstall`
	// and are joined into a single template flags entry.
	got := convertHelmDeleteCommandFlags([]string{"--debug", "", "--ignore-not-found"})
	want := []HelmFlag{{Command: "uninstall", Flag: "--debug --ignore-not-found"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("convertHelmDeleteCommandFlags() = %v, want %v", got, want)
	}

	if got := convertHelmDeleteCommandFlags(nil); got != nil {
		t.Errorf("convertHelmDeleteCommandFlags(nil) = %v, want nil", got)
	}
}

func TestConvertStepHelmDelete(t *testing.T) {
	step := &v0.Step{
		Spec: &v0.StepHelmDelete{
			ReleaseName:          "my-release",
			CommandFlags:         []string{"--debug", "--ignore-not-found"},
			EnvironmentVariables: map[string]string{"B": "2", "A": "1"},
		},
	}

	tmpl := ConvertStepHelmDelete(step)
	if tmpl == nil {
		t.Fatal("ConvertStepHelmDelete() returned nil")
	}
	if tmpl.Uses != v1.StepTypeHelmDelete {
		t.Errorf("Uses = %q, want %q", tmpl.Uses, v1.StepTypeHelmDelete)
	}

	with, ok := tmpl.With.(HelmDeleteWith)
	if !ok {
		t.Fatalf("With has unexpected type %T", tmpl.With)
	}
	if with.ReleaseName != "my-release" {
		t.Errorf("release = %q, want %q", with.ReleaseName, "my-release")
	}
	wantFlags := []HelmFlag{{Command: "uninstall", Flag: "--debug --ignore-not-found"}}
	if !reflect.DeepEqual(with.Flags, wantFlags) {
		t.Errorf("flags = %v, want %v", with.Flags, wantFlags)
	}
	wantEnv := []map[string]string{{"key": "A", "value": "1"}, {"key": "B", "value": "2"}}
	if !reflect.DeepEqual(with.Envvars, wantEnv) {
		t.Errorf("env_vars = %v, want %v", with.Envvars, wantEnv)
	}
}

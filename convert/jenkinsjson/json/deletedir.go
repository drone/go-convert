package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

func ConvertDeleteDir(node Node, variables map[string]string) *harness.Step {
	step := &harness.Step{
		Name: "DeleteDir",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "script",
		Spec: &harness.StepExec{
			Shell: "sh",
			Run:   "dir_to_delete=$(pwd) && cd .. && rm -rf $dir_to_delete",
		},
	}
	if len(variables) > 0 {
		step.Spec.(*harness.StepExec).Envs = variables
	}
	return step
}

func ConvertDir(node Node, variables map[string]string) (*harness.Step, bool) {
	dirPath, ok := node.ParameterMap["path"].(string)
	if !ok {
		return nil, false
	}
	step := &harness.Step{
		Name: "Dir",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "script",
		Spec: &harness.StepExec{
			Shell: "sh",
			Run:   fmt.Sprintf("mkdir -p %s && cd %s", dirPath, dirPath),
		},
	}
	if len(variables) > 0 {
		step.Spec.(*harness.StepExec).Envs = variables
	}
	return step, true
}

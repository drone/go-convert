package json

import (
	harness "github.com/drone/spec/dist/go"
)

func ConvertBuild(node Node, variables map[string]string, timeout string, dockerImage string) *harness.Step {
	shStep := ConvertSh(node, variables, timeout, dockerImage, "")

	stepExec := shStep.Spec.(*harness.StepExec)
	stepExec.Run = "curl https://replace_with_webhook_url"

	return shStep
}

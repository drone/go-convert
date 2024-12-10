package json

import (
	harness "github.com/drone/spec/dist/go"
)

var JenkinsToDroneSlackUserIdFromEmailParamMapperList = []JenkinsToDroneParamMapper{
	{"email", "slack_user_email_id", StringType, nil},
}

func ConvertSlackUserIdFromEmail(node Node, variables map[string]string) *harness.Step {
	step := ConvertToStepWithProperties(&node, variables, JenkinsToDroneSlackUserIdFromEmailParamMapperList,
		SlackPluginImage)
	tmpStepPlugin := step.Spec.(*harness.StepPlugin)
	tmpStepPlugin.With["access_token"] = "<+input>"
	return step
}

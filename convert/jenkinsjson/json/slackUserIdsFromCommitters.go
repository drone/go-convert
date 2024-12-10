package json

import harness "github.com/drone/spec/dist/go"

var JenkinsToDroneSlackUserIdsFromCommittersParamMapperList = []JenkinsToDroneParamMapper{
	{"email", "slack_user_email_id", StringType, nil},
}

func ConvertSlackUserIdsFromCommitters(node Node, variables map[string]string) *harness.Step {
	step := ConvertToStepWithProperties(&node, variables, JenkinsToDroneSlackUserIdsFromCommittersParamMapperList,
		SlackPluginImage)
	tmpStepPlugin := step.Spec.(*harness.StepPlugin)
	tmpStepPlugin.With["access_token"] = "<+input>"
	tmpStepPlugin.With["git_repo_path"] = "<+input>"
	tmpStepPlugin.With["committers_slack_id"] = true
	return step
}

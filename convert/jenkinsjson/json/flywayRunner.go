package json

import (
	harness "github.com/drone/spec/dist/go"
)

const (
	FlywayRunnerPluginImage = "plugins/drone-flyway-runner"
)

var JenkinsToFlywayRunnerParamMapperList = []JenkinsToDroneParamMapper{
	{"url", "url", StringType, nil},
	{"commandLineArgs", "command_line_args", StringType, nil},
	{"flywayCommand", "flyway_command", StringType, nil},
	{"locations", "locations", StringType, nil},
	{"commandLineArgs", "command_line_args", StringType, nil},
}

func ConvertFlywayRunner(node Node, variables map[string]string) *harness.Step {
	step := GetStepUsingParameterMapDelegate(&node, JenkinsToFlywayRunnerParamMapperList, FlywayRunnerPluginImage)
	tmpStepPlugin := step.Spec.(*harness.StepPlugin)
	tmpStepPlugin.With["username"] = "<+input>"
	tmpStepPlugin.With["password"] = "<+input>"

	return step
}

package json

import (
	harness "github.com/drone/spec/dist/go"
)

const (
	SlackPluginImage = "plugins/slack"
)

var JenkinsToDroneSlackUploadParamMapperList = []JenkinsToDroneParamMapper{
	{"filePath", "file_path", StringType, nil},
	{"channel", "channel", StringType, nil},
	{"credentialId", "credential_id", StringType, nil},
	{"initialComment", "initial_comment", StringType, nil},
	{"failOnError", "fail_on_error", BoolType, nil},
}

func ConvertSlackUploadFile(node Node, variables map[string]string) *harness.Step {

	step := ConvertToStepWithProperties(&node, variables, JenkinsToDroneSlackUploadParamMapperList,
		SlackPluginImage)

	return step
}

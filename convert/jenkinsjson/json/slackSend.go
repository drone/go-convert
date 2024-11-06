package json

import (
	harness "github.com/drone/spec/dist/go"
)

var JenkinsToDroneSlackSendParamMapperList = []JenkinsToDroneParamMapper{
	{"channel", "channel", StringType, nil},
	{"tokenCredentialId", "access_token", StringType, nil},
	{"username", "username", StringType, nil},
	{"color", "color", StringType, nil},
	{"message", "message", StringType, nil},
	{"iconEmoji", "icon_emoji", StringType, nil},
	// "teamDomain" e.g https://<teamDomain>.slack.com - Unsupported Jenkins field
}

func ConvertSlackSend(node Node, variables map[string]string) *harness.Step {
	step := ConvertToStepWithProperties(&node, variables, JenkinsToDroneSlackSendParamMapperList,
		SlackPluginImage)
	return step
}

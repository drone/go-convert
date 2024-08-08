package json

import (
	"fmt"
	harness "github.com/drone/spec/dist/go"
	"strings"
)

func ConvertSleep(node Node, variables map[string]string) *harness.Step {
	timeValue, ok := node.ParameterMap["time"].(float64)
	if !ok {
		// Handle error or set a default value
		timeValue = 0
	}

	unit, ok := node.ParameterMap["unit"].(string)
	if !ok {
		unit = ""
	}

	sleepCommand := fmt.Sprintf("sleep %.0f", timeValue)

	switch strings.ToUpper(unit) {
	case "SECONDS":
		sleepCommand = fmt.Sprintf("sleep %.0fs", timeValue)
	case "MINUTES":
		sleepCommand = fmt.Sprintf("sleep %.0fm", timeValue)
	case "HOURS":
		sleepCommand = fmt.Sprintf("sleep %.0fh", timeValue)
	case "DAYS":
		sleepCommand = fmt.Sprintf("sleep %.0fd", timeValue)
	case "NANOSECONDS", "MICROSECONDS", "MILLISECONDS":
		// Unsupported units, use default seconds
		sleepCommand = fmt.Sprintf("sleep %.0f", timeValue)
	}

	sleepStep := &harness.Step{
		Name: strings.TrimPrefix(node.SpanName, "Stage: "),
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "script",
		Spec: &harness.StepExec{
			Shell: "sh",
			Run:   sleepCommand,
			Image: "busybox",
		},
	}
	if len(variables) > 0 {
		sleepStep.Spec.(*harness.StepExec).Envs = variables
	}
	return sleepStep
}

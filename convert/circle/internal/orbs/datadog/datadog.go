package datadog

import (
	"fmt"
	"strings"

	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"
)

func Convert(command string, step *circle.Custom) *harness.Step {
	switch command {
	case "":
		return nil // not supported
	case "setup":
		return convertSetup(step)
	case "stop":
		return convertStop(step)
	default:
		return nil // not supported
	}
}

func convertSetup(step *circle.Custom) *harness.Step {
	version := "7" // default agent version

	envs := map[string]string{
		"DD_HOSTNAME":     "none",
		"DD_INSTALL_ONLY": "true",
		"DD_APM_ENABLED":  "true",
	}

	if s, _ := step.Params["site"].(string); s != "" {
		envs["DD_SITE"] = s
	}
	if s, _ := step.Params["agent_major_version"].(string); s != "" {
		envs["DD_AGENT_MAJOR_VERSION"] = s
		version = s
	}
	if s, _ := step.Params["api_key"].(string); s != "" {
		envs["DD_API_KEY"] = s
	}

	commands := []string{
		fmt.Sprintf(`bash -c "$(curl -L "https://s3.amazonaws.com/dd-agent/scripts/install_script_agent%s.sh")"`, version),
		`find /etc/datadog-agent/conf.d/ -iname "*.yaml.default" -delete`,
		`service datadog-agent start`,
		`echo "waiting for data-dog agent to start"`,
		`sleep 30`,
	}

	return &harness.Step{
		Name: "datadog_setup",
		Type: "script",
		Spec: &harness.StepBackground{
			Run:  strings.Join(commands, "\n"),
			Envs: envs,
		},
	}
}

func convertStop(step *circle.Custom) *harness.Step {
	commands := []string{
		// Stop Datadog Agent
		"service datadog-agent stop",
	}

	return &harness.Step{
		Name: "datadog_stop",
		Type: "script",
		Spec: &harness.StepBackground{
			Run: strings.Join(commands, "\n"),
		},
	}
}

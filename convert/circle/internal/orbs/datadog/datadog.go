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
	agentMajorVersion, _ := step.Params["agent_major_version"].(string)
	var scriptURL string

	switch agentMajorVersion {
	case "6":
		scriptURL = "https://s3.amazonaws.com/dd-agent/scripts/install_script_agent6.sh"
	case "7":
		scriptURL = "https://s3.amazonaws.com/dd-agent/scripts/install_script_agent7.sh"
	default:
		scriptURL = "https://s3.amazonaws.com/dd-agent/scripts/install_script_agent7.sh" // default to agent 7
	}
	commands := []string{
		// Install Datadog
		"PARAM_DD_API_KEY=$(eval echo \"\\$PARAM_DD_API_KEY\")",
		"if [ -n \"${DD_SITE}\" ]; then",
		"  PARAM_DD_SITE=${DD_SITE}",
		"fi",
		"DD_API_KEY=${PARAM_DD_API_KEY} DD_AGENT_MAJOR_VERSION=${PARAM_DD_AGENT_MAJOR_VERSION} DD_SITE=${PARAM_DD_SITE} \\",
		"DD_HOSTNAME=\"none\" DD_INSTALL_ONLY=\"true\" DD_APM_ENABLED=\"true\" \\",
		fmt.Sprintf("bash -c \"$(curl -L %s)\"", scriptURL),
		// Delete Default YAML Files
		"if [ \"$UID\" = \"0\" ]; then export SUDO=''; else export SUDO='sudo'; fi",
		"$SUDO find /etc/datadog-agent/conf.d/ -iname \"*.yaml.default\" -delete",
		// Start Datadog and Check Health
		"$SUDO service datadog-agent start",
		"set +e",
		"attempts=0",
		"until [ $attempts -eq 10 ] || $SUDO datadog-agent health; do",
		"attempts=$((attempts+1))",
		"sleep_time=$(( attempts*5 < 30 ? attempts*5 : 30 ))",
		"echo \"Waiting for agent to start up sleeping for ${sleep_time} seconds\"",
		"sleep $sleep_time",
		"done",
		"if [ $attempts -eq 10 ]; then",
		"echo \"Could not start the agent\"",
		"exit 1",
		"else",
		"echo \"Agent is ready\"",
		"fi",
	}

	envs := map[string]string{}
	if s, _ := step.Params["agent_major_version"].(string); s != "" {
		envs["PARAM_DD_AGENT_MAJOR_VERSION"] = s
	}
	if s, _ := step.Params["api_key"].(string); s != "" {
		envs["PARAM_DD_API_KEY"] = s
	}
	if s, _ := step.Params["site"].(string); s != "" {
		envs["PARAM_DD_SITE"] = s
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
		"if [ \"$UID\" = \"0\" ]; then export SUDO=''; else export SUDO='sudo'; fi",
		"$SUDO service datadog-agent stop",
	}

	return &harness.Step{
		Name: "datadog_stop",
		Type: "script",
		Spec: &harness.StepBackground{
			Run: strings.Join(commands, "\n"),
		},
	}
}

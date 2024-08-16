package json

import (
	"encoding/json"
	"fmt"
	harness "github.com/drone/spec/dist/go"
	"strconv"
)

type anchoreArguments struct {
	BailOnFail              string `json:"bailOnFail"`
	ForceAnalyze            bool   `json:"forceAnalyze"`
	Name                    string `json:"name"`
	PolicyBundleId          string `json:"policyBundleId"`
	EngineCredentialsId     string `json:"engineCredentialsId"`
	EngineRetries           string `json:"engineRetries"`
	Engineurl               string `json:"engineurl"`
	Engineverify            bool   `json:"engineverify"`
	ExcludeFromBaseImage    bool   `json:"excludeFromBaseImage"`
	BailOnPluginFail        bool   `json:"bailOnPluginFail"`
	AutoSubscribeTagUpdates bool   `json:"autoSubscribeTagUpdates"`
	Engineaccount           string `json:"engineaccount"`
}

func (a *anchoreArguments) UnmarshalJSON(data []byte) error {
	type Alias anchoreArguments
	aux := &struct {
		ForceAnalyze            string `json:"forceAnalyze"`
		Engineverify            bool   `json:"engineverify"`
		ExcludeFromBaseImage    bool   `json:"excludeFromBaseImage"`
		BailOnPluginFail        bool   `json:"bailOnPluginFail"`
		AutoSubscribeTagUpdates bool   `json:"autoSubscribeTagUpdates"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	a.ForceAnalyze, err = strconv.ParseBool(aux.ForceAnalyze)
	if err != nil {
		return err
	}
	return nil
}

func ConvertAnchore(node Node, variables map[string]string) *harness.Step {
	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var anchoreAttr struct {
			Delegate struct {
				Arguments anchoreArguments `json:"arguments"`
			} `json:"delegate"`
		}

		err := json.Unmarshal([]byte(attr), &anchoreAttr)
		if err != nil {
			fmt.Println("Error parsing anchore attribute:", err)
			return nil
		}

		args := anchoreAttr.Delegate.Arguments

		baseCmd := "anchorectl"
		if args.Engineurl != "" {
			baseCmd += fmt.Sprintf(" --url %s", args.Engineurl)
		}
		if args.EngineCredentialsId != "" {
			baseCmd += fmt.Sprintf(" --username %s", args.EngineCredentialsId)
		}
		if args.EngineRetries != "" {
			baseCmd += fmt.Sprintf(" --max-retries %s", args.EngineRetries)
		}
		if args.Engineverify {
			baseCmd += " --insecure-skip-verify"
		}
		if args.Engineaccount != "" {
			baseCmd += fmt.Sprintf(" --account %s", args.Engineaccount)
		}

		cmd := fmt.Sprintf(`
set -e
curl -sSfL https://anchorectl-releases.anchore.io/anchorectl/install.sh | sh -s -- -b ${HOME}/.local/bin
export PATH="${HOME}/.local/bin/:${PATH}"
anchorectl --version

ANCHORE_IMAGE=$(cat $ANCHORE_FILE_NAME)

%s image add --wait $ANCHORE_IMAGE

%s image vulnerabilities`, baseCmd, baseCmd)

		if args.ExcludeFromBaseImage {
			cmd += " --exclude-from-base"
		}
		cmd += " $ANCHORE_IMAGE\n\n"

		cmd += baseCmd + " image check --detail"
		if args.PolicyBundleId != "" {
			cmd += fmt.Sprintf(" --policy %s", args.PolicyBundleId)
		}
		if args.ForceAnalyze {
			cmd += " --force"
		}
		if args.BailOnPluginFail {
			cmd += " --fail-on-plugin-error"
		}
		cmd += " $ANCHORE_IMAGE\n\n"

		if args.AutoSubscribeTagUpdates {
			cmd += baseCmd + " subscription activate $ANCHORE_IMAGE\n\n"
		}

		if args.BailOnFail == "true" {
			cmd += "exit $?\n"
		} else {
			cmd += "exit 0\n"
		}

		envs := map[string]string{
			"ANCHORECTL_FAIL_BASED_ON_RESULTS":   anchoreAttr.Delegate.Arguments.BailOnFail,
			"ANCHORECTL_FORCE":                   strconv.FormatBool(anchoreAttr.Delegate.Arguments.ForceAnalyze),
			"ANCHORE_FILE_NAME":                  anchoreAttr.Delegate.Arguments.Name,
			"ANCHORECTL_POLICY":                  anchoreAttr.Delegate.Arguments.PolicyBundleId,
			"ANCHORECTL_ENGINECREDENTIALS":       anchoreAttr.Delegate.Arguments.EngineCredentialsId,
			"ANCHORECTL_ENGINERETRIES":           anchoreAttr.Delegate.Arguments.EngineRetries,
			"ANCHORECTL_URL":                     anchoreAttr.Delegate.Arguments.Engineurl,
			"ANCHORECTL_ENGINEVERIFY":            strconv.FormatBool(anchoreAttr.Delegate.Arguments.Engineverify),
			"ANCHORECTL_EXCLUDEFROMBASEIMAGE":    strconv.FormatBool(anchoreAttr.Delegate.Arguments.ExcludeFromBaseImage),
			"ANCHORECTL_BAILONPLUGINFAIL":        strconv.FormatBool(anchoreAttr.Delegate.Arguments.BailOnPluginFail),
			"ANCHORECTL_AUTOSUBSCRIBETAGUPDATES": strconv.FormatBool(anchoreAttr.Delegate.Arguments.AutoSubscribeTagUpdates),
			"ANCHORECTL_ENGINEACCOUNT":           anchoreAttr.Delegate.Arguments.Engineaccount,
		}

		for k, v := range envs {
			if v == "" {
				delete(envs, k)
			}
		}

		for k, v := range variables {
			envs[k] = v
		}

		//runCommand := fmt.Sprintf("curl -sSfL https://anchorectl-releases.anchore.io/anchorectl/install.sh | sh -s -- -b ${HOME}/.local/bin\nexport PATH=\"${HOME}/.local/bin/:${PATH}\"\nanchorectl --version\n\nANCHORE_IMAGE=$(cat $ANCHORE_FILE_NAME)\n\nANCHORE_CMD=\"anchorectl\"\n[ -n \"$ANCHORECTL_URL\" ] && ANCHORE_CMD+=\" --url $ANCHORECTL_URL\"\n[ -n \"$ANCHORECTL_ENGINECREDENTIALS\" ] && ANCHORE_CMD+=\" --username $ANCHORECTL_ENGINECREDENTIALS\"\n[ -n \"$ANCHORECTL_ENGINERETRIES\" ] && ANCHORE_CMD+=\" --max-retries $ANCHORECTL_ENGINERETRIES\"\n[ \"$ANCHORECTL_ENGINEVERIFY\" = \"true\" ] && ANCHORE_CMD+=\" --insecure-skip-verify\"\n[ -n \"$ANCHORECTL_ENGINEACCOUNT\" ] && ANCHORE_CMD+=\" --account $ANCHORECTL_ENGINEACCOUNT\"\n\n$ANCHORE_CMD image add --wait $ANCHORE_IMAGE\n\nVULN_CMD=\"$ANCHORE_CMD image vulnerabilities\"\n[ \"$ANCHORECTL_EXCLUDEFROMBASEIMAGE\" = \"true\" ] && VULN_CMD+=\" --exclude-from-base\"\n$VULN_CMD $ANCHORE_IMAGE\n\nCHECK_CMD=\"$ANCHORE_CMD image check --detail\"\n[ -n \"$ANCHORECTL_POLICY\" ] && CHECK_CMD+=\" --policy $ANCHORECTL_POLICY\"\n[ \"$ANCHORECTL_FORCE\" = \"true\" ] && CHECK_CMD+=\" --force\"\n[ \"$ANCHORECTL_BAILONPLUGINFAIL\" = \"true\" ] && CHECK_CMD+=\" --fail-on-plugin-error\"\n$CHECK_CMD $ANCHORE_IMAGE\n\n[ \"$ANCHORECTL_AUTOSUBSCRIBETAGUPDATES\" = \"true\" ] && $ANCHORE_CMD subscription activate $ANCHORE_IMAGE\n\nexit_code=$?\n[ \"$ANCHORECTL_FAIL_BASED_ON_RESULTS\" = \"true\" ] && exit $exit_code || exit 0")

		step := &harness.Step{
			Name: node.SpanName,
			Id:   SanitizeForId(node.SpanName, node.SpanId),
			Type: "script",
			Spec: &harness.StepExec{
				Shell: "sh",
				Run:   cmd,
				Envs:  envs,
			},
		}

		return step
	}
	fmt.Println("Error: harness-attribute not found for anchore step")
	return nil
}

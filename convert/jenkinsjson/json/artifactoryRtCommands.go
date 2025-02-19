package json

import (
	"errors"
	"fmt"
	harness "github.com/drone/spec/dist/go"
)

const (
	ArtifactoryRtCommandsPluginImage = "plugins/artifactory"
	InputPlaceHolder                 = "<+input>"
	MvnTool                          = "mvn"
	GradleTool                       = "gradle"
	UrlStr                           = "url"
	UseNameStr                       = "username"
	AccessTokenStr                   = "access_token"
	ModuleStr                        = "module"
	ResolverIdStr                    = "resolver_id"
	DeployerIdStr                    = "deployer_id"
	ProjectStr                       = "project"
	BuildNameStr                     = "build_name"
	BuildNumberStr                   = "build_number"
	BuildToolStr                     = "build_tool"
	SpecPath                         = "spec_path"
	TargetStr                        = "target"
	CopyStr                          = "copy"

	// command defs
	rtDownloadCmd       = "rtDownload"
	rtMavenRunCmd       = "rtMavenRun"
	rtGradleRunCmd      = "rtGradleRun"
	publishBuildInfoCmd = "rtPublishBuildInfo"
	rtPromoteCmd        = "rtPromote"
	xrayScanCmd         = "xrayScan"
)

var (
	// rtDownloadAttributesList defines
	rtDownloadAttributesList = []string{UrlStr, UseNameStr, AccessTokenStr, ModuleStr,
		ProjectStr, BuildNameStr, BuildNumberStr, SpecPath}
	rtMavenRunAttributesList = []string{UrlStr, UseNameStr, AccessTokenStr, ResolverIdStr, DeployerIdStr,
		"resolve_release_repo", "resolve_snapshot_repo", BuildNameStr, BuildNumberStr}
	rtGradleRunAttributesList = []string{UrlStr, UseNameStr, AccessTokenStr, ResolverIdStr, DeployerIdStr,
		"repo_resolve", "repo_deploy", BuildNameStr, BuildNumberStr}
	publishBuildInfoAttributesList = []string{BuildToolStr, UrlStr, UseNameStr, AccessTokenStr, BuildNameStr,
		BuildNumberStr, DeployerIdStr, "deploy_release_repo", "deploy_snapshot_repo"}
	rtPromoteAttributesList = []string{UrlStr, UseNameStr, AccessTokenStr, BuildNameStr,
		BuildNumberStr, TargetStr, CopyStr}
	xrayScanAttributesList = []string{UrlStr, UseNameStr, AccessTokenStr, BuildNameStr, BuildNumberStr}

	// ConvertRtMavenRunParamMapperList defines
	ConvertRtMavenRunParamMapperList = []JenkinsToDroneParamMapper{
		{"pom", "source", StringType, nil},
		{"goals", "goals", StringType, nil},
		{"buildName", BuildNameStr, StringType, nil},
		{"buildNumber", BuildNumberStr, StringType, nil},
	}
	ConvertRtGradleRunParamMapperList = []JenkinsToDroneParamMapper{
		{"buildName", BuildNameStr, StringType, nil},
		{"buildNumber", BuildNumberStr, StringType, nil},
		{"tasks", "tasks", StringType, nil},
	}
	ConvertRtDownloadParamMapperList = []JenkinsToDroneParamMapper{
		{"buildName", BuildNameStr, StringType, nil},
		{"buildNumber", BuildNumberStr, StringType, nil},
		{ModuleStr, ModuleStr, StringType, nil},
		{"specPath", SpecPath, StringType, nil},
	}
	ConvertRtPromoteParamMapperList = []JenkinsToDroneParamMapper{
		{"buildName", BuildNameStr, StringType, nil},
		{"buildNumber", BuildNumberStr, StringType, nil},
		{"targetRepo", TargetStr, StringType, nil},
		{CopyStr, CopyStr, StringType, nil},
	}
	ConvertXrayScanParamMapperList = []JenkinsToDroneParamMapper{
		{"buildName", BuildNameStr, StringType, nil},
		{"buildNumber", BuildNumberStr, StringType, nil},
	}
)

type RtCommandParams struct {
	StepType                      string
	Tool                          string
	Command                       string
	JenkinsToDroneParamMapperList []JenkinsToDroneParamMapper
	AttributesList                []string
}

var RtCommandParamsDef = map[string]RtCommandParams{
	rtDownloadCmd:       {rtDownloadCmd, "", "download", ConvertRtDownloadParamMapperList, rtDownloadAttributesList},
	rtMavenRunCmd:       {rtMavenRunCmd, MvnTool, "", ConvertRtMavenRunParamMapperList, rtMavenRunAttributesList},
	rtGradleRunCmd:      {rtGradleRunCmd, GradleTool, "", ConvertRtGradleRunParamMapperList, rtGradleRunAttributesList},
	publishBuildInfoCmd: {publishBuildInfoCmd, "", "publish", nil, publishBuildInfoAttributesList},
	rtPromoteCmd:        {rtPromoteCmd, "", "promote", ConvertRtPromoteParamMapperList, rtPromoteAttributesList},
	xrayScanCmd:         {xrayScanCmd, "", "scan", ConvertXrayScanParamMapperList, xrayScanAttributesList},
}

func ConvertArtifactoryRtCommand(stepType string, node Node, variables map[string]string) *harness.Step {
	rtCommandParams, ok := RtCommandParamsDef[stepType]
	if !ok {
		fmt.Println("Error: failed to convert ", stepType)
		return nil
	}
	step := convertRtStep(rtCommandParams.StepType, rtCommandParams.Tool, rtCommandParams.Command, node,
		rtCommandParams.JenkinsToDroneParamMapperList, rtCommandParams.AttributesList)
	return step
}

func convertRtStep(stepType string, tool string, command string, node Node,
	jenkinsToDroneParamMapper []JenkinsToDroneParamMapper, attributesList []string) *harness.Step {
	step := GetStepWithProperties(&node, jenkinsToDroneParamMapper, ArtifactoryRtCommandsPluginImage)
	if step == nil {
		fmt.Println("Error: failed to convert ", stepType)
		return nil
	}
	tmpStepPlugin, ok := step.Spec.(*harness.StepPlugin)
	if !ok {
		fmt.Println("Error: failed to convert to StepPlugin ", stepType)
		return nil
	}
	if tmpStepPlugin.With == nil {
		tmpStepPlugin.With = map[string]interface{}{}
	}
	err := SetRtCommandAttributes(tool, command, tmpStepPlugin, attributesList)
	if err != nil {
		fmt.Println("Error: failed to set attributes to input placeholder")
		return nil
	}
	return step
}

func SetRtCommandAttributes(toolName string, command string,
	tmpStepPlugin *harness.StepPlugin, attributeValues []string) error {
	if tmpStepPlugin == nil {
		errStr := "error: rtCommand StepPlugin is nil"
		fmt.Println(errStr)
		return errors.New(errStr)
	}

	if toolName != "" {
		tmpStepPlugin.With[BuildToolStr] = toolName
	}

	if command != "" {
		tmpStepPlugin.With["command"] = command
	}
	for _, attribute := range attributeValues {
		if _, ok := tmpStepPlugin.With[attribute]; !ok {
			tmpStepPlugin.With[attribute] = InputPlaceHolder
		}
	}
	return nil
}

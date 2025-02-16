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
)

func ConvertArtifactoryRtCommand(stepType string, node Node, variables map[string]string) *harness.Step {
	switch stepType {
	case "rtDownload":
		return convertRtDownload(node, variables)
	case "rtMavenRun":
		return convertRtMavenRun(node, variables)
	case "rtGradleRun":
		return convertRtGradleRun(node, variables)
	case "publishBuildInfo":
		return convertPublishBuildInfo(node, variables)
	case "rtPromote":
		return convertRtPromote(node, variables)
	case "xrayScan":
		return convertXrayScan(node, variables)
	}
	return nil
}

var ConvertRtMavenRunParamMapperList = []JenkinsToDroneParamMapper{
	{"pom", "source", StringType, nil},
	{"goals", "goals", StringType, nil},
	{"buildName", "build_name", StringType, nil},
	{"buildNumber", "build_number", StringType, nil},
}

func convertRtMavenRun(node Node, variables map[string]string) *harness.Step {
	step := GetStepWithProperties(&node, ConvertRtMavenRunParamMapperList, ArtifactoryRtCommandsPluginImage)
	if step == nil {
		fmt.Println("Error: failed to convert rtMavenRun")
		return nil
	}
	tmpStepPlugin, ok := step.Spec.(*harness.StepPlugin)
	if !ok {
		fmt.Println("Error: failed to convert rtMavenRun")
		return nil
	}
	if tmpStepPlugin.With == nil {
		tmpStepPlugin.With = map[string]interface{}{}
	}
	tmpStepPlugin.With["build_tool"] = MvnTool
	attributesList := []string{"url", "username", "password", "access_token",
		"resolver_id", "deployer_id", "resolve_release_repo", "resolve_snapshot_repo"}

	if _, ok := tmpStepPlugin.With["build_name"]; !ok {
		attributesList = append(attributesList, "build_name")
	}
	if _, ok := tmpStepPlugin.With["build_number"]; !ok {
		attributesList = append(attributesList, "build_number")
	}

	err := SetRtCommandAttributesToInputPlaceHolder(tmpStepPlugin, attributesList)
	if err != nil {
		fmt.Println("Error: failed to set attributes to input placeholder")
		return nil
	}
	return step
}

var ConvertRtGradleRunParamMapperList = []JenkinsToDroneParamMapper{
	{"buildName", "build_name", StringType, nil},
	{"buildNumber", "build_number", StringType, nil},
	{"tasks", "tasks", StringType, nil},
}

func convertRtGradleRun(node Node, variables map[string]string) *harness.Step {
	step := GetStepWithProperties(&node, ConvertRtGradleRunParamMapperList, ArtifactoryRtCommandsPluginImage)
	if step == nil {
		fmt.Println("Error: failed to convert rtGradleRun")
		return nil
	}
	tmpStepPlugin, ok := step.Spec.(*harness.StepPlugin)
	if !ok {
		fmt.Println("Error: failed to convert rtGradleRun")
		return nil
	}
	if tmpStepPlugin.With == nil {
		tmpStepPlugin.With = map[string]interface{}{}
	}

	tmpStepPlugin.With["build_tool"] = GradleTool

	attributesList := []string{"url", "username", "password", "access_token", "build_name",
		"build_number", "resolver_id", "deployer_id", "repo_resolve", "repo_deploy"}
	if _, ok := tmpStepPlugin.With["build_name"]; !ok {
		attributesList = append(attributesList, "build_name")
	}
	if _, ok := tmpStepPlugin.With["build_number"]; !ok {
		attributesList = append(attributesList, "build_number")
	}
	err := SetRtCommandAttributesToInputPlaceHolder(tmpStepPlugin, attributesList)
	if err != nil {
		fmt.Println("Error: failed to set attributes to input placeholder")
		return nil
	}
	return step
}

func convertPublishBuildInfo(node Node, variables map[string]string) *harness.Step {
	step := GetStepWithProperties(&node, nil, ArtifactoryRtCommandsPluginImage)
	if step == nil {
		fmt.Println("Error: failed to convert publishBuildInfo")
		return nil
	}
	tmpStepPlugin, ok := step.Spec.(*harness.StepPlugin)
	if !ok {
		fmt.Println("Error: failed to convert publishBuildInfo")
		return nil
	}
	if tmpStepPlugin.With == nil {
		tmpStepPlugin.With = map[string]interface{}{}
	}
	tmpStepPlugin.With["command"] = "publish"
	attributesList := []string{"build_tool", "url", "username", "password", "access_token", "build_name",
		"build_number", "deployer_id", "deploy_release_repo", "deploy_snapshot_repo"}

	err := SetRtCommandAttributesToInputPlaceHolder(tmpStepPlugin, attributesList)
	if err != nil {
		fmt.Println("Error: failed to set attributes to input placeholder")
		return nil
	}
	return step
}

var ConvertRtDownloadParamMapperList = []JenkinsToDroneParamMapper{
	{"buildName", "build_name", StringType, nil},
	{"buildNumber", "build_number", StringType, nil},
	{"module", "module", StringType, nil},
	{"specPath", "spec_path", StringType, nil},
}

func convertRtDownload(node Node, variables map[string]string) *harness.Step {
	step := GetStepWithProperties(&node, ConvertRtDownloadParamMapperList, ArtifactoryRtCommandsPluginImage)
	if step == nil {
		fmt.Println("Error: failed to convert rtDownload")
		return nil
	}
	tmpStepPlugin, ok := step.Spec.(*harness.StepPlugin)
	if !ok {
		fmt.Println("Error: failed to convert rtDownload")
		return nil
	}
	if tmpStepPlugin.With == nil {
		tmpStepPlugin.With = map[string]interface{}{}
	}

	tmpStepPlugin.With["command"] = "download"
	attributesList := []string{"url", "username", "password", "module", "project"}
	if _, ok := tmpStepPlugin.With["build_name"]; !ok {
		attributesList = append(attributesList, "build_name")
	}
	if _, ok := tmpStepPlugin.With["build_number"]; !ok {
		attributesList = append(attributesList, "build_number")
	}
	if _, ok := tmpStepPlugin.With["spec_path"]; !ok {
		attributesList = append(attributesList, "spec_path")
	}

	err := SetRtCommandAttributesToInputPlaceHolder(tmpStepPlugin, attributesList)
	if err != nil {
		fmt.Println("Error: failed to set attributes to input placeholder")
		return nil
	}
	return step
}

var ConvertRtPromoteParamMapperList = []JenkinsToDroneParamMapper{
	{"buildName", "build_name", StringType, nil},
	{"buildNumber", "build_number", StringType, nil},
	{"targetRepo", "target", StringType, nil},
	{"copy", "copy", StringType, nil},
}

func convertRtPromote(node Node, variables map[string]string) *harness.Step {
	step := GetStepWithProperties(&node, ConvertRtPromoteParamMapperList, ArtifactoryRtCommandsPluginImage)
	if step == nil {
		fmt.Println("Error: failed to convert rtPromote")
		return nil
	}
	tmpStepPlugin, ok := step.Spec.(*harness.StepPlugin)
	if !ok {
		fmt.Println("Error: failed to convert rtPromote")
		return nil
	}
	if tmpStepPlugin.With == nil {
		tmpStepPlugin.With = map[string]interface{}{}
	}
	tmpStepPlugin.With["command"] = "promote"
	attributesList := []string{"url", "username", "password", "access_token"}

	if _, ok := tmpStepPlugin.With["build_name"]; !ok {
		attributesList = append(attributesList, "build_name")
	}
	if _, ok := tmpStepPlugin.With["build_number"]; !ok {
		attributesList = append(attributesList, "build_number")
	}
	if _, ok := tmpStepPlugin.With["target"]; !ok {
		attributesList = append(attributesList, "target")
	}
	if _, ok := tmpStepPlugin.With["copy"]; !ok {
		attributesList = append(attributesList, "copy")
	}

	err := SetRtCommandAttributesToInputPlaceHolder(tmpStepPlugin, attributesList)
	if err != nil {
		fmt.Println("Error: failed to set attributes to input placeholder")
		return nil
	}
	return step
}

var ConvertXrayScanParamMapperList = []JenkinsToDroneParamMapper{
	{"buildName", "build_name", StringType, nil},
	{"buildNumber", "build_number", StringType, nil},
}

func convertXrayScan(node Node, variables map[string]string) *harness.Step {
	step := GetStepWithProperties(&node, ConvertXrayScanParamMapperList, ArtifactoryRtCommandsPluginImage)
	if step == nil {
		fmt.Println("Error: failed to convert xrayScan")
		return nil
	}
	tmpStepPlugin, ok := step.Spec.(*harness.StepPlugin)
	if !ok {
		fmt.Println("Error: failed to convert xrayScan")
		return nil
	}
	if tmpStepPlugin.With == nil {
		tmpStepPlugin.With = map[string]interface{}{}
	}
	tmpStepPlugin.With["command"] = "scan"
	if _, ok := tmpStepPlugin.With["build_name"]; !ok {
		tmpStepPlugin.With["build_name"] = InputPlaceHolder
	}
	if _, ok := tmpStepPlugin.With["build_number"]; !ok {
		tmpStepPlugin.With["build_number"] = InputPlaceHolder
	}
	attributesList := []string{"url", "username", "password", "access_token"}
	err := SetRtCommandAttributesToInputPlaceHolder(tmpStepPlugin, attributesList)
	if err != nil {
		fmt.Println("Error: failed to set attributes to input placeholder")
		return nil
	}
	return step
}

func SetRtCommandAttributesToInputPlaceHolder(tmpStepPlugin *harness.StepPlugin, attributeValues []string) error {
	if tmpStepPlugin == nil {
		errStr := "error: rtCommand StepPlugin is nil"
		fmt.Println(errStr)
		return errors.New(errStr)
	}

	for _, attribute := range attributeValues {
		tmpStepPlugin.With[attribute] = InputPlaceHolder
	}
	return nil
}

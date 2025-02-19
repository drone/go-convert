package json

import (
	"encoding/json"
	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

var expectedRtDownloadStep = `{
    "id": "rtUploadc4112c",
    "name": "rtUpload",
    "type": "plugin",
    "spec": {
        "image": "plugins/artifactory",
        "with": {
            "build_name": "\u003c+input\u003e",
            "build_number": "\u003c+input\u003e",
            "command": "download",
            "module": "\u003c+input\u003e",
            "access_token": "\u003c+input\u003e",
            "project": "\u003c+input\u003e",
            "spec_path": "\u003c+input\u003e",
            "url": "\u003c+input\u003e",
            "username": "\u003c+input\u003e"
        },
        "inputs": {
            "build_name": "\u003c+input\u003e",
            "build_number": "\u003c+input\u003e",
            "command": "download",
            "module": "\u003c+input\u003e",
            "access_token": "\u003c+input\u003e",
            "project": "\u003c+input\u003e",
            "spec_path": "\u003c+input\u003e",
            "url": "\u003c+input\u003e",
            "username": "\u003c+input\u003e"
        }
    }
}`

func TestRtDownload(t *testing.T) {
	jsonFilePath := "../convertTestFiles/artifactoryRtCommands/rtDownload.json"

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, jsonFilePath)

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var node Node
	if err := json.Unmarshal(jsonData, &node); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	tmpTestStep := convertRtStep("rtDownload", "", "download", node,
		ConvertRtDownloadParamMapperList, rtDownloadAttributesList)

	wantStep, err := ToStructFromJsonString[harness.Step](expectedRtDownloadStep)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("Failed to convert JSON to struct: %v", diffs)
	}
}

const expectedRtMavenRunStep = `{
    "id": "rtMavenRundac366",
    "name": "rtMavenRun",
    "type": "plugin",
    "spec": {
        "image": "plugins/artifactory",
        "with": {
            "access_token": "\u003c+input\u003e",
            "build_name": "\u003c+input\u003e",
            "build_number": "\u003c+input\u003e",
            "build_tool": "mvn",
            "deployer_id": "\u003c+input\u003e",
            "goals": "clean install",
            "resolve_release_repo": "\u003c+input\u003e",
            "resolve_snapshot_repo": "\u003c+input\u003e",
            "resolver_id": "\u003c+input\u003e",
            "source": "pom.xml",
            "url": "\u003c+input\u003e",
            "username": "\u003c+input\u003e"
        },
        "inputs": {
            "access_token": "\u003c+input\u003e",
            "build_name": "\u003c+input\u003e",
            "build_number": "\u003c+input\u003e",
            "build_tool": "mvn",
            "deployer_id": "\u003c+input\u003e",
            "goals": "clean install",
            "resolve_release_repo": "\u003c+input\u003e",
            "resolve_snapshot_repo": "\u003c+input\u003e",
            "resolver_id": "\u003c+input\u003e",
            "source": "pom.xml",
            "url": "\u003c+input\u003e",
            "username": "\u003c+input\u003e"
        }
    }
}`

func TestRtMavenRun(t *testing.T) {
	jsonFilePath := "../convertTestFiles/artifactoryRtCommands/rtMavenRun.json"

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, jsonFilePath)

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var node Node
	if err := json.Unmarshal(jsonData, &node); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	tmpTestStep := convertRtStep("rtMavenRun", "mvn", "",
		node, ConvertRtMavenRunParamMapperList, rtMavenRunAttributesList)
	wantStep, err := ToStructFromJsonString[harness.Step](expectedRtMavenRunStep)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("Failed to convert JSON to struct: %v", diffs)
	}
}

const expectedRtGradleRunStep = `{
    "id": "rtGradleRun110fc6",
    "name": "rtGradleRun",
    "type": "plugin",
    "spec": {
        "image": "plugins/artifactory",
        "with": {
            "access_token": "\u003c+input\u003e",
            "build_name": "\u003c+input\u003e",
            "build_number": "\u003c+input\u003e",
            "build_tool": "gradle",
            "deployer_id": "\u003c+input\u003e",
            "repo_deploy": "\u003c+input\u003e",
            "repo_resolve": "\u003c+input\u003e",
            "resolver_id": "\u003c+input\u003e",
            "tasks": "clean build",
            "url": "\u003c+input\u003e",
            "username": "\u003c+input\u003e"
        },
        "inputs": {
            "access_token": "\u003c+input\u003e",
            "build_name": "\u003c+input\u003e",
            "build_number": "\u003c+input\u003e",
            "build_tool": "gradle",
            "deployer_id": "\u003c+input\u003e",
            "repo_deploy": "\u003c+input\u003e",
            "repo_resolve": "\u003c+input\u003e",
            "resolver_id": "\u003c+input\u003e",
            "tasks": "clean build",
            "url": "\u003c+input\u003e",
            "username": "\u003c+input\u003e"
        }
    }
}`

func TestRtGradleRun(t *testing.T) {
	jsonFilePath := "../convertTestFiles/artifactoryRtCommands/rtGradleRun.json"

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, jsonFilePath)

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var node Node
	if err := json.Unmarshal(jsonData, &node); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	tmpTestStep := convertRtStep("rtGradleRun", "gradle", "", node,
		ConvertRtGradleRunParamMapperList, rtGradleRunAttributesList)
	wantStep, err := ToStructFromJsonString[harness.Step](expectedRtGradleRunStep)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("Failed to convert JSON to struct: %v", diffs)
	}
}

const expectedPublishBuildInfoStep = `{
    "id": "publishBuildInfo66c280",
    "name": "publishBuildInfo",
    "type": "plugin",
    "spec": {
        "image": "plugins/artifactory",
        "with": {
            "access_token": "\u003c+input\u003e",
            "build_name": "\u003c+input\u003e",
            "build_number": "\u003c+input\u003e",
            "build_tool": "\u003c+input\u003e",
            "command": "publish",
            "deploy_release_repo": "\u003c+input\u003e",
            "deploy_snapshot_repo": "\u003c+input\u003e",
            "deployer_id": "\u003c+input\u003e",
            "url": "\u003c+input\u003e",
            "username": "\u003c+input\u003e"
        },
        "inputs": {
            "access_token": "\u003c+input\u003e",
            "build_name": "\u003c+input\u003e",
            "build_number": "\u003c+input\u003e",
            "build_tool": "\u003c+input\u003e",
            "command": "publish",
            "deploy_release_repo": "\u003c+input\u003e",
            "deploy_snapshot_repo": "\u003c+input\u003e",
            "deployer_id": "\u003c+input\u003e",
            "url": "\u003c+input\u003e",
            "username": "\u003c+input\u003e"
        }
    }
}`

func TestPublishBuildInfo(t *testing.T) {
	jsonFilePath := "../convertTestFiles/artifactoryRtCommands/publishBuildInfo.json"

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, jsonFilePath)

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var node Node
	if err := json.Unmarshal(jsonData, &node); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	tmpTestStep := convertRtStep("publishBuildInfo", "", "", node, nil, publishBuildInfoAttributesList)
	wantStep, err := ToStructFromJsonString[harness.Step](expectedPublishBuildInfoStep)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("Failed to convert JSON to struct: %v", diffs)
	}
}

const expectedRtPromoteStep = `{
    "id": "rtPromoted6c65e",
    "name": "rtPromote",
    "type": "plugin",
    "spec": {
        "image": "plugins/artifactory",
        "with": {
            "access_token": "\u003c+input\u003e",
            "build_name": "mvn02",
            "build_number": "2",
            "command": "promote",
            "copy": "\u003c+input\u003e",
            "target": "tst-libs-snapshot",
            "url": "\u003c+input\u003e",
            "username": "\u003c+input\u003e"
        },
        "inputs": {
            "access_token": "\u003c+input\u003e",
            "build_name": "mvn02",
            "build_number": "2",
            "command": "promote",
            "copy": "\u003c+input\u003e",
            "target": "tst-libs-snapshot",
            "url": "\u003c+input\u003e",
            "username": "\u003c+input\u003e"
        }
    }
}`

func TestRtPromote(t *testing.T) {
	jsonFilePath := "../convertTestFiles/artifactoryRtCommands/rtPromote.json"

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, jsonFilePath)

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var node Node
	if err := json.Unmarshal(jsonData, &node); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	tmpTestStep := convertRtStep("rtPromote", "", "", node, ConvertRtPromoteParamMapperList, rtPromoteAttributesList)
	wantStep, err := ToStructFromJsonString[harness.Step](expectedRtPromoteStep)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("Failed to convert JSON to struct: %v", diffs)
	}
}

const expectedXrayScanStep = `{
    "id": "xrayScan673cd5",
    "name": "xrayScan",
    "type": "plugin",
    "spec": {
        "image": "plugins/artifactory",
        "with": {
            "access_token": "\u003c+input\u003e",
            "build_name": "mvn02",
            "build_number": "2",
            "command": "scan",
            "url": "\u003c+input\u003e",
            "username": "\u003c+input\u003e"
        },
        "inputs": {
            "access_token": "\u003c+input\u003e",
            "build_name": "mvn02",
            "build_number": "2",
            "command": "scan",
            "url": "\u003c+input\u003e",
            "username": "\u003c+input\u003e"
        }
    }
}`

func TestXrayScan(t *testing.T) {
	jsonFilePath := "../convertTestFiles/artifactoryRtCommands/xrayScan.json"

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, jsonFilePath)

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var node Node
	if err := json.Unmarshal(jsonData, &node); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	tmpTestStep := convertRtStep("xrayScan", "", "", node, ConvertXrayScanParamMapperList, xrayScanAttributesList)

	wantStep, err := ToStructFromJsonString[harness.Step](expectedXrayScanStep)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("Failed to convert JSON to struct: %v", diffs)
	}
}

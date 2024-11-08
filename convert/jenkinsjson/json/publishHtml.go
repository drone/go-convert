package json

import (
	"fmt"
	harness "github.com/drone/spec/dist/go"
)

func ConvertPublishHtml(node Node, variables map[string]string) *harness.Step {

	s, _ := ToJsonStringFromStruct[Node](node)

	fmt.Println(s)
	publishHtmlParameterMap, err := ToStructFromJsonString[PublishHtmlParameterMap](s)
	if err != nil {
		fmt.Println(err)
	}
	pmt := publishHtmlParameterMap.ParameterMap.Target

	fmt.Println(publishHtmlParameterMap)
	step := &harness.Step{
		Id:   SanitizeForId("UploadPublish", node.SpanId),
		Name: "Upload and Publish",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Connector: "<+input>",
			Image:     DroneS3UploadPublish,
			With: map[string]interface{}{
				"aws_access_key_id":     "<+input>",
				"aws_secret_access_key": "<+input>",
				"aws_bucket":            "<+input>",
				"default_region":        "<+input>",
				"source":                pmt.ReportDir,
				"target":                "<+pipeline.sequenceId>",
				"include":               pmt.Include,
				"artifact_file":         "artifact.txt",
			},
		},
	}
	return step
}

type PublishHtmlParameterMap struct {
	ParameterMap struct {
		Target struct {
			Include     string `json:"includes"`
			ReportDir   string `json:"reportDir"`
			ReportFiles string `json:"reportFiles"`
			ReportName  string `json:"reportName"`
		} `json:"target"`
	} `json:"ParameterMap"`
}

const DroneS3UploadPublish = "harnesscommunity/drone-s3-upload-publish"

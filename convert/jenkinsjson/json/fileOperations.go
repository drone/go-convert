package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

// createFileCopyStep creates a Harness step for file copy operations.
func ConvertFileCopy(node Node, operation map[string]interface{}) *harness.Step {
	args := operation["arguments"].(map[string]interface{})
	includes, _ := args["includes"].(string)
	targetLocation, _ := args["targetLocation"].(string)
	copyFileStep := &harness.Step{
		Name: operation["symbol"].(string),
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   fmt.Sprintf("cp -r %s %s", includes, targetLocation),
		},
	}
	return copyFileStep
}

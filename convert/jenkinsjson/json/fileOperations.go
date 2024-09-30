package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

// createFileDeleteStep creates a Harness step for file Delete operations.
func ConvertFileDelete(node Node, operation map[string]interface{}) *harness.Step {
	args := operation["arguments"].(map[string]interface{})
	includes, _ := args["includes"].(string)
	deleteFileStep := &harness.Step{
		Name: operation["symbol"].(string),
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   fmt.Sprintf("rm -rf %s", includes),
		},
	}
	return deleteFileStep
}

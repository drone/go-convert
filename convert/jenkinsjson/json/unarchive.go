package json

import (
	"fmt"
	"strings"

	harness "github.com/drone/spec/dist/go"
)

// ConvertUnarchive creates a Harness step for nunit plugin.
func ConvertUnarchive(node Node, paramMap map[string]interface{}) *harness.Step {

	// Extract the mapping from the parameterMap
	mapping := paramMap["mapping"].(map[string]interface{})

	// Initialize the script
	var scriptBuilder strings.Builder

	// Loop through each entry in the mapping
	for source, target := range mapping {
		// Determine the unarchiving command based on the file extension
		switch {
		case strings.HasSuffix(source, ".zip"):
			// Unzip command
			scriptBuilder.WriteString(fmt.Sprintf("mkdir -p %s && unzip -o %s -d %s\n", target, source, target))
		case strings.HasSuffix(source, ".tar"):
			// Tar extraction command
			scriptBuilder.WriteString(fmt.Sprintf("mkdir -p %s && tar -xvf %s -C %s\n", target, source, target))
		case strings.HasSuffix(source, ".tar.gz"):
			// Tar.gz extraction command
			scriptBuilder.WriteString(fmt.Sprintf("mkdir -p %s && tar -xzvf %s -C %s\n", target, source, target))
		default:
			// Unsupported format
			scriptBuilder.WriteString(fmt.Sprintf("echo 'Unsupported archive format for %s'\n", source))
		}
	}

	// Trim the trailing newline
	script := strings.TrimSuffix(scriptBuilder.String(), "\n")

	// Create the Harness step with the generated script
	convertUnArchive := &harness.Step{
		Name: "UnArchive",
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Shell: "sh",
			Run:   script,
		},
	}

	return convertUnArchive
}

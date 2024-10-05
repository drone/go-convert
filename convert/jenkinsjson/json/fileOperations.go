package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

// createFileCreateStep creates a Harness step for file Create operations.
func ConvertFileCreate(node Node, operation map[string]interface{}) *harness.Step {
	args := operation["arguments"].(map[string]interface{})
	fileName, _ := args["fileName"].(string)
	fileContent, _ := args["fileContent"].(string)
	createFileStep := &harness.Step{
		Name: operation["symbol"].(string),
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   fmt.Sprintf("echo '%s' > %s", fileContent, fileName),
		},
	}
	return createFileStep
}

// createFileDownloadStep creates a Harness step for file Download operations.
func ConvertFileDownload(node Node, operation map[string]interface{}) *harness.Step {
	args := operation["arguments"].(map[string]interface{})
	url, _ := args["url"].(string)
	targetLocation, _ := args["targetLocation"].(string)
	downloadFileStep := &harness.Step{
		Name: operation["symbol"].(string),
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Image: "alpine/curl",
			Run:   fmt.Sprintf("wget -P %s %s", targetLocation, url),
		},
	}
	return downloadFileStep
}

// createFileJoinStep creates a Harness step for file Join operations.
func ConvertFileJoin(node Node, operation map[string]interface{}) *harness.Step {
	args := operation["arguments"].(map[string]interface{})
	sourceFile, _ := args["sourceFile"].(string)
	targetFile, _ := args["targetFile"].(string)
	joinFileStep := &harness.Step{
		Name: operation["symbol"].(string),
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   fmt.Sprintf("cat %s >> %s", sourceFile, targetFile),
		},
	}
	return joinFileStep
}

// createFileJsonStep creates a Harness step for file file to json operations.
func ConvertFileJson(node Node, operation map[string]interface{}) *harness.Step {
	args := operation["arguments"].(map[string]interface{})
	sourceFile, _ := args["sourceFile"].(string)
	targetFile, _ := args["targetFile"].(string)
	jsonFileStep := &harness.Step{
		Name: operation["symbol"].(string),
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Image: "alpine",
			Run: fmt.Sprintf("stat -c '{\"size\": %%s, \"permissions\": \"%%A\", \"owner\": %%U, \"group\": %%G, \"last_modified\": \"%%y\"}' %s > %s",
				sourceFile, targetFile),
		},
	}
	return jsonFileStep
}

// createFileRenameStep creates a Harness step for file Rename operations.
func ConvertFileRename(node Node, operation map[string]interface{}) *harness.Step {
	args := operation["arguments"].(map[string]interface{})
	source, _ := args["source"].(string)
	destination, _ := args["destination"].(string)
	renameFileStep := &harness.Step{
		Name: operation["symbol"].(string),
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   fmt.Sprintf("mv %s %s", source, destination),
		},
	}
	return renameFileStep
}

// ConvertFileTranform creates a Harness step for file transform operations.
func ConvertFileTranform(node Node, operation map[string]interface{}) *harness.Step {
	args := operation["arguments"].(map[string]interface{})
	includes, _ := args["includes"].(string)
	excludes, _ := args["excludes"].(string)
	transformFileStep := &harness.Step{
		Name: operation["symbol"].(string),
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   fmt.Sprintf("find . -type f -name '%s' ! -name '%s' -exec sh -c 'iconv -f <source_encoding> -t UTF-8 \"$0\" -o \"${0%%.txt}.utf8\"' {} \\;", includes, excludes),
		},
	}
	return transformFileStep
}

// ConvertFolderCopy creates a Harness step for Folder Copy operations.
func ConvertFolderCopy(node Node, operation map[string]interface{}) *harness.Step {
	args := operation["arguments"].(map[string]interface{})
	sourceFolderPath, _ := args["sourceFolderPath"].(string)
	destinationFolderPath, _ := args["destinationFolderPath"].(string)
	folderCopyStep := &harness.Step{
		Name: operation["symbol"].(string),
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   fmt.Sprintf("cp -r %s %s", sourceFolderPath, destinationFolderPath),
		},
	}
	return folderCopyStep
}

// ConvertFolderCreate creates a Harness step for Folder Create operations.
func ConvertFolderCreate(node Node, operation map[string]interface{}) *harness.Step {
	args := operation["arguments"].(map[string]interface{})
	folderPath, _ := args["folderPath"].(string)
	folderCreateStep := &harness.Step{
		Name: operation["symbol"].(string),
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   fmt.Sprintf("mkdir -p %s", folderPath),
		},
	}
	return folderCreateStep
}

// ConvertFolderDelete creates a Harness step for folder delete operations.
func ConvertFolderDelete(node Node, operation map[string]interface{}) *harness.Step {
	args := operation["arguments"].(map[string]interface{})
	folderPath, _ := args["folderPath"].(string)
	folderDeleteStep := &harness.Step{
		Name: operation["symbol"].(string),
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   fmt.Sprintf("rm -rf %s", folderPath),
		},
	}
	return folderDeleteStep
}

// ConvertFolderRename creates a Harness step for folder rename operations.
func ConvertFolderRename(node Node, operation map[string]interface{}) *harness.Step {
	args := operation["arguments"].(map[string]interface{})
	source, _ := args["source"].(string)
	destination, _ := args["destination"].(string)
	folderRenameStep := &harness.Step{
		Name: operation["symbol"].(string),
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   fmt.Sprintf("mv %s %s", source, destination),
		},
	}
	return folderRenameStep
}

// ConvertFileUntar creates a Harness step for file untar operations.
func ConvertFileUntar(node Node, operation map[string]interface{}) *harness.Step {
	var format string
	var source string
	var target string
	args := operation["arguments"].(map[string]interface{})

	if filePath, ok := args["filePath"].(string); ok {
		source = filePath
	}
	if targetLocation, ok := args["targetLocation"].(string); ok {
		target = targetLocation
	}
	if isGZIP, ok := args["isGZIP"].(bool); ok {
		if isGZIP {
			format = "gzip"
		} else {
			format = "tar"
		}
	}
	fileUntarStep := &harness.Step{
		Name: operation["symbol"].(string),
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image: "plugins/archive:latest",
			With: map[string]interface{}{
				"action": "extract",
				"format": format,
				"source": source,
				"target": target,
			},
		},
	}
	return fileUntarStep
}

// ConvertFileUnzip creates a Harness step for file unzip operations.
func ConvertFileUnzip(node Node, operation map[string]interface{}) *harness.Step {
	var source string
	var target string
	args := operation["arguments"].(map[string]interface{})

	if filePath, ok := args["filePath"].(string); ok {
		source = filePath
	}
	if targetLocation, ok := args["targetLocation"].(string); ok {
		target = targetLocation
	}

	fileUnzipStep := &harness.Step{
		Name: operation["symbol"].(string),
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image: "plugins/archive:latest",
			With: map[string]interface{}{
				"action": "extract",
				"format": "zip",
				"source": source,
				"target": target,
			},
		},
	}
	return fileUnzipStep
}

// ConvertFileZip creates a Harness step for file zip operations.
func ConvertFileZip(node Node, operation map[string]interface{}) *harness.Step {
	var source string
	var target string
	args := operation["arguments"].(map[string]interface{})

	if folderPath, ok := args["folderPath"].(string); ok {
		source = folderPath
	}
	if outputFolderPath, ok := args["outputFolderPath"].(string); ok {
		target = outputFolderPath
	}

	fileZipStep := &harness.Step{
		Name: operation["symbol"].(string),
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image: "plugins/archive:latest",
			With: map[string]interface{}{
				"action":    "archive",
				"format":    "zip",
				"overwrite": "true",
				"source":    source,
				"target":    target,
			},
		},
	}
	return fileZipStep
}

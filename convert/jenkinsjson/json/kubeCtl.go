package json

import (
	harness "github.com/drone/spec/dist/go"
)

// ConvertNunit creates a Harness step for nunit plugin.
func ConvertKubeCtl(node Node, paramMap map[string]interface{}) *harness.Step {
	var contextName string = ""
	var serverUrl string = ""
	var clusterName string = ""
	var namespace string = ""
	var credentialsId string = ""
	var caCertificate string = ""

	script := findScript(node)

	if contextNameValue, ok := node.ParameterMap["contextName"]; ok {
		contextName = contextNameValue.(string)
	}

	if serverUrlValue, ok := node.ParameterMap["serverUrl"]; ok {
		serverUrl = serverUrlValue.(string)
	}

	if clusterNameValue, ok := node.ParameterMap["clusterName"]; ok {
		clusterName = clusterNameValue.(string)
	}

	if namespaceValue, ok := node.ParameterMap["namespace"]; ok {
		namespace = namespaceValue.(string)
	}

	if credentialsIdValue, ok := node.ParameterMap["credentialsId"]; ok {
		credentialsId = credentialsIdValue.(string)
	}

	if caCertificateValue, ok := node.ParameterMap["caCertificate"]; ok {
		caCertificate = caCertificateValue.(string)
	}

	convertNunit := &harness.Step{
		Name: "Apply Kubernetes Files",
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Connector: "docker_hub_connector",
			Image:     "bitnami/kubectl",
			Shell:     "sh",
			Envs: map[string]string{
				"KUBE_SERVER":         serverUrl,
				"KUBE_CONTEXT":        contextName,
				"KUBE_CLUSTER":        clusterName,
				"KUBE_TOKEN":          "<+secrets.getValue(\"kube_token\")>",
				"KUBE_NAMESPACE":      namespace,
				"KUBE_CA_CERTIFICATE": caCertificate,
				"KUBE_CREDENTIALS_ID": credentialsId,
			},
			Run: script,
		},
	}

	return convertNunit
}

// Function to recursively search for "script" in a specific child node
func findScript(node Node) string {
	// Check if "jenkins.pipeline.step.type" is "withKubeConfig"
	if node.AttributesMap["jenkins.pipeline.step.type"] == "withKubeConfig" {
		// Traverse through children
		for _, child := range node.Children {
			// Check if "jenkins.pipeline.step.type" in the child is "sh"
			if child.AttributesMap["jenkins.pipeline.step.type"] == "sh" {
				// Check if "script" exists in ParameterMap and return it
				if scriptValue, ok := child.ParameterMap["script"]; ok {
					return scriptValue.(string)
				}
			} else {
				// Recursively check in case the child has further descendants
				if result := findScript(child); result != "" {
					return result
				}
			}
		}
	}
	// Return empty string if no valid "script" found
	return ""
}

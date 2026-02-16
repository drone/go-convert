package converthelpers

import (
	"bytes"
	"encoding/json"
	"fmt"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// HTTPStepWith represents the 'with' configuration for httpStep@1.0.0 template
type HTTPStepWith struct {
	URL             string            `json:"url,omitempty"`
	Method          string            `json:"method,omitempty"`
	Headers         string            `json:"headers,omitempty"`
	Body            string            `json:"body,omitempty"`
	DisableRedirect bool              `json:"disable_redirect,omitempty"`
	Assertion       string            `json:"assertion,omitempty"`
	OutputVars      string            `json:"output_vars,omitempty"`
	EnvVars         map[string]string `json:"env_vars,omitempty"`
	ClientCert      string            `json:"client_cert,omitempty"`
	ClientKey       string            `json:"client_key,omitempty"`
}

// ConvertStepHTTP converts a v0 HTTP step to v1 template spec
func ConvertStepHTTP(src *v0.Step) *v1.StepRun {
	if src == nil || src.Spec == nil {
		return nil
	}

	// Extract the typed spec
	spec, ok := src.Spec.(*v0.StepHTTP)
	if !ok {
		return nil
	}

	container := &v1.Container{
		Image:     "harnessdev/harness-http:harness-http-v0.0.1",
		Connector: "account.harnessImage",
	}

	env := map[string]interface{}{
		"PLUGIN_URL":              spec.URL,
		"PLUGIN_METHOD":           spec.Method,
		"PLUGIN_DISABLE_REDIRECT": false,
	}
	if spec.Assertion != "" {
		env["PLUGIN_ASSERTION"] = spec.Assertion
	}

	// with := HTTPStepWith{
	// 	URL:             spec.URL,
	// 	Method:          spec.Method,
	// 	DisableRedirect: false,
	// 	Assertion:       spec.Assertion,
	// }

	// Convert headers from []*Variable to JSON string format: "{h1: v1, h2: v2}"
	if len(spec.Headers) > 0 {
		headersMap := make(map[string]interface{})
		for _, h := range spec.Headers {
			if h != nil && h.Key != "" {
				// Handle both string and non-string values
				if h.Value != "" {
					headersMap[h.Key] = h.Value
				}
			}
		}
		if len(headersMap) > 0 {
			// Use custom encoder to prevent HTML escaping of angle brackets
			var buf bytes.Buffer
			encoder := json.NewEncoder(&buf)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(headersMap); err == nil {
				// Remove trailing newline added by encoder
				env["PLUGIN_HEADERS"] = string(bytes.TrimSpace(buf.Bytes()))
			}
		}
	}

	// Convert request body - wrap in single quotes if present
	if spec.RequestBody != "" {
		env["PLUGIN_BODY"] = spec.RequestBody
	}

	// Convert output variables from []*Variable to JSON string format: "{o1: v1, o2: v2}"
	if len(spec.OutputVariables) > 0 {
		outputMap := make(map[string]interface{})
		for _, ov := range spec.OutputVariables {
			if ov != nil && ov.Name != "" {
				// Handle both string and non-string values
				if ov.Value != nil {
					outputMap[ov.Name] = ov.Value
				}
			}
		}
		if len(outputMap) > 0 {
			// Use custom encoder to prevent HTML escaping of angle brackets
			var buf bytes.Buffer
			encoder := json.NewEncoder(&buf)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(outputMap); err == nil {
				// Remove trailing newline added by encoder
				env["PLUGIN_OUTPUT_VARS"] = string(bytes.TrimSpace(buf.Bytes()))
			}
		}
	}

	// Convert input variables to env_vars map
	if len(spec.InputVariables) > 0 {
		envVars := make(map[string]string)
		for _, iv := range spec.InputVariables {
			if iv != nil && iv.Name != "" {
				// Convert value to string
				if iv.Value != nil {
					envVars[iv.Name] = fmt.Sprintf("%v", iv.Value)
				}
			}
		}
		if len(envVars) > 0 {
			env["PLUGIN_ENV_VARS"] = envVars
		}
	}

	// Handle certificate fields
	if spec.Certificate != "" {
		env["PLUGIN_CLIENT_CERT"] = spec.Certificate
	}
	if spec.CertificateKey != "" {
		env["PLUGIN_CLIENT_KEY"] = spec.CertificateKey
	}

	return &v1.StepRun{
		Container: container,
		Env:       &flexible.Field[map[string]interface{}]{Value: env},
	}
}

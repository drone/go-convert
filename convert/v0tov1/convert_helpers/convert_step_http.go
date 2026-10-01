package converthelpers

import (
	"bytes"
	"encoding/json"
	"fmt"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// encodeJSONNoEscape marshals v to a JSON string without HTML-escaping
// characters such as '<', '>' and '&' (e.g. in CEL assertions).
func encodeJSONNoEscape(v interface{}) (string, bool) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return "", false
	}
	// Remove trailing newline added by encoder
	return string(bytes.TrimSpace(buf.Bytes())), true
}

// ConvertStepHTTP converts a v0 HTTP step to the v1 httpStep template
func ConvertStepHTTP(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}

	// Extract the typed spec
	spec, ok := src.Spec.(*v0.StepHTTP)
	if !ok {
		return nil
	}

	with := make(map[string]interface{})

	if spec.URL != "" {
		with["url"] = spec.URL
	}
	if spec.Method != "" {
		with["method"] = spec.Method
	}

	// Convert headers from []*KeyValuePair to a JSON object string, e.g. {"h1":"v1","h2":"v2"}
	if len(spec.Headers) > 0 {
		headersMap := make(map[string]interface{})
		for _, h := range spec.Headers {
			if h != nil && h.Key != "" && h.Value != "" {
				headersMap[h.Key] = h.Value
			}
		}
		if len(headersMap) > 0 {
			if headers, ok := encodeJSONNoEscape(headersMap); ok {
				with["headers"] = headers
			}
		}
	}

	if spec.RequestBody != "" {
		with["body"] = spec.RequestBody
	}

	if spec.Assertion != "" {
		with["assertion"] = spec.Assertion
	}

	// Convert output variables from []*Variable to a JSON object string, e.g. {"o1":"response.body.id"}
	if len(spec.OutputVariables) > 0 {
		outputMap := make(map[string]interface{})
		for _, ov := range spec.OutputVariables {
			if ov != nil && ov.Name != "" && ov.Value != nil {
				outputMap[ov.Name] = ov.Value
			}
		}
		if len(outputMap) > 0 {
			if outputVars, ok := encodeJSONNoEscape(outputMap); ok {
				with["output_vars"] = outputVars
			}
		}
	}

	// Convert input variables to a JSON object string for env_vars
	if len(spec.InputVariables) > 0 {
		envVars := make(map[string]interface{})
		for _, iv := range spec.InputVariables {
			if iv != nil && iv.Name != "" && iv.Value != nil {
				envVars[iv.Name] = fmt.Sprintf("%v", iv.Value)
			}
		}
		if len(envVars) > 0 {
			if env, ok := encodeJSONNoEscape(envVars); ok {
				with["env_vars"] = env
			}
		}
	}

	if spec.Certificate != "" {
		with["client_cert"] = spec.Certificate
	}
	if spec.CertificateKey != "" {
		with["client_key"] = spec.CertificateKey
	}

	// body_file, work_dir, skip_verify, proxy, disable_redirect, log_level, output_file:
	// no corresponding v0 StepHTTP field; template defaults apply.

	return &v1.StepTemplate{
		Uses: "httpStep",
		With: with,
	}
}

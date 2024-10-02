package json

import (
	harness "github.com/drone/spec/dist/go"
	"log"
	"strings"
)

const (
	HttpRequestPluginImage = "plugins/httpRequest"
)

var JenkinsToDroneParamMapperList = []JenkinsToDroneParamMapper{
	{"url", "url", StringType, nil},
	{"httpMode", "http_method", StringType, nil},
	{"validResponseCodes", "valid_response_codes", StringType, nil},
	{"timeout", "timeout", Float64Type, nil},
	{"validResponseContent", "valid_response_body", StringType, nil},
	{"wrapAsMultipart", "wrap_as_multipart", BoolType, nil},
	{"outputFile", "output_file", StringType, nil},
	{"ignoreSslErrors", "ignore_ssl", BoolType, nil},
	{"requestBody", "request_body", StringType, nil},
	{"consoleLogResponseBody", "log_response", BoolType, nil},
	{"quiet", "quiet", BoolType, nil},
	{"contentType", "content_type", StringType, nil},
	{"customHeaders", "headers", DontCare, HeaderToStrCsv}, // HeaderToStrCsv overrides the default behavior
	{"authentication", "auth_basic", StringType, nil},
	{"acceptType", "accept_type", StringType, nil},
	{"uploadFile", "upload_file", StringType, nil},
	{"multipartName", "multipart_name", StringType, nil},
}

func ConvertHttpRequest(node Node, variables map[string]string) *harness.Step {
	step := ConvertToStepWithProperties(&node, variables, JenkinsToDroneParamMapperList,
		HttpRequestPluginImage)

	return step
}

func HeaderToStrCsv(node *Node, attrMap map[string]interface{}, jenkinsKey string) (interface{}, error) {

	var ret interface{}
	var ls []string

	valList, ok := attrMap[jenkinsKey].([]interface{})
	if !ok {
		log.Printf("jenkins param %s missing for node %s", jenkinsKey, node.SpanId)
		return "", nil
	}

	for _, v := range valList {
		vm := v.(map[string]interface{})

		name, ok := vm["name"].(string)
		if !ok {
			log.Printf("jenkins param %s missing for node %s", jenkinsKey, node.SpanId)
			continue
		}

		value, ok := vm["value"].(string)
		if !ok {
			log.Printf("jenkins param %s missing for node %s", jenkinsKey, node.SpanId)
			continue
		}

		kvPair := name + ":" + value
		ls = append(ls, kvPair)

	}
	ret = strings.Join(ls, ",")

	return ret, nil
}

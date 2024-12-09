package json

import (
	"encoding/json"
	"errors"
	harness "github.com/drone/spec/dist/go"
	"log"
	"strconv"
	"strings"
)

type ParamTransform func(node *Node, attrMap map[string]interface{}, jenkinsKey string) (interface{}, error)

type JenkinsToDroneParamMapper struct {
	JenkinsParam     string
	DroneParam       string
	JenkinsParamType string
	TransformFunc    ParamTransform
}

func ConvertToStepWithProperties(node *Node, variables map[string]string,
	tmpJenkinsToDroneParamMapperList []JenkinsToDroneParamMapper, imageName string) *harness.Step {

	step := GetStepWithProperties(node, tmpJenkinsToDroneParamMapperList, imageName)

	if len(variables) > 0 {
		step.Spec.(*harness.StepPlugin).Envs = variables
	}

	return step
}

func ToJsonStringFromMap[T any](m T) (string, error) {
	outBytes, err := json.Marshal(m)
	if err == nil {
		return string(outBytes), nil
	}
	return "", err
}

func GetStepWithProperties(node *Node,
	tmpJenkinsToDroneParamMapperList []JenkinsToDroneParamMapper, imageName string) *harness.Step {

	withProperties := map[string]interface{}{}

	attr, ok := node.AttributesMap[HarnessAttribute]
	if !ok {
		log.Printf("harness-attribute missing for spanName %s", node.SpanName)
		return nil
	}

	attrMap, err := ToMapFromJsonString[map[string]interface{}](attr)
	if err != nil {
		log.Printf("failed to unmarshal harness-attribute for node %s: %v", node.SpanName, err)
		return nil
	}

	for _, val := range tmpJenkinsToDroneParamMapperList {
		SafeAssignWithPropertiesTyped(node, &withProperties, attrMap, val.JenkinsParam,
			val.DroneParam, val.JenkinsParamType, val.TransformFunc, false)
	}

	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image:  imageName,
			Inputs: withProperties,
			With:   withProperties,
		},
	}

	return step
}

func ConvertToStepUsingParameterMapDelegate(node *Node, variables map[string]string,
	tmpJenkinsToDroneParamMapperList []JenkinsToDroneParamMapper, imageName string) *harness.Step {

	step := GetStepUsingParameterMapDelegate(node, tmpJenkinsToDroneParamMapperList, imageName)

	if len(variables) > 0 {
		step.Spec.(*harness.StepPlugin).Envs = variables
	}

	return step
}

func GetStepUsingParameterMapDelegate(node *Node,
	tmpJenkinsToDroneParamMapperList []JenkinsToDroneParamMapper, imageName string) *harness.Step {

	withProperties := map[string]interface{}{}

	delegateMapIfce, ok := node.ParameterMap["delegate"]
	if !ok {
		return &harness.Step{}
	}

	delegateMap, err := delegateMapIfce.(map[string]interface{})
	if !err {
		return &harness.Step{}
	}

	arguments, ok := delegateMap["arguments"]
	if !ok {
		return &harness.Step{}
	}

	argumentsMap, err := arguments.(map[string]interface{})
	if !err {
		return &harness.Step{}
	}

	for _, val := range tmpJenkinsToDroneParamMapperList {
		SafeAssignWithPropertiesTyped(node, &withProperties, argumentsMap, val.JenkinsParam,
			val.DroneParam, val.JenkinsParamType, val.TransformFunc, false)
	}

	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image:  imageName,
			Inputs: withProperties,
			With:   withProperties,
		},
	}

	return step
}

func SafeAssignWithPropertiesTyped(node *Node, withProperties *map[string]interface{},
	attrMap map[string]interface{}, jenkinsKey, droneKey, jenkinsParamType string,
	paramTransformFunc ParamTransform, isWarn bool) {

	var valOk bool
	var newVal interface{}

	if paramTransformFunc != nil { // paramTransformFunc overrides the default behavior
		retVal, err := paramTransformFunc(node, attrMap, jenkinsKey)
		if err != nil {
			return
		}
		(*withProperties)[droneKey] = retVal
		return
	}

	val, found := attrMap[jenkinsKey]
	if !found {
		if isWarn {
			log.Printf("jenkins param -- %s missing for node %s", jenkinsKey, droneKey)
		}
		return
	}

	switch jenkinsParamType {
	case StringType:
		newVal, valOk = CastTo[string](val)
	case Float64Type:
		newVal, valOk = CastTo[float64](val)
	case BoolType:
		newVal, valOk = CastTo[bool](val)
	case InterfaceListType:
		newVal, valOk = CastTo[[]interface{}](val)
	}

	if !valOk {
		//log.Printf("jenkins param %s is not a %s for node %s", droneKey, jenkinsParamType, droneKey)
		return
	}

	(*withProperties)[droneKey] = newVal
}

func CastTo[T any](jenkinsParamInterface interface{}) (T, bool) {
	tmpParamVal, ok := jenkinsParamInterface.(T)
	return tmpParamVal, ok
}

func ToMapFromJsonString[T any](jsonString string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(jsonString), &result)
	return result, err
}

func ToJsonStringFromStruct[T any](v T) (string, error) {
	jsonBytes, err := json.Marshal(v)

	if err == nil {
		return string(jsonBytes), nil
	}

	return "", err
}

func ToStructFromJsonString[T any](jsonStr string) (T, error) {
	var v T
	err := json.Unmarshal([]byte(jsonStr), &v)
	return v, err
}

func ToStringArrayFromCsvString(csv string) ([]string, error) {
	if csv == "" {
		return nil, errors.New("input string is empty")
	}

	parts := strings.Split(csv, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
		if parts[i] == "" {
			return nil, errors.New("invalid CSV string: contains empty values")
		}
	}

	return parts, nil
}

func ToFloat64FromString(s string) (float64, error) {
	valFloat64, err := strconv.ParseFloat(s, 64)
	return valFloat64, err
}

const (
	Float64Type       = "float64"
	StringType        = "string"
	BoolType          = "bool"
	InterfaceListType = "InterfaceList"
	DontCare          = "DontCare"
	HarnessAttribute  = "harness-attribute"
)

//
//

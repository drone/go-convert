package json

import (
	"encoding/json"
	harness "github.com/drone/spec/dist/go"
	"log"
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
			val.DroneParam, val.JenkinsParamType, val.TransformFunc)
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
	attrMap map[string]interface{}, jenkinsKey, droneKey, jenkinsParamType string, paramTransformFunc ParamTransform) {

	var valOk bool
	var newVal interface{}

	if paramTransformFunc != nil { // paramTransformFunc overrides the default behavior
		retVal, err := paramTransformFunc(node, attrMap, jenkinsKey)
		if err != nil {
			// log.Printf("jenkins parameter %s is not a string for node %s", droneKey, droneKey)
			return
		}
		(*withProperties)[droneKey] = retVal
		return
	}

	val, found := attrMap[jenkinsKey]
	if !found {
		log.Printf("jenkins param -- %s missing for node %s", jenkinsKey, droneKey)
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

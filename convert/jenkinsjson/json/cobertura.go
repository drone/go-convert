package json

import (
	"fmt"
	harness "github.com/drone/spec/dist/go"
)

var CorberturaJenkinsToDroneParamMapperList = []JenkinsToDroneParamMapper{
	{"coberturaReportFile", "reports_path_pattern", StringType, nil},
	{"failUnstable", "fail_on_threshold", BoolType, nil},
	{"failUnhealthy", "fail_on_threshold", BoolType, nil},
	{"conditionalCoverageTargets", "threshold_branch", Float64Type, MakeThresholdGetter("conditionalCoverageTargets")},
	{"classCoverageTargets", "threshold_class", Float64Type, MakeThresholdGetter("classCoverageTargets")},
	{"fileCoverageTargets", "threshold_file", Float64Type, MakeThresholdGetter("fileCoverageTargets")},
	{"packageCoverageTargets", "threshold_package", Float64Type, MakeThresholdGetter("packageCoverageTargets")},
	{"lineCoverageTargets", "threshold_line", Float64Type, MakeThresholdGetter("lineCoverageTargets")},
	{"methodCoverageTargets", "threshold_method", Float64Type, MakeThresholdGetter("methodCoverageTargets")},
	{"tool", "tool", StringType, SetCoberturaTool},
	// runAlways - Missing convert logic: When: parametersMap.delegate.arguments.changeBuildStatus
	// {"runAlways", "run_always", BoolType, nil},
}

func MakeThresholdGetter(attrName string) func(node *Node,
	attrMap map[string]interface{}, jenkinsKey string) (interface{}, error) {
	return func(node *Node, attrMap map[string]interface{}, jenkinsKey string) (interface{}, error) {
		classThreshold, _, err := GetCoberturaThreshold(attrName, attrMap)
		if err != nil {
			return 0.0, err
		}
		return classThreshold, nil
	}
}

func GetCoberturaThreshold(attrKey string, attrMap map[string]interface{}) (interface{}, []string, error) {

	var valuesList []string
	var err error

	csvParams, ok := attrMap[attrKey]
	if !ok {
		return "", valuesList, fmt.Errorf("Error in GetCoberturaThreshold: %s not found", attrKey)
	}
	valuesList, err = ToStringArrayFromCsvString(csvParams.(string))
	if err != nil {
		return "", valuesList, err
	}

	if len(valuesList) != 3 {
		return "", valuesList, fmt.Errorf("Error in GetCoberturaThreshold: %s must have 3 values", attrKey)
	}

	retFloat64Val, err := ToFloat64FromString(valuesList[2])
	if err != nil {
		return "", valuesList, err
	}

	return retFloat64Val, valuesList, nil
}

func SetCoberturaTool(node *Node, attrMap map[string]interface{}, jenkinsKey string) (interface{}, error) {
	return CoberturaToolName, nil
}

func ConvertCobertura(node Node, variables map[string]string) *harness.Step {

	step := ConvertToStepUsingParameterMapDelegate(&node, variables, CorberturaJenkinsToDroneParamMapperList,
		CoverageReportImage)

	return step
}

const (
	CoberturaToolName = "cobertura"
)

//

package json

import (
	"fmt"
	harness "github.com/drone/spec/dist/go"
)

var QualityGatesMap = map[string]interface{}{
	"MODULE":             "threshold_module",
	"CLASS":              "threshold_class",
	"FILE":               "threshold_file",
	"PACKAGE":            "threshold_package",
	"LINE":               "threshold_line",
	"METHOD":             "threshold_method",
	"INSTRUCTION":        "threshold_instruction",
	"BRANCH":             "threshold_branch",
	"COMPLEXITY":         "threshold_complexity",
	"COMPLEXITY_DENSITY": "threshold_complexity_density",
	"LOC":                "threshold_loc",
}

func ConvertRecordCoverage(node Node, variables map[string]string) []*harness.Step {

	s, err := ToJsonStringFromStruct(node)
	if err != nil {
		fmt.Println("Error converting json to struct:", err)
		return nil
	}

	rcn, err := ToStructFromJsonString[RecordCoverageNode](s)
	if err != nil {
		fmt.Println("Error converting json to struct:", err)
		return nil
	}

	var stepsMap []map[string]interface{}

	for _, cvg := range rcn.ParameterMap.Tools {
		stepMap := map[string]interface{}{}
		switch cvg.Parser {
		case "JACOCO":
			stepMap["tool"] = "jacoco-xml"
		case "COBERTURA":
			stepMap["tool"] = "cobertura"
		default:
			fmt.Println("Unknown coverage parser:", cvg.Parser)
		}
		if cvg.Pattern != "" {
			stepMap["reports_path_pattern"] = cvg.Pattern
		}
		stepMap["source_code_encoding"] = "UTF-8"
		stepMap["fail_on_threshold"] = rcn.ParameterMap.EnabledForFailure

		for _, qualityGate := range rcn.ParameterMap.QualityGates {
			if val, ok := QualityGatesMap[qualityGate.Metric]; ok {
				stepMap[val.(string)] = qualityGate.Threshold
			} else {
				fmt.Println("Unknown metric:", qualityGate.Metric)
			}
		}

		stepsMap = append(stepsMap, stepMap)
	}

	var stepsList []*harness.Step
	for _, stepMap := range stepsMap {
		withProperties := map[string]interface{}{}
		step := &harness.Step{
			Name: node.SpanName,
			Id:   SanitizeForId(node.SpanName, node.SpanId),
			Type: "plugin",
			Spec: &harness.StepPlugin{
				Image:  "plugins/coverage-report",
				Inputs: withProperties,
				With:   withProperties,
			},
		}
		step.Spec.(*harness.StepPlugin).With = stepMap
		stepsList = append(stepsList, step)
	}

	return stepsList
}

type RecordCoverageNode struct {
	ParameterMap struct {
		SourceCodeRetention string `json:"sourceCodeRetention"`
		Name                string `json:"name"`
		QualityGates        []struct {
			Unstable  bool    `json:"unstable"`
			Metric    string  `json:"metric"`
			Threshold float64 `json:"threshold"`
			Baseline  string  `json:"baseline"`
		} `json:"qualityGates"`
		EnabledForFailure bool   `json:"enabledForFailure"`
		Id                string `json:"id"`
		Tools             []struct {
			Parser  string `json:"parser"`
			Pattern string `json:"pattern,omitempty"`
		} `json:"tools"`
	} `json:"parameterMap"`
}

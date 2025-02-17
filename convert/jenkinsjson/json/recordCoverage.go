package json

import (
	"fmt"
	harness "github.com/drone/spec/dist/go"
)

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
			if qualityGate.Metric == "MODULE" {
				stepMap["threshold_module"] = qualityGate.Threshold
			} else if qualityGate.Metric == "CLASS" {
				stepMap["threshold_class"] = qualityGate.Threshold
			} else if qualityGate.Metric == "FILE" {
				stepMap["threshold_file"] = qualityGate.Threshold
			} else if qualityGate.Metric == "PACKAGE" {
				stepMap["threshold_package"] = qualityGate.Threshold
			} else if qualityGate.Metric == "LINE" {
				stepMap["threshold_line"] = qualityGate.Threshold
			} else if qualityGate.Metric == "METHOD" {
				stepMap["threshold_method"] = qualityGate.Threshold
			} else if qualityGate.Metric == "INSTRUCTION" {
				stepMap["threshold_instruction"] = qualityGate.Threshold
			} else if qualityGate.Metric == "BRANCH" {
				stepMap["threshold_branch"] = qualityGate.Threshold
			} else if qualityGate.Metric == "COMPLEXITY" {
				stepMap["threshold_complexity"] = qualityGate.Threshold
			} else if qualityGate.Metric == "COMPLEXITY_DENSITY" {
				stepMap["threshold_complexity_density"] = qualityGate.Threshold
			} else if qualityGate.Metric == "LOC" {
				stepMap["threshold_loc"] = qualityGate.Threshold
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

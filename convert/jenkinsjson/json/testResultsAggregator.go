package json

import (
	harness "github.com/drone/spec/dist/go"
)

const (
	testResultsAggregatorImage = "plugins/test-results-aggregator"
)

var ConvertTestResultsAggregatorParamMapperList = []JenkinsToDroneParamMapper{
	{"compareWithPreviousRun", "compare_build_results", BoolType, nil},
	{"influxdbBucket", "influxdb_bucket", StringType, nil},
	{"influxdbOrg", "influxdb_org", StringType, nil},
	{"influxdbToken", "influxdb_token", StringType, nil},
	{"influxdbUrl", "influxdb_url", StringType, nil},
}

func ConvertTestResultsAggregator(node Node, variables map[string]string) *harness.Step {
	step := ConvertToStepWithProperties(&node, variables, ConvertTestResultsAggregatorParamMapperList,
		testResultsAggregatorImage)
	tmpStepPlugin := step.Spec.(*harness.StepPlugin)
	tmpStepPlugin.With["tool"] = "<+input>"
	tmpStepPlugin.With["group"] = "<+input>"
	tmpStepPlugin.With["include_pattern"] = "<+input>"
	tmpStepPlugin.With["reports_dir"] = "<+input>"
	return step
}

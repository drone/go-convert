package convertexpressions

var FailureStrategiesConversionRules = []ConversionRule{
	{"failureStrategies[i].onFailure.errors", "onFailure[i].-.errors"},
	// {"failureStrategies[i].onFailure.action.type", "onFailure[i].action.ACTION_TYPE.type"}, // not directly
	{"failureStrategies[i].onFailure.action.specConfig.retryCount", "onFailure[i].-.action.-.retry.attempts"},
	{"failureStrategies[i].onFailure.action.specConfig.retryInterval", "onFailure[i].-.action.-.retry.interval"},
	{"failureStrategies[i].onFailure.action.specConfig.onRetryFailure.action.type", "onFailure[i].-.action.-.retry.failureAction.type"}, // map back to action node

	{"failureStrategies[i].onFailure.action.specConfig.timeout", "onFailure[i].-.action.-.manualIntervention.timeout"},
	{"failureStrategies[i].onFailure.action.specConfig.onTimeout.action.type", "onFailure[i].-.action.-.manualIntervention.timeoutAction.type"}, //map back to action node
}
// <+pipeline.stages.cd.spec.execution.steps.sh.failureStrategies[0].onFailure.action.specConfig.onRetryFailure.action.type>
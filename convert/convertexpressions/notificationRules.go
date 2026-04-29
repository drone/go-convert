package convertexpressions

var NotificationRulesConversionRules = []ConversionRule{
	{"notificationRules[i].identifier", "notification[i].id"},
	{"notificationRules[i].notificationMethod.type", "notification[i].uses"},
	// common
	{"notificationRules[i].notificationMethod.spec.delegateSelectors", "notification[i].with.delegate-selectors"},
	{"notificationRules[i].notificationMethod.spec.executeOnDelegate", "notification[i].with.execute-on-delegate"},
	
	{"notificationRules[i].notificationMethod.spec.userGroup", "notification[i].with.user-group"}, // slack,email,pagerduty,ms-teams
	{"notificationRules[i].notificationMethod.spec.webhookUrl", "notification[i].with.webhook"}, //slack
	{"notificationRules[i].notificationMethod.spec.recipients", "notification[i].with.recipients"}, // email
	{"notificationRules[i].notificationMethod.spec.webhookUrl", "notification[i].with.url"}, // webhook 
	{"notificationRules[i].notificationMethod.spec.headers", "notification[i].with.headers"}, // webhook,datadog
	{"notificationRules[i].notificationMethod.spec.integrationKey", "notification[i].with.integration-key"}, //pagerduty
	{"notificationRules[i].notificationMethod.spec.msTeamKeys", "notification[i].with.msteam-keys"}, // ms-teams
	{"notificationRules[i].notificationMethod.spec.apiKey", "notification[i].with.api-key"}, // datadog
	{"notificationRules[i].notificationMethod.spec.url", "notification[i].with.url"}, // datadog

}
package json

import (
	"fmt"
	"log"
)

func ConvertWithCredentials(node Node) map[string]string {
	envVars := make(map[string]string)

	if bindings, ok := node.ParameterMap["bindings"].([]interface{}); ok {
		for _, binding := range bindings {
			bindingMap, okBind := binding.(map[string]interface{})
			symbol, okSym := bindingMap["symbol"].(string)
			arguments, okArgs := bindingMap["arguments"].(map[string]interface{})
			if !okBind || !okSym || !okArgs {
				continue
			}

			switch symbol {
			case "usernamePassword":
				addVariable("usernameVariable", arguments, envVars)
				addVariable("passwordVariable", arguments, envVars)
			case "string":
				addVariable("variable", arguments, envVars)
			default:
				log.Printf("withCredentials unsupported symbol: %s", symbol)
			}
		}
	}

	return envVars
}

func addVariable(key string, arguments map[string]interface{}, envVars map[string]string) {
	if passwordVariable, ok := arguments[key].(string); ok {
		envVars[passwordVariable] = fmt.Sprintf("<+pipeline.variables.%s>", passwordVariable)
	}
}

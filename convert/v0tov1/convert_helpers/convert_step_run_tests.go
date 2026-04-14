package converthelpers

import (
	"strings"
	"fmt"
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// ConvertStepRunTests converts a v0 RunTests step to v1 run-test step
func ConvertStepRunTests(src *v0.Step) *v1.StepTest {
	if src == nil || src.Spec == nil {
		return nil
	}
	sp, ok := src.Spec.(*v0.StepRunTests)
	if !ok {
		return nil
	}

	// Build the script from language, buildTool, preCommand, args, postCommand
	script := buildTestScript(sp.Language, sp.BuildTool, sp.PreCommand, sp.Args, sp.PostCommand)

	// Container mapping
	var container *v1.Container
	if sp.Image != "" || sp.ConnectorRef != "" || sp.Privileged != nil || (sp.Resources != nil && sp.Resources.Limits != nil) || sp.RunAsUser != nil {
		pull := convertImagePullPolicy(sp.ImagePullPolicy)
		cpu := ""
		memory := ""
		if sp.Resources != nil && sp.Resources.Limits != nil {
			cpu = sp.Resources.Limits.GetCPUString()
			memory = sp.Resources.Limits.GetMemoryString()
		}
		container = &v1.Container{
			Image:      sp.Image,
			Connector:  sp.ConnectorRef,
			Privileged: sp.Privileged,
			Pull:       pull,
			Cpu:        cpu,
			Memory:     memory,
			User:       sp.RunAsUser,
		}
	}

	// Reports mapping
	var report *v1.Reports
	if sp.Reports != nil {
		report = &v1.Reports{}
		report.Type = strings.ToLower(sp.Reports.Type)
		if sp.Reports.Spec != nil {
			report.Paths = sp.Reports.Spec.Paths
		}
	}

	// Shell mapping - lower-case common values
	shell := strings.ToLower(sp.Shell)

	var intelligence *v1.TestIntelligence
	if sp.RunOnlySelectedTests != nil {
		intelligence = &v1.TestIntelligence{
			Disabled: flexible.NegateBool(sp.RunOnlySelectedTests),
		}
	}

	// Match globs (testGlobs -> match)
	var match v1.Stringorslice
	if sp.TestGlobs != "" {
		globs := strings.Split(sp.TestGlobs, ",")
		for _, g := range globs {
			trimmed := strings.TrimSpace(g)
			if trimmed != "" {
				match = append(match, trimmed)
			}
		}
	}

	// Outputs
	outputs := ConvertOutputVariables(sp.OutputVariables)

	dst := &v1.StepTest{
		Container:    container,
		Env:          sp.EnvVariables,
		Intelligence: intelligence,
		Match:        match,
		Outputs:      outputs,
		Report:       report,
		Shell:        shell,
	}

	if script != "" {
		dst.Script = v1.Stringorslice{script}
	}

	return dst
}

// buildTestScript generates the complete test script from v0 RunTests fields
// Combines preCommand + buildTool command + args + postCommand
func buildTestScript(language, buildTool, preCommand, args, postCommand string) string {
	var parts []string

	// Add preCommand if present
	if preCommand != "" {
		parts = append(parts, preCommand)
	}

	// Generate the test command based on language and buildTool
	testCmd := generateTestCommand(language, buildTool, args)
	if testCmd != "" {
		parts = append(parts, testCmd)
	}

	// Add postCommand if present
	if postCommand != "" {
		parts = append(parts, postCommand)
	}

	return strings.Join(parts, "\n")
}

// generateTestCommand generates the appropriate test command based on language and buildTool
func generateTestCommand(language, buildTool, args string) string {
	buildTool = strings.ToLower(buildTool)

	var cmd string
	switch buildTool {
	// Java/Kotlin/Scala build tools
	case "maven":
		cmd = "mvn test"
	case "gradle":
		cmd = "./gradlew test"
	case "bazel":
		cmd = "bazel test"
	case "sbt":
		cmd = "sbt test"

	// C# build tools
	case "dotnet":
		cmd = "dotnet test"
	case "nunitconsole":
		cmd = "nunit3-console"

	// Python build tools
	case "pytest":
		cmd = "pytest"
	case "unittest":
		cmd = "python -m unittest"

	// Ruby build tools
	case "rspec":
		cmd = "rspec"

	default:
		// check for expression
		if strings.Contains(buildTool, "<+") {
			fmt.Printf("Expression buildTool %v is not supported for conversion\n", buildTool)
		} else {
			fmt.Printf("Unknown build tool %v for language %v in RunTests step\n", buildTool, language)
		}
		return ""
	}

	// Append args if present
	if args != "" {
		cmd = cmd + " " + args
	}

	return cmd
}

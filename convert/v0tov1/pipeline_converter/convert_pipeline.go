package pipelineconverter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	convert_helpers "github.com/drone/go-convert/convert/v0tov1/convert_helpers"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

type StepInfo struct {
	Type   string // v0 step type (e.g., "Run", "Action", "K8sRollingDeploy")
	V0Path string // full v0 path (e.g., "pipeline.stages.build.spec.execution.steps.compile")
	V1Path string // full v1 path (e.g., "pipeline.stages.build.steps.compile")
}

type PipelineConverter struct {
	stageCtx    *convert_helpers.StageConversionContext
	stepTypeMap map[string]*StepInfo // maps step ID to step info (type + v0 path)
}

func NewPipelineConverter() *PipelineConverter {
	return &PipelineConverter{
		stageCtx:    convert_helpers.NewStageConversionContext(),
		stepTypeMap: make(map[string]*StepInfo),
	}
}

// GetStepTypeMap returns the accumulated step ID to step info mapping.
func (c *PipelineConverter) GetStepTypeMap() map[string]*StepInfo {
	return c.stepTypeMap
}

// ConvertPipeline converts a v0 Pipeline to a v1 Pipeline.
func (c *PipelineConverter) ConvertPipeline(src *v0.Pipeline) *v1.Pipeline {
	if src == nil {
		return nil
	}

	dst := &v1.Pipeline{
		Id:   src.ID,
		Name: src.Name,
	}

	// Check for Template - if template exists, set it and continue converting rest
	if src.Template != nil {
		dst.Template = c.convertPipelineTemplate(src)
	}

	var barriers []string
	if src.FlowControl != nil {
		barriers = convertBarriers(src.FlowControl.Barriers)
	}

	inputs := c.convertVariables(src.Variables)
	stages := c.convertStages(src.Stages, "pipeline")

	clone := c.convertCodebase(src.Props.CI.Codebase)
	dst.Inputs = inputs
	dst.Stages = stages
	dst.Barriers = barriers
	dst.Clone = clone
	dst.Notifications = convert_helpers.ConvertNotifications(src.NotificationRules)
	dst.Delegate = convert_helpers.ConvertDelegate(src.DelegateSelectors, nil)

	return dst
}

func (c *PipelineConverter) convertCodebase(src *v0.Codebase) *v1.Clone {
	if src == nil {
		return &v1.Clone{
			Enabled: false,
		}
	}

	clone := &v1.Clone{
		Enabled:   true,
		Repo:      src.Name,
		Connector: src.Conn,
	}

	// Handle Build field - can be either a string expression or a Build struct
	if !src.Build.IsNil() {
		if build, ok := src.Build.AsStruct(); ok {
			// Build is a struct with Type and Spec
			cloneRef := &v1.CloneRef{}

			// Extract name from Spec based on type
			if build.Type == "branch" && build.Spec.Branch != "" {
				cloneRef.Name = build.Spec.Branch
				cloneRef.Type = "branch"
			} else if build.Type == "tag" && build.Spec.Tag != "" {
				cloneRef.Name = build.Spec.Tag
				cloneRef.Type = "tag"
			} else if build.Type == "PR" && build.Spec.Number != nil {
				cloneRef.Number = build.Spec.Number
				cloneRef.Type = "pull-request"
			} else if build.Type == "commitSha" && build.Spec.CommitSha != "" {
				cloneRef.Sha = build.Spec.CommitSha
				cloneRef.Type = "commit"
			}

			clone.Ref = cloneRef
		}
	}
	clone.Depth = src.Depth
	clone.Lfs = src.Lfs

	clone.Tags = src.FetchTags
	clone.Trace = src.Debug
	clone.CloneDir = src.CloneDirectory

	clone.Submodules = src.SubmoduleStrategy

	switch src.PrCloneStrategy {
	case "MergeCommit":
		clone.Strategy = "merge"
	case "SourceBranch":
		clone.Strategy = "source-branch"
	}

	clone.Insecure = src.SslVerify

	if src.Resources != nil && src.Resources.Limits != nil {
		clone.Resources = &v1.Resources{
			Limits: &v1.Limits{
				CPU:    src.Resources.Limits.GetCPUString(),
				Memory: src.Resources.Limits.GetMemoryString(),
			},
		}
	}

	return clone
}

// convertVariables converts a list of v0 Variables to v1 Inputs.
func (c *PipelineConverter) convertVariables(src []*v0.Variable) map[string]*v1.Input {
	if len(src) == 0 {
		return nil
	}

	dst := make(map[string]*v1.Input)

	for _, variable := range src {
		if variable == nil || variable.Name == "" {
			continue
		}

		v1Type := convertVariableType(variable.Type)

		input := &v1.Input{
			Type:     v1Type,
			Required: variable.Required,
		}

		if !isEmptyValue(variable.Value) {
			valueStr, ok := variable.Value.(string)
			if ok {
				parsed := parseInputExpression(valueStr)

				if parsed.value != "" {
					input.Value = formatValueForV1Type(parsed.value, v1Type)
				}
				if parsed.defaultVal != "" {
					input.Default = formatValueForV1Type(parsed.defaultVal, v1Type)
				}
				if len(parsed.enum) > 0 {
					input.Enum = formatEnumForV1Type(parsed.enum, v1Type)
				}
				if parsed.regex != "" {
					input.Pattern = parsed.regex
				}
				if parsed.executionInput {
					input.ExecutionInput = true
				}
			} else {
				input.Value = formatValueForV1Type(variable.Value, v1Type)
			}
		}

		// YAML-level default is a fallback when inline .default() wasn't found
		if !isEmptyValue(variable.Default) && isEmptyValue(input.Default) {
			input.Default = formatValueForV1Type(variable.Default, v1Type)
		}

		if variable.Type == "Secret" {
			input.Mask = true
		}

		dst[variable.Name] = input
	}

	return dst
}

type parsedExpression struct {
	value          string
	defaultVal     string
	enum           []string
	regex          string
	executionInput bool
}

var lastFuncPattern = regexp.MustCompile(`^(.*)\.(default|allowedValues|selectOneFrom|selectManyFrom|regex)\((.+)\)$`)

func parseInputExpression(valueStr string) *parsedExpression {
	result := &parsedExpression{}

	for {
		// Try .executionInput() (no args)
		if strings.HasSuffix(valueStr, ".executionInput()") {
			result.executionInput = true
			valueStr = strings.TrimSuffix(valueStr, ".executionInput()")
			continue
		}

		matches := lastFuncPattern.FindStringSubmatch(valueStr)
		if matches == nil {
			break
		}

		prefix := matches[1]
		funcName := matches[2]
		args := matches[3]

		switch funcName {
		case "default":
			result.defaultVal = args
		case "allowedValues", "selectOneFrom", "selectManyFrom":
			result.enum = parseCommaSeparatedParams(args)
		case "regex":
			result.regex = args
		}

		valueStr = prefix
	}

	result.value = valueStr
	return result
}

func parseCommaSeparatedParams(paramsStr string) []string {
	var params []string
	var current strings.Builder

	for i := 0; i < len(paramsStr); i++ {
		char := paramsStr[i]

		if char == ',' {
			param := strings.TrimSpace(current.String())
			if param != "" {
				params = append(params, param)
			}
			current.Reset()
		} else {
			current.WriteByte(char)
		}
	}

	param := strings.TrimSpace(current.String())
	if param != "" {
		params = append(params, param)
	}

	return params
}

func isEmptyValue(v interface{}) bool {
	if v == nil {
		return true
	}
	if s, ok := v.(string); ok && s == "" {
		return true
	}
	return false
}

// formatValueForV1Type coerces a value to match the declared v1 type.
// For string/secret types, ensures the result is always a string.
// For number type, attempts numeric conversion.
func formatValueForV1Type(value interface{}, v1Type string) interface{} {
	if value == nil {
		return nil
	}
	switch v1Type {
	case "string", "secret":
		switch v := value.(type) {
		case string:
			return v
		case bool:
			return fmt.Sprintf("%v", v)
		case int:
			return strconv.Itoa(v)
		case int64:
			return strconv.FormatInt(v, 10)
		case float64:
			if v == float64(int64(v)) {
				return strconv.FormatInt(int64(v), 10)
			}
			return strconv.FormatFloat(v, 'f', -1, 64)
		default:
			return fmt.Sprintf("%v", v)
		}
	case "number":
		switch v := value.(type) {
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
			return v
		default:
			return v
		}
	default:
		return value
	}
}

// formatEnumForV1Type coerces enum values to match the declared v1 type.
func formatEnumForV1Type(params []string, v1Type string) interface{} {
	switch v1Type {
	case "number":
		floatParams := make([]float64, 0, len(params))
		for _, p := range params {
			if f, err := strconv.ParseFloat(p, 64); err == nil {
				floatParams = append(floatParams, f)
			}
		}
		if len(floatParams) > 0 {
			return floatParams
		}
		return params
	default:
		return params
	}
}

func convertVariableType(v0Type string) string {
	switch v0Type {
	case "Secret":
		return "secret"
	case "Text":
		return "string"
	case "String":
		return "string"
	case "Number":
		return "number"
	default:
		return "string"
	}
}

// convertBarriers converts a list of v0 Barriers to v1 Barriers.
func convertBarriers(src []*v0.Barrier) []string {
	if len(src) == 0 {
		return nil
	}
	dst := make([]string, 0, len(src))
	for _, barrier := range src {
		if barrier == nil {
			continue
		}
		dst = append(dst, barrier.Name)
	}
	return dst
}

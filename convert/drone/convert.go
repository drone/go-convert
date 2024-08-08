// Copyright 2022 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package drone converts Drone pipelines to Harness pipelines.
package drone

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	v1 "github.com/drone/go-convert/convert/drone/yaml"
	v2 "github.com/drone/spec/dist/go"

	"github.com/drone/go-convert/internal/store"
	"github.com/ghodss/yaml"
)

// conversion context
type context struct {
	pipeline []*v1.Pipeline
	stage    *v1.Pipeline
}

// Converter converts a Drone pipeline to a Harness
// v1 pipeline.
type Converter struct {
	kubeEnabled   bool
	kubeNamespace string
	kubeConnector string
	dockerhubConn string
	identifiers   *store.Identifiers
	orgSecrets    []string
}

var variableMap = map[string]string{
	"DRONE_BRANCH":             "<+codebase.branch>",
	"DRONE_BUILD_NUMBER":       "<+pipeline.sequenceId>",
	"DRONE_COMMIT_AUTHOR":      "<+codebase.gitUserId>",
	"DRONE_COMMIT_BRANCH":      "<+codebase.branch>",
	"DRONE_COMMIT_SHA":         "<+codebase.commitSha>",
	"DRONE_PULL_REQUEST":       "<+codebase.prNumber>",
	"DRONE_PULL_REQUEST_TITLE": "<+codebase.prTitle>",
	"DRONE_REMOTE_URL":         "<+codebase.repoUrl>",
	"DRONE_REPO_NAME":          "<+<+codebase.repoUrl>.substring(<+codebase.repoUrl>.lastIndexOf('/') + 1)>",
	"CI_BUILD_NUMBER":          "<+pipeline.sequenceId>",
	"CI_COMMIT_AUTHOR":         "<+codebase.gitUserId>",
	"CI_COMMIT_BRANCH":         "<+codebase.branch>",
	"CI_COMMIT_SHA":            "<+codebase.commitSha>",
	"CI_REMOTE_URL":            "<+codebase.repoUrl>",
	"CI_REPO_NAME":             "<+<+codebase.repoUrl>.substring(<+codebase.repoUrl>.lastIndexOf('/') + 1)>",
}

// New creates a new Converter that converts a Drone
// pipeline to a Harness v1 pipeline.
func New(options ...Option) *Converter {
	d := new(Converter)

	// create the unique identifier store. this store
	// is used for registering unique identifiers to
	// prevent duplicate names, unique index violations.
	d.identifiers = store.New()

	// loop through and apply the options.
	for _, option := range options {
		option(d)
	}

	// set the default kubernetes namespace.
	if d.kubeNamespace == "" {
		d.kubeNamespace = "default"
	}

	// set the runtime to kubernetes if the kubernetes
	// connector is configured.
	if d.kubeConnector != "" {
		d.kubeEnabled = true
	}

	return d
}

// Convert downgrades a v1 pipeline.
func (d *Converter) Convert(r io.Reader) ([]byte, error) {
	src, err := v1.Parse(r)
	if err != nil {
		return nil, err
	}
	return d.convert(&context{
		pipeline: src,
	})
}

// ConvertString downgrades a v1 pipeline.
func (d *Converter) ConvertBytes(b []byte) ([]byte, error) {
	return d.Convert(
		bytes.NewBuffer(b),
	)
}

// ConvertString downgrades a v1 pipeline.
func (d *Converter) ConvertString(s string) ([]byte, error) {
	return d.Convert(
		bytes.NewBufferString(s),
	)
}

// ConvertFile downgrades a v1 pipeline.
func (d *Converter) ConvertFile(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return d.Convert(f)
}

// converts converts a Drone pipeline to a Harness pipeline.
func (d *Converter) convert(ctx *context) ([]byte, error) {

	//
	// TODO convert env substitution to expression
	//

	//
	// TODO convert from_secret to expression
	//

	// create the pipeline spec
	pipeline := &v2.Pipeline{
		Options: &v2.Default{
			Registry: convertRegistry(ctx.pipeline),
		},
	}

	// create the harness pipeline resource
	config := &v2.Config{
		Version: 1,
		Kind:    "pipeline",
		Spec:    pipeline,
	}

	for _, from := range ctx.pipeline {
		if from == nil {
			continue
		}

		switch from.Kind {
		case v1.KindSecret: // TODO
		case v1.KindSignature: // TODO
		case v1.KindPipeline:
			// TODO pipeline.name removed from spec
			// pipeline.Name = from.Name

			pipeline.Stages = append(pipeline.Stages, &v2.Stage{
				Name:     from.Name,
				Type:     "ci",
				When:     convertCond(from.Trigger),
				Delegate: convertNode(from.Node),
				Spec: &v2.StageCI{
					Clone:    convertClone(from.Clone),
					Envs:     copyenv(from.Environment),
					Platform: convertPlatform(from.Platform),
					Runtime:  convertRuntime(from),
					Steps:    convertSteps(from, d.orgSecrets),
					Volumes:  convertVolumes(from.Volumes),

					// TODO support for delegate.selectors from from.Node
					// TODO support for stage.variables
				},
			})
		}
	}

	// marshal the harness yaml
	out, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	// drone had a bug where it required double-escaping
	// which can be eliminated going forward.
	// TODO this will probably require some tweaking.
	out = bytes.ReplaceAll(out, []byte(`\\\\`), []byte(`\\`))

	// Replace all occurrences of /drone/src with /harness
	out = bytes.ReplaceAll(out, []byte("/drone/src"), []byte("/harness"))

	return out, nil
}

func convertRegistry(src []*v1.Pipeline) *v2.Registry {
	// note that registry credentials in Drone are stored
	// at the stage level, but in Harness, we are proposing
	// they are stored at the pipeline level (this could
	// change in the future).
	//
	// this means we need to the combined, unique list
	// of pull secrets across all stages.

	set := map[string]struct{}{}
	for _, v := range src {
		if v == nil {
			continue
		}
		if len(v.PullSecrets) == 0 {
			continue
		}
		for _, s := range v.PullSecrets {
			set[s] = struct{}{}
		}
	}
	if len(set) == 0 {
		return nil
	}
	dst := &v2.Registry{}
	for k := range set {
		dst.Connector = append(dst.Connector, &v2.RegistryConnector{
			Name: k,
		})
	}
	return dst
}

func convertSteps(src *v1.Pipeline, orgSecrets []string) []*v2.Step {
	var dst []*v2.Step
	for _, v := range src.Services {
		if v != nil {
			dst = append(dst, convertBackground(v, orgSecrets))
		}
	}
	for _, v := range src.Steps {
		if v != nil {
			switch {
			case v.Detach:
				dst = append(dst, convertBackground(v, orgSecrets))
			case isPlugin(v):
				dst = append(dst, convertPlugin(v, orgSecrets))
			default:
				dst = append(dst, convertRun(v, orgSecrets))
			}
		}
	}
	return dst
}

func convertPlugin(src *v1.Step, orgSecrets []string) *v2.Step {
	return &v2.Step{
		Name: src.Name,
		Type: "plugin",
		When: convertCond(src.When),
		Spec: &v2.StepPlugin{
			Image:      src.Image,
			Mount:      convertMounts(src.Volumes),
			Privileged: src.Privileged,
			Pull:       convertPull(src.Pull),
			User:       src.User,
			Envs:       convertVariables(src.Environment, orgSecrets),
			With:       convertSettings(src.Settings, orgSecrets),
			Resources:  convertResourceLimits(&src.Resource),
			// Volumes       // FIX
		},
	}
}

func convertBackground(src *v1.Step, orgSecrets []string) *v2.Step {
	return &v2.Step{
		Name: src.Name,
		Type: "background",
		When: convertCond(src.When),
		Spec: &v2.StepBackground{
			Image:      src.Image,
			Mount:      convertMounts(src.Volumes),
			Privileged: src.Privileged,
			Pull:       convertPull(src.Pull),
			Shell:      convertShell(src.Shell),
			User:       src.User,
			Entrypoint: convertEntrypoint(src.Entrypoint),
			Args:       convertArgs(src.Entrypoint, src.Command),
			Run:        convertScript(src.Commands),
			Envs:       convertVariables(src.Environment, orgSecrets),
			Resources:  convertResourceLimits(&src.Resource),
			// Volumes       // FIX
		},
	}
}

func convertRun(src *v1.Step, orgSecrets []string) *v2.Step {
	return &v2.Step{
		Name: src.Name,
		Type: "script",
		When: convertCond(src.When),
		Spec: &v2.StepExec{
			Image:      src.Image,
			Mount:      convertMounts(src.Volumes),
			Privileged: src.Privileged,
			Pull:       convertPull(src.Pull),
			Shell:      convertShell(src.Shell),
			User:       src.User,
			Entrypoint: convertEntrypoint(src.Entrypoint),
			Args:       convertArgs(src.Entrypoint, src.Command),
			Run:        convertScript(src.Commands),
			Envs:       convertVariables(src.Environment, orgSecrets),
			Resources:  convertResourceLimits(&src.Resource),
			// Volumes       // FIX
		},
	}
}

func convertResourceLimits(src *v1.Resources) *v2.Resources {
	if src.Limits.CPU == 0 && src.Limits.Memory == 0 {
		return nil
	}
	return &v2.Resources{
		Limits: &v2.Resource{
			Cpu:    v2.StringorInt(src.Requests.CPU),
			Memory: v2.MemStringorInt(src.Requests.Memory),
		},
	}
}

func convertResourceRequests(src *v1.Resources) *v2.Resources {
	if src.Requests.CPU == 0 && src.Requests.Memory == 0 {
		return nil
	}
	return &v2.Resources{
		Requests: &v2.Resource{
			Cpu:    v2.StringorInt(src.Requests.CPU),
			Memory: v2.MemStringorInt(src.Requests.Memory),
		},
	}
}

func convertEntrypoint(src []string) string {
	if len(src) == 0 {
		return ""
	} else {
		return src[0]
	}
}

func convertVariables(src map[string]*v1.Variable, orgSecrets []string) map[string]string {
	dst := map[string]string{}

	orgSecretsMap := make(map[string]bool, len(orgSecrets))
	for _, secret := range orgSecrets {
		orgSecretsMap[secret] = true
	}

	for k, v := range src {
		switch {
		case v.Value != "":
			dst[sanitizeString(k)] = replaceVars(v.Value)
		case v.Secret != "":
			secretID := sanitizeString(v.Secret)
			if _, exists := orgSecretsMap[secretID]; exists {
				secretID = "org." + secretID
			}
			dst[k] = fmt.Sprintf("<+secrets.getValue(%q)>", secretID) // TODO figure out secret syntax
		}
	}
	return dst
}

func sanitizeString(input string) string {
	// Regular expression to match all characters that are not a letter, number, or underscore
	reg, _ := regexp.Compile("[^a-zA-Z0-9_]+")
	sanitized := reg.ReplaceAllString(input, "_")
	return sanitized
}

func convertVolumes(src []*v1.Volume) []*v2.Volume {
	var dst []*v2.Volume
	for _, v := range src {
		if v == nil || v.Name == "" {
			continue
		}
		switch {
		case v.EmptyDir != nil:
			dst = append(dst, &v2.Volume{
				Name: v.Name,
				Type: "temp",
				Spec: &v2.VolumeTemp{
					// TODO convert medium and limit
				},
			})
		case v.HostPath != nil:
			dst = append(dst, &v2.Volume{
				Name: v.Name,
				Type: "host",
				Spec: &v2.VolumeHost{
					Path: v.HostPath.Path,
				},
			})
		}
	}
	if len(dst) == 0 {
		return nil
	}
	return dst
}

func convertMounts(src []*v1.VolumeMount) []*v2.Mount {
	var dst []*v2.Mount
	for _, v := range src {
		if v == nil || v.Name == "" || v.MountPath == "" {
			continue
		}
		dst = append(dst, &v2.Mount{
			Name: v.Name,
			Path: v.MountPath,
		})
	}
	if len(dst) == 0 {
		return nil
	}
	return dst
}

func convertSettings(src map[string]*v1.Parameter, orgSecrets []string) map[string]interface{} {
	dst := map[string]interface{}{}

	orgSecretsMap := make(map[string]bool, len(orgSecrets))
	for _, secret := range orgSecrets {
		orgSecretsMap[secret] = true
	}

	for k, v := range src {
		switch {
		case v.Secret != "":
			secretID := sanitizeString(v.Secret)
			if _, exists := orgSecretsMap[secretID]; exists {
				secretID = "org." + secretID
			}
			dst[k] = fmt.Sprintf("<+secrets.getValue(%q)>", secretID)
		case v.Value != nil:
			dst[k] = convertInterface(v.Value)
		}
	}
	return dst
}

func convertInterface(i interface{}) interface{} {
	switch v := i.(type) {
	case string:
		return replaceVars(v)
	case map[interface{}]interface{}:
		newMap := make(map[string]interface{})
		for key, value := range v {
			keyStr, ok := key.(string)
			if !ok {
				continue
			}
			newMap[keyStr] = convertInterface(value)
		}
		return newMap
	case []interface{}:
		for i, value := range v {
			v[i] = convertInterface(value)
		}
	}
	return i
}

func replaceVars(val string) string {
	vars := strings.Split(val, " ") // required for combine vars

	for i, v := range vars {
		if containsReplacement(v) {
			vars[i] = replaceCharacters(v)
		} else if containsSubstring(v) {
			vars[i] = processSubstring(v)
		} else {
			// simple variable substitution
			vars[i] = replaceSimpleVar(v)
		}
	}

	return strings.Join(vars, " ")
}

func replaceSimpleVar(val string) string {
	var re = regexp.MustCompile(`\$\$?({)?(\w+)(})?`)

	return re.ReplaceAllStringFunc(val, func(match string) string {
		varName := strings.Trim(match, "${}$")
		if harnessVar, ok := variableMap[varName]; ok {
			return harnessVar
		}
		return match
	})
}

// Check if variable contains a replacement operation
func containsReplacement(v string) bool {
	var re = regexp.MustCompile(`\$\{(\w+)(/|//)([^/]+)/([^}]+)\}`)
	return re.MatchString(v)
}

// Perform the replacement operation
func replaceCharacters(match string) string {
	// Initialize the regular expression to match "${...}"
	var re = regexp.MustCompile(`\$\{([^}/]+)(/|//)([^/]+)/([^}]+)\}`)

	return re.ReplaceAllStringFunc(match, func(m string) string {
		groups := re.FindStringSubmatch(m)

		if len(groups) < 5 {
			return m
		}

		varName := groups[1]
		separator := groups[2]
		oldChar := strings.ReplaceAll(groups[3], `\\`, `\`) // Replace escaped backslash
		newChar := strings.ReplaceAll(groups[4], `\/`, `/`) // Replace escaped forward slash

		if strings.HasPrefix(newChar, "/") && len(newChar) > 1 {
			newChar = strings.TrimPrefix(newChar, "/")
		}

		if oldChar == "\\" {
			oldChar = "/"
		}

		if harnessVar, ok := variableMap[varName]; ok {
			if separator == "//" {
				return "<+" + strings.Trim(harnessVar, "<+>") + ".replace('" + oldChar + "', '" + newChar + "')>"
			} else {
				return "<+" + strings.Trim(harnessVar, "<+>") + ".replaceFirst('" + oldChar + "', '" + newChar + "')>"
			}
		}

		return m
	})
}

// Check if variable contains a substring operation
func containsSubstring(v string) bool {
	var re = regexp.MustCompile(`\$\$?({)?(\w+)((:\d+)?(:\d+)?)?(})?`)
	return re.MatchString(v)
}

// Perform the substring operation
func processSubstring(val string) string {
	var re = regexp.MustCompile(`\$\$?({)?(\w+)((:\d+)?(:\d+)?)?(})?`)

	return re.ReplaceAllStringFunc(val, func(match string) string {
		// Remove special characters from match and split on ":"
		parts := strings.Split(strings.Trim(match, "${}$"), ":")

		varName := parts[0]
		if harnessVar, ok := variableMap[varName]; ok {
			// If there are substring operations
			if len(parts) > 1 {
				offset := parts[1]
				length := ""
				if len(parts) > 2 {
					length = parts[2]
				}
				// Modify the harnessVar to add the substring operation
				return "<+" + strings.Trim(harnessVar, "<+>") + ".substring(" + offset + "," + length + ")>"
			}
			return harnessVar
		}
		return match
	})
}

func convertScript(src []string) string {
	if len(src) == 0 {
		return ""
	}

	for i, cmd := range src {
		src[i] = replaceVars(cmd)
	}

	return strings.Join(src, "\n")
}

func convertArgs(src1, src2 []string) []string {
	if len(src1) == 0 {
		return src2
	} else {
		return append(src1[:1], src2...)
	}
}

func convertPull(src string) string {
	switch src {
	case "always":
		return "always"
	case "never":
		return "never"
	case "if-not-exists":
		return "if-not-exists"
	default:
		return ""
	}
}

func convertShell(src string) string {
	switch src {
	case "bash":
		return "bash"
	case "sh", "posix":
		return "sh"
	case "pwsh", "powershell":
		return "powershell"
	default:
		return ""
	}
}

func convertRuntime(src *v1.Pipeline) *v2.Runtime {
	if src.Type == "kubernetes" {
		return &v2.Runtime{
			Type: "kubernetes",
			Spec: &v2.RuntimeKube{
				// TODO should harness support `dns_config`
				// TODO should harness support `host_aliases`
				// TODO support for `tolerations`
				Annotations:    src.Metadata.Annotations,
				Labels:         src.Metadata.Labels,
				Namespace:      src.Metadata.Namespace,
				NodeSelector:   src.NodeSelector,
				Node:           src.NodeName,
				ServiceAccount: src.ServiceAccount,
				Resources:      convertResourceRequests(&src.Resource),
			},
		}
	}
	return &v2.Runtime{
		Type: "machine",
		Spec: v2.RuntimeMachine{},
	}
}

func convertClone(src v1.Clone) *v2.CloneStage {
	dst := new(v2.CloneStage)
	if v := src.Depth; v != 0 {
		dst.Depth = int64(v)
	}
	if v := src.Disable; v {
		dst.Disabled = true
	}
	if v := src.SkipVerify; v {
		dst.Insecure = true
	}
	if v := src.Trace; v {
		dst.Trace = true
	}
	return dst
}

func convertNode(src map[string]string) []string {
	if len(src) == 0 {
		return nil
	}
	var dst []string
	for k, v := range src {
		dst = append(
			dst, k+":"+v)
	}
	return dst
}

func convertPlatform(src v1.Platform) *v2.Platform {
	if src.Arch == "" && src.OS == "" {
		return nil
	}
	dst := new(v2.Platform)
	switch src.OS {
	case "windows", "win", "win32":
		dst.Os = v2.OSWindows.String()
	case "darwin", "macos", "mac":
		dst.Os = v2.OSDarwin.String()
	default:
		dst.Os = v2.OSLinux.String()
	}
	switch src.Arch {
	case "arm", "arm64":
		dst.Arch = v2.ArchArm64.String()
	default:
		dst.Arch = v2.ArchAmd64.String()
	}
	return dst
}

func convertCond(src v1.Conditions) *v2.When {
	if isCondsEmpty(src) {
		return nil
	}

	exprs := map[string]*v2.Expr{}
	if expr := convertExpr(src.Action); expr != nil {
		exprs["action"] = expr
	}
	if expr := convertExpr(src.Branch); expr != nil {
		exprs["branch"] = expr
	}
	if expr := convertExpr(src.Cron); expr != nil {
		exprs["cron"] = expr
	}
	if expr := convertExpr(src.Event); expr != nil {
		exprs["event"] = expr
	}
	if expr := convertExpr(src.Instance); expr != nil {
		exprs["instance"] = expr
	}
	if expr := convertExpr(src.Paths); expr != nil {
		exprs["paths"] = expr
	}
	if expr := convertExpr(src.Ref); expr != nil {
		exprs["ref"] = expr
	}
	if expr := convertExpr(src.Repo); expr != nil {
		exprs["repo"] = expr
	}
	if expr := convertExpr(src.Status); expr != nil {
		exprs["status"] = expr
	}
	if expr := convertExpr(src.Target); expr != nil {
		exprs["target"] = expr
	}

	dst := new(v2.When)
	dst.Cond = []map[string]*v2.Expr{exprs}
	return dst
}

func convertExpr(src v1.Condition) *v2.Expr {
	if len(src.Include) != 0 {
		return &v2.Expr{In: src.Include}
	}
	if len(src.Exclude) != 0 {
		return &v2.Expr{
			Not: &v2.Expr{In: src.Exclude},
		}
	}
	return nil
}

func isCondsEmpty(src v1.Conditions) bool {
	return isCondEmpty(src.Action) &&
		isCondEmpty(src.Action) &&
		isCondEmpty(src.Branch) &&
		isCondEmpty(src.Cron) &&
		isCondEmpty(src.Event) &&
		isCondEmpty(src.Instance) &&
		isCondEmpty(src.Paths) &&
		isCondEmpty(src.Ref) &&
		isCondEmpty(src.Repo) &&
		isCondEmpty(src.Status) &&
		isCondEmpty(src.Target)
}

func isCondEmpty(src v1.Condition) bool {
	return len(src.Exclude) == 0 && len(src.Include) == 0
}

func isPlugin(src *v1.Step) bool {
	return len(src.Settings) > 0
}

// copyenv returns a copy of the environment variable map.
func copyenv(src map[string]string) map[string]string {
	dst := map[string]string{}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

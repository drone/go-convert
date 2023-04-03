package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/drone/go-convert/convert/circle/converter/circleci/config"
)

func ConvertEnvs(in map[string]*config.Xcode) map[string]string {
	env := make(map[string]string)
	for k, v := range in {
		if v.String != nil {
			env[k] = *v.String
		}
		if v.Double != nil {
			env[k] = fmt.Sprintf("%v", *v.Double)
		}
	}
	return env
}

func ConvertEnvs_(in map[string]*config.Environment) map[string]string {
	env := make(map[string]string)
	for k, v := range in {
		if v.String != nil {
			env[k] = *v.String
		}
		if v.Double != nil {
			env[k] = fmt.Sprintf("%v", *v.Double)
		}
		if v.Bool != nil {
			env[k] = fmt.Sprintf("%v", *v.Bool)
		}
	}
	return env
}

func ReplaceSecret(s string, prefix string) string {
	secret := fmt.Sprintf("replace-%s", prefix)
	if s != "" {
		secret = s
	}
	return fmt.Sprintf("<+secrets.getValue(\"%s\")>", secret)
}

func ReplaceString(s string, prefix string) string {
	val := fmt.Sprintf("replace-%s", prefix)
	if s != "" {
		val = s
	}
	return val
}

func ConvertIfcMapToStrMap(in map[string]interface{}) map[string]string {
	m := make(map[string]string)
	for k, v := range in {
		m[k] = fmt.Sprintf("%v", v)
	}
	return m
}

func GetParamDefaults(in map[string]interface{}) map[string]string {
	m := make(map[string]string)
	for k, v := range in {
		if vv_, ok := v.(map[string]interface{}); ok {
			if d, ok := vv_["default"]; ok {
				m[k] = fmt.Sprintf("%v", d)
			}
		}
	}
	return m
}

func ResolveParams(params map[string]string, inputs map[string]string) map[string]string {
	m := make(map[string]string)
	for k, v := range params {
		if v_, ok := inputs[k]; ok {
			m[k] = v_
		} else {
			m[k] = v
		}
	}
	return m
}

func ResolveStrExpr(s string, in map[string]string) string {
	re := regexp.MustCompile(`<<\s*(parameters|pipeline.parameters)\.([A-z0-9-_]*)\s*>>`)
	matches := re.FindAllStringSubmatch(s, -1)

	out := s
	for _, m := range matches {
		if len(m) == 3 {
			mstr := m[0]
			param := m[2]

			if m[1] == "parameters" {
				if val, ok := in[param]; ok {
					out = strings.Replace(out, mstr, val, -1)
				}
			} else if m[1] == "pipeline.parameters" {
				d := fmt.Sprintf("<+pipeline.variables.%s", param)
				out = strings.Replace(out, mstr, d, -1)
			}
		}
	}
	return out
}

func ResolveMapExpr(m map[string]string, in map[string]string) map[string]string {
	o := make(map[string]string, 0)
	for k, v := range m {
		o[ResolveStrExpr(k, in)] = ResolveStrExpr(v, in)
	}
	return o
}

func ResolveListExpr(m []string, in map[string]string) []string {
	o := make([]string, 0)
	for _, v := range m {
		o = append(o, ResolveStrExpr(v, in))
	}
	return o
}

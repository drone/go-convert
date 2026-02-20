package converthelpers

import (
	"fmt"
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

func ConvertStepBitrise(src *v0.Step) (*v1.StepRun) {
	if src == nil || src.Spec == nil {
	return nil
	}
	sp, ok := src.Spec.(*v0.StepBitrise)
	if !ok {
		return nil
	}
	script := fmt.Sprintf("plugin -kind bitrise -name %v", sp.Uses)
	
	env_map := map[string]interface{}{}
	var env *flexible.Field[map[string]interface{}]
	for k, v := range sp.With {
		env_map[k] = v
	}
	for k, v := range sp.Envs {
		env_map[k] = v
	} 	
	if len(env_map)>0 {
		env = &flexible.Field[map[string]interface{}]{Value: env_map}
	}

	return &v1.StepRun{
		Script: v1.Stringorslice{script},
		Env:    env,
	}

}

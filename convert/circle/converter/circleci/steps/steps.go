package steps

import (
	"fmt"
	"strings"

	"github.com/drone/go-convert/convert/circle/commons"
	"github.com/drone/go-convert/convert/circle/converter/circleci/config"
	"github.com/drone/go-convert/convert/circle/converter/circleci/orbs"
	"github.com/drone/go-convert/convert/circle/converter/circleci/steps/artifact"
	"github.com/drone/go-convert/convert/circle/converter/circleci/steps/cache"
	"github.com/drone/go-convert/convert/circle/converter/circleci/steps/run"
	"github.com/drone/go-convert/convert/circle/converter/circleci/steps/testreport"
	"github.com/drone/go-convert/convert/circle/converter/circleci/utils"

	harness "github.com/drone/spec/dist/go"
)

func Convert(opts commons.Opts, s config.Step, image string,
	inParams map[string]string, jorbName string,
	orbMapper map[string]config.ORBClass) ([]*harness.Step, error) {
	steps := make([]*harness.Step, 0)

	name, sinputs := getStepCmdInfo(s)
	if name != "" {
		orbCmd, ok := getOrbCmd(name, jorbName, orbMapper)
		if !ok {
			return nil, nil
		}

		sinputs = utils.ResolveMapExpr(sinputs, inParams)
		cmdParams := utils.ResolveParams(
			utils.GetParamDefaults(orbCmd.Parameters), sinputs)
		for _, s_ := range orbCmd.Steps {
			ss, err := Convert(opts, s_, image, cmdParams, jorbName, orbMapper)
			if err != nil {
				return nil, err
			}

			steps = append(steps, ss...)
		}
		return steps, nil
	}

	if s.StepClass != nil && s.StepClass.When != nil {
		// TODO: Add support for when condition at step level
		fmt.Println("When condition ignored")
		if len(s.StepClass.When.Steps.StepArray) != 0 {
			for _, s_ := range s.StepClass.When.Steps.StepArray {
				ss, err := Convert(opts, s_, image, inParams, jorbName, orbMapper)
				if err != nil {
					return nil, err
				}

				steps = append(steps, ss...)
			}
		} else if s.StepClass.When.Steps.MapSteps != nil {
			ss, err := convert(opts, s.StepClass.When.Steps.MapSteps, image, inParams)
			if err != nil {
				return nil, err
			}

			steps = append(steps, ss)
		}
		return steps, nil
	}

	step, err := convert(opts, s.StepClass, image, inParams)
	if err != nil {
		return nil, err
	}
	if step == nil {
		return nil, nil
	}
	return []*harness.Step{step}, nil
}

func convert(opts commons.Opts, s *config.StepClass, image string, inParams map[string]string) (
	*harness.Step, error) {
	if s != nil {
		if s.Run != nil {
			return run.Convert(*s.Run, image, inParams)
		} else if s.StoreTestResults != nil {
			return testreport.Convert(*s.StoreTestResults, inParams)
		} else if s.RestoreCache != nil {
			return cache.ConvertRestore(opts, *s.RestoreCache, inParams)
		} else if s.SaveCache != nil {
			return cache.ConvertSave(opts, *s.SaveCache, inParams)
		} else if s.StoreArtifacts != nil {
			return artifact.Convert(opts, *s.StoreArtifacts, inParams)
		}
	}
	return nil, nil
}

func getOrbCmd(stepName, jorbName string, orbMapper map[string]config.ORBClass) (
	*config.CommandValue, bool) {
	orbName, orbCmd := "", ""

	s := strings.Split(stepName, "/")
	if len(s) == 2 {
		orbName, orbCmd = s[0], s[1]
	} else if jorbName != "" {
		orbName, orbCmd = jorbName, stepName
	} else {
		orbCmd = stepName
	}

	return orbs.FindCmdInOrbMapper(orbName, orbCmd, orbMapper)
}

func getStepCmdInfo(s config.Step) (string, map[string]string) {
	if s.String != nil && *s.String != "" {
		return *s.String, nil
	}

	if len(s.MapClass) == 1 {
		for k_, v_ := range s.MapClass {
			return k_, utils.ConvertIfcMapToStrMap(v_)
		}
	}
	return "", nil
}

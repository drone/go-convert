package jobs

import (
	"fmt"
	"strings"

	"github.com/drone/go-convert/convert/circle/commons"
	"github.com/drone/go-convert/convert/circle/converter/circleci/config"
	"github.com/drone/go-convert/convert/circle/converter/circleci/steps"
	"github.com/drone/go-convert/convert/circle/converter/circleci/steps/background"
	"github.com/drone/go-convert/convert/circle/converter/circleci/utils"

	harness "github.com/drone/spec/dist/go"
)

func Convert(opts commons.Opts, j config.JobValue, orbName string,
	orbMapper map[string]config.ORBClass, jobInputs map[string]string,
	executors map[string]config.ExecutorValue) (
	*harness.Stage, error) {
	hci := harness.StageCI{}
	jobParams := utils.ResolveParams(utils.GetParamDefaults(j.Parameters), jobInputs)
	envs := utils.ResolveMapExpr(utils.ConvertEnvs(j.Environment), jobParams)

	executor, executorInputs := getExecutor(j, orbName, orbMapper, executors)
	// TODO: Handle executors other than docker
	var dExecutors []config.DockerExecutor
	if len(j.Docker) != 0 {
		dExecutors = j.Docker
		executorInputs = jobInputs
	} else if executor != nil && len(executor.Docker) != 0 {
		dExecutors = executor.Docker
	}

	machine := j.Machine
	if machine == nil && executor != nil {
		machine = executor.Machine
	}
	if executor != nil {
		if j.Macos != nil || executor.Macos != nil {
			hci.Platform = &harness.Platform{Os: harness.OSMacos, Arch: harness.ArchArm64}
			hci.Runtime = &harness.Runtime{Type: "cloud"}
		} else if machine != nil && machine.MachineExecutor != nil && strings.HasPrefix(machine.MachineExecutor.Image, "windows") {
			hci.Platform = &harness.Platform{Os: harness.OSWindows, Arch: harness.ArchAmd64}
			hci.Runtime = &harness.Runtime{Type: "cloud"}
		}
	}

	image := ""
	if len(dExecutors) != 0 {
		image = utils.ResolveStrExpr(dExecutors[0].Image, executorInputs)
		for k, v := range utils.ResolveMapExpr(utils.ConvertEnvs_(dExecutors[0].Environment), executorInputs) {
			envs[k] = v
		}
		hci.Steps = append(hci.Steps, &harness.Step{
			Name: "setup",
			Type: "script",
			Spec: harness.StepExec{
				Run: "mkdir -p /github  /__w/_temp; sudo chmod -R 777 /harness /github /__w/_temp",
			},
		})

		for i := 1; i < len(dExecutors); i++ {
			s, err := background.Convert(dExecutors[i], executorInputs)
			if err != nil {
				return nil, err
			}
			hci.Steps = append(hci.Steps, s)
		}
	}

	for _, step := range j.Steps {
		s, err := steps.Convert(opts, step, image, jobParams, orbName, orbMapper)
		if err != nil {
			return nil, err
		}

		if s != nil {
			hci.Steps = append(hci.Steps, s...)
		}
	}
	hci.Envs = envs

	return &harness.Stage{
		Name: "ci stage",
		Type: "ci",
		Spec: hci,
	}, nil
}

func getExecutor(j config.JobValue, orbName string,
	orbMapper map[string]config.ORBClass, executors map[string]config.ExecutorValue) (
	*config.ExecutorValue, map[string]string) {
	if j.Executor == nil {
		return nil, nil
	}

	name := ""
	inputs := make(map[string]interface{})
	if j.Executor.String != nil {
		name = *j.Executor.String
	} else if len(j.Executor.MapClass) != 0 {
		v, ok := j.Executor.MapClass["name"]
		if !ok {
			return nil, nil
		}
		name = fmt.Sprintf("%s", v)

		for k, v := range j.Executor.MapClass {
			if k == "name" {
				continue
			}
			inputs[k] = v
		}
	}

	var executor *config.ExecutorValue
	if e, ok := executors[name]; ok {
		executor = &e
	} else {
		if e_, ok := orbMapper[orbName].Executors[name]; ok {
			executor = &e_
		} else {
			fmt.Printf("Warn: Executor %s not found", name)
			return nil, nil
		}
	}

	in := utils.GetParamDefaults(executor.Parameters)
	for k, v := range utils.ConvertIfcMapToStrMap(inputs) {
		in[k] = v
	}
	return executor, in
}

package orbs

import (
	"encoding/json"

	"github.com/CircleCI-Public/circleci-yaml-language-server/pkg/parser"
	"github.com/CircleCI-Public/circleci-yaml-language-server/pkg/utils"
	"github.com/drone/go-convert/convert/circle/converter/circleci/config"
	"sigs.k8s.io/yaml"
)

func Fetch(orbs map[string]*config.ORBValue) (map[string]config.ORBClass, error) {
	context := &utils.LsContext{
		Api: utils.ApiContext{
			Token:   "XXXXXXXXXXXX",
			HostUrl: "https://circleci.com",
		},
	}
	cache := utils.CreateCache()

	m := make(map[string]config.ORBClass)
	for k, v := range orbs {
		if v.String != nil {
			orbClass, err := convert(*v.String, cache, context)
			if err != nil {
				return nil, err
			}
			m[k] = *orbClass
		}

		if v.ORBClass != nil {
			m[k] = *v.ORBClass
		}
	}
	return m, nil
}

func convert(orbVersionCode string, cache *utils.Cache, context *utils.LsContext) (*config.ORBClass, error) {
	data, err := parser.GetOrbInfo(orbVersionCode, cache, context)
	if err != nil {
		return nil, err
	}
	jdata, err := yaml.YAMLToJSON([]byte(data.Source))
	if err != nil {
		return nil, err
	}

	var r config.ORBClass
	if err := json.Unmarshal(jdata, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func FindJobInOrbMapper(orbName, jobName string,
	mapper map[string]config.ORBClass) (*config.JobValue, bool) {
	orbClass, ok := mapper[orbName]
	if !ok {
		return nil, false
	}

	jobVal, ok := orbClass.Jobs[jobName]
	if !ok {
		return nil, false
	}
	return &jobVal, true
}

func FindCmdInOrbMapper(orbName, cmd string, mapper map[string]config.ORBClass) (
	*config.CommandValue, bool) {
	orbClass, ok := mapper[orbName]
	if !ok {
		return nil, false
	}

	cmdVal, ok := orbClass.Commands[cmd]
	if !ok {
		return nil, false
	}
	return &cmdVal, true
}

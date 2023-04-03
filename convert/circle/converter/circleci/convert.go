package circleci

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/drone/go-convert/convert/circle/commons"
	"github.com/drone/go-convert/convert/circle/converter/circleci/config"
	"github.com/drone/go-convert/convert/circle/converter/circleci/jobs"
	"github.com/drone/go-convert/convert/circle/converter/circleci/orbs"
	"github.com/drone/go-convert/convert/circle/converter/circleci/utils"

	harness "github.com/drone/spec/dist/go"
	"sigs.k8s.io/yaml"
)

func Convert(opts commons.Opts, d []byte) ([]*harness.Pipeline, error) {
	jdata, err := yaml.YAMLToJSON(d)
	if err != nil {
		return nil, err
	}

	cfg, err := config.UnmarshalConfig(jdata)
	if err != nil {
		return nil, err
	}

	orbMapper, err := orbs.Fetch(cfg.Orbs)
	if err != nil {
		return nil, err
	}
	if len(cfg.Commands) != 0 {
		orbMapper[""] = config.ORBClass{
			Commands: cfg.Commands,
			Jobs:     make(map[string]config.JobValue),
		}
	}

	var pipelines []*harness.Pipeline
	for k, w := range cfg.Workflows {
		if w.Workflow == nil {
			continue
		}

		p := &harness.Pipeline{
			Version: 1,
			Name:    k,
			Inputs:  convertParams(cfg.Parameters),
		}
		for _, j := range w.Workflow.Jobs {
			name, _, inputs := getJobInfo(j)
			if name == "" {
				continue
			}

			jobVal, orbName := getJobVal(name, cfg, orbMapper)
			if jobVal == nil {
				continue
			}

			s, err := jobs.Convert(opts, *jobVal, orbName, orbMapper, inputs, cfg.Executors)
			if err != nil {
				return nil, err
			}

			s.Name = name
			p.Stages = append(p.Stages, s)
		}
		pipelines = append(pipelines, p)
	}
	return pipelines, nil
}

func getJobInfo(j config.Job) (string, *config.JobRef,
	map[string]string) {
	if j.String != nil {
		return *j.String, nil, nil
	}

	if j.MapClass != nil && len(j.MapClass) == 1 {
		for k_, v_ := range j.MapClass {
			var jobRef *config.JobRef
			parseJobRef(v_, jobRef)

			delete(v_, "matrix")
			delete(v_, "context")
			delete(v_, "filters")
			delete(v_, "requires")
			delete(v_, "type")
			return k_, jobRef, utils.ConvertIfcMapToStrMap(v_)
		}
	}
	return "", nil, nil
}

func getJobVal(name string, cfg config.Config, orbMapper map[string]config.ORBClass) (*config.JobValue, string) {
	if strings.Contains(name, "/") {
		s := strings.Split(name, "/")
		if len(s) != 2 {
			return nil, ""
		}

		orbName, jobName := s[0], s[1]
		jobVal, ok := orbs.FindJobInOrbMapper(orbName, jobName, orbMapper)
		if !ok {
			fmt.Printf("Warn: failed to find job: %s in orbs", name)
			return nil, ""
		}
		return jobVal, orbName
	}

	for k, j := range cfg.Jobs {
		if k == name {
			return &j, ""
		}
	}
	return nil, ""
}

func parseJobRef(in interface{}, p interface{}) error {
	d, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(d, p)
}

func convertParams(p map[string]interface{}) map[string]*harness.Input {
	m := make(map[string]*harness.Input)
	for k, v := range p {
		m[k] = &harness.Input{}
		if vv_, ok := v.(map[string]interface{}); ok {
			if d, ok := vv_["default"]; ok {
				m[k].Default = fmt.Sprintf("%v", d)
			}
			if d, ok := vv_["type"]; ok {
				m[k].Type = fmt.Sprintf("%v", d)
			}
		}
	}
	return m
}

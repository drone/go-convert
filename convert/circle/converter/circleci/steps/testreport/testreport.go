package testreport

import (
	"github.com/drone/go-convert/convert/circle/converter/circleci/config"
	"github.com/drone/go-convert/convert/circle/converter/circleci/utils"
	harness "github.com/drone/spec/dist/go"
)

const (
	stepType   = "script"
	namePrefix = "report"
)

func Convert(c config.StoreTestResults, inputs map[string]string) (*harness.Step, error) {
	name := namePrefix
	if c.Name != nil && *c.Name != "" {
		name = *c.Name
	}

	return &harness.Step{
		Name: name,
		Type: stepType,
		Spec: harness.StepExec{
			Run: "echo Storing test report",
			Reports: []*harness.Report{
				{
					Path: []string{utils.ResolveStrExpr(c.Path, inputs)},
					Type: "junit",
				},
			},
		},
	}, nil
}

package background

import (
	"github.com/drone/go-convert/convert/circle/converter/circleci/config"
	"github.com/drone/go-convert/convert/circle/converter/circleci/utils"

	harness "github.com/drone/spec/dist/go"
)

const (
	stepType   = "background"
	namePrefix = "background"
)

func Convert(c config.DockerExecutor, inParams map[string]string) (
	*harness.Step, error) {
	return &harness.Step{
		Name: utils.ResolveStrExpr(getName(c), inParams),
		Type: "background",
		Spec: harness.StepBackground{
			Envs:  utils.ResolveMapExpr(utils.ConvertEnvs_(c.Environment), inParams),
			Image: c.Image,
			// TODO (shubham): Set entrypoint after it is fixed to list in simplified yaml
			Entrypoint: utils.ResolveStrExpr(getEntrypoint(c.Entrypoint), inParams),
			Args:       utils.ResolveListExpr(getCmd(c.Command), inParams),
		},
	}, nil
}

func getName(c config.DockerExecutor) string {
	name := namePrefix
	if c.Name != nil && *c.Name != "" {
		name = *c.Name
	}
	return name
}

func getCmd(c *config.CommandUnion) []string {
	if c == nil {
		return nil
	}

	if c.String != nil {
		return []string{*c.String}
	}
	return c.StringArray
}

func getEntrypoint(c *config.CommandUnion) string {
	cmd := getCmd(c)
	if len(cmd) != 1 {
		return ""
	}
	return cmd[0]
}

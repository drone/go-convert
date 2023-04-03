package run

import (
	"github.com/drone/go-convert/convert/circle/converter/circleci/config"
	"github.com/drone/go-convert/convert/circle/converter/circleci/utils"

	harness "github.com/drone/spec/dist/go"
)

const (
	scriptType     = "script"
	backgroundType = "background"
	namePrefix     = "run"
)

func Convert(r config.Deploy, image string, inParams map[string]string) (
	*harness.Step, error) {
	if r.String != nil && *r.String != "" {
		return &harness.Step{
			Name: namePrefix,
			Type: scriptType,
			Spec: harness.StepExec{
				Run:   utils.ResolveStrExpr(*r.String, inParams),
				Image: image,
			},
		}, nil
	}

	if r.Confi == nil {
		return nil, nil
	}

	background := isBackground(r.Confi)
	envs := utils.ConvertEnvs(r.Confi.Environment)
	cmd := r.Confi.Command
	name := utils.ResolveStrExpr(getName(r.Confi), inParams)
	shell := getShell(r.Confi)

	if background {
		return &harness.Step{
			Name: name,
			Type: backgroundType,
			Spec: harness.StepBackground{
				Envs:  utils.ResolveMapExpr(envs, inParams),
				Run:   utils.ResolveStrExpr(cmd, inParams),
				Shell: utils.ResolveStrExpr(shell, inParams),
				Image: utils.ResolveStrExpr(image, inParams),
			},
		}, nil
	}
	return &harness.Step{
		Name: name,
		Type: scriptType,
		Spec: harness.StepExec{
			Envs:  utils.ResolveMapExpr(envs, inParams),
			Run:   utils.ResolveStrExpr(cmd, inParams),
			Shell: utils.ResolveStrExpr(shell, inParams),
			Image: utils.ResolveStrExpr(image, inParams),
		},
	}, nil
}

func getShell(c *config.Confi) string {
	shell := ""
	if c.Shell != nil && *c.Shell != "" {
		shell = *c.Shell
	}
	return shell
}

func getName(c *config.Confi) string {
	name := namePrefix
	if c.Name != nil && *c.Name != "" {
		name = *c.Name
	}
	return name
}

func isBackground(c *config.Confi) bool {
	background := false
	if c.Background != nil && *c.Background {
		background = true
	}
	return background
}

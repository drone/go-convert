package cache

import (
	"github.com/drone/go-convert/convert/circle/commons"
	"github.com/drone/go-convert/convert/circle/converter/circleci/config"
	"github.com/drone/go-convert/convert/circle/converter/circleci/utils"
	harness "github.com/drone/spec/dist/go"
)

const (
	saveStepType   = "plugin"
	saveNamePrefix = "save cache"
)

func ConvertSave(opts commons.Opts, c config.SaveCache,
	inputs map[string]string) (*harness.Step, error) {
	name := saveNamePrefix
	if c.Name != nil && *c.Name != "" {
		name = *c.Name
	}

	backend := getBackend(opts)

	m := make(map[string]interface{})
	m["bucket"] = getBucket(opts)
	m["cache_key"] = utils.ResolveStrExpr(getKey(&c.Key, []string{}), inputs)
	m["rebuild"] = "true"
	m["mount"] = utils.ResolveListExpr(c.Paths, inputs)
	m["exit_code"] = "true"
	m["archive_format"] = "tar"
	m["backend"] = backend
	m["backend_operation_timeout"] = "1800s"
	m["fail_restore_if_key_not_present"] = "false"
	if backend == "s3" {
		m["region"] = getRegion(opts)
		m["access_key"] = utils.ReplaceSecret(getS3AccessKey(opts), "access-key")
		m["secret_key"] = utils.ReplaceSecret(getS3SecretKey(opts), "secret-key")
	} else {
		m["json_key"] = utils.ReplaceSecret(getGCSJSONKey(opts), "json-key")
	}

	return &harness.Step{
		Name: name,
		Type: saveStepType,
		Spec: harness.StepPlugin{
			Image: "plugins/cache:1.4.6",
			With:  m,
		},
	}, nil
}

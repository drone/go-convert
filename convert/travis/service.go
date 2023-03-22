// Copyright 2022 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package travis

import harness "github.com/drone/spec/dist/go"

func (d *Converter) convertServices(ctx *context) []*harness.Step {

	// TODO support for addons.postgres
	// TODO support for addons.postgresql

	// TODO document authentication differences for postgres (password set to "postgres")
	// TODO document authentication differences for mysql (username set to "root", not "travis")
	var dst []*harness.Step
	for _, name := range ctx.config.Services {
		if _, ok := defaultServiceImage[name]; ok {
			dst = append(dst, d.convertService(name))
		}
	}
	if addons := ctx.config.Addons; addons != nil {
		if v := addons.Rethinkdb; v != "" {
			// TODO support addons.rethinkdb version
			dst = append(dst, d.convertService("rethinkdb"))
		}
		if v := addons.Mariadb; v != "" {
			// TODO support addons.mariadb version
			dst = append(dst, d.convertService("mariadb"))
		}
	}
	return dst
}

func (d *Converter) convertService(name string) *harness.Step {
	return &harness.Step{
		Name: d.identifiers.Generate(name),
		Type: "background",
		Spec: &harness.StepBackground{
			Image: defaultServiceImage[name],
			Ports: defaultServicePorts[name],
			Envs:  defaultServiceEnvs[name],
			// TODO support for adding an image connector to a background step
			// Connector: "",
		},
	}
}

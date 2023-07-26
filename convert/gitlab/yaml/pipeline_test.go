package yaml

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestPipelineYaml(t *testing.T) {
	tests, err := filepath.Glob("testdata/job_keywords/*/*.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			// parse the yaml file
			tmp1, err := ParseFile(test)
			if err != nil {
				t.Error(err)
				return
			}

			// marshal the yaml file
			tmp2, err := yaml.Marshal(tmp1)
			if err != nil {
				t.Error(err)
				return
			}

			// unmarshal the yaml file to a map
			got := map[string]interface{}{}
			if err := yaml.Unmarshal(tmp2, &got); err != nil {
				t.Error(err)
				return
			}

			// parse the golden yaml file and unmarshal
			data, err := ioutil.ReadFile(test + ".golden")
			if err != nil {
				// skip tests with no golden files
				// TODO these should be re-enabled
				return
			}

			// unmarshal the golden yaml file
			want := map[string]interface{}{}
			if err := yaml.Unmarshal(data, &want); err != nil {
				t.Error(err)
				return
			}

			// compare the parsed yaml to the golden file
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("Unexpected parsing result")
				t.Log(diff)
			}
		})
	}
}

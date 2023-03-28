package yaml

import (
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestEnvironment(t *testing.T) {
	tests := []struct {
		yaml string
		want Environment
	}{
		// string value
		{
			yaml: `staging_environment`,
			want: Environment{Name: "staging_environment"},
		},
		// struct value
		{
			yaml: `
  name: production_environment
  url: https://github.com
`,
			want: Environment{
				Name: "production_environment",
				URL:  "https://github.com",
			},
		},
	}

	for i, test := range tests {
		got := new(Environment)
		if err := yaml.Unmarshal([]byte(test.yaml), got); err != nil {
			t.Log(test.yaml)
			t.Error(err)
			return
		}
		if diff := cmp.Diff(got, &test.want); diff != "" {
			t.Log(test.yaml)
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}

func TestEnvironment_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Environment))
	if err == nil || err.Error() != "failed to unmarshal environment" {
		t.Errorf("Expect error, got %s", err)
	}
}

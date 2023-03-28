package yaml

import (
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestSecrets(t *testing.T) {
	tests := []struct {
		yaml string
		want Secrets
	}{
		// string value
		{
			yaml: `inherit`,
			want: Secrets{Inherit: true},
		},
		// struct value
		{
			yaml: `
 access-token: token
`,
			want: Secrets{
				Values: map[string]string{
					"access-token": "token",
				},
			},
		},
	}

	for i, test := range tests {
		got := new(Secrets)
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

func TestSecrets_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Secrets))
	if err == nil || err.Error() != "failed to unmarshal secrets" {
		t.Errorf("Expect error, got %s", err)
	}
}

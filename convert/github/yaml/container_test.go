package yaml

import (
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestContainer(t *testing.T) {
	tests := []struct {
		yaml string
		want Container
	}{
		// string value
		{
			yaml: `node:14.16`,
			want: Container{Image: "node:14.16"},
		},
		// struct value
		{
			yaml: `
      image: node:14.16
      env:
        NODE_ENV: development
      ports:
        - 80
      volumes:
        - my_docker_volume:/volume_mount
      options: --cpus 1
`,
			want: Container{
				Image: "node:14.16",
				Env: map[string]string{
					"NODE_ENV": "development",
				},
				Ports: []string{
					"80",
				},
				Volumes: []string{
					"my_docker_volume:/volume_mount",
				},
				Options: "--cpus 1",
			},
		},
		// struct value
		{
			yaml: `
  image: ghcr.io/owner/image
  credentials:
     username: username
     password: password
`,
			want: Container{
				Image: "ghcr.io/owner/image",
				Credentials: &Credentials{
					Username: "username",
					Password: "password",
				},
			},
		},
	}

	for i, test := range tests {
		got := new(Container)
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

func TestContainer_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Container))
	if err == nil || err.Error() != "failed to unmarshal container" {
		t.Errorf("Expect error, got %s", err)
	}
}

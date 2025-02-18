package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertKubeCtl(t *testing.T) {

	var tests []runner
	tests = append(tests, prepare(t, "/kubeCtl/kubeCtl_snippet", &harness.Step{
		Id:   "nulled0166",
		Name: "Apply Kubernetes Files",
		Type: "script",
		Spec: &harness.StepExec{
			Connector: "docker_hub_connector",
			Image:     "bitnami/kubectl",
			Shell:     "sh",
			Run:       "kubectl get pods",
			Envs: map[string]string{
				"KUBE_SERVER":         "https://127.0.0.1:56227",
				"KUBE_CONTEXT":        "context-name",
				"KUBE_CLUSTER":        "cluster-name",
				"KUBE_TOKEN":          "<+secrets.getValue(\"kube_token\")>",
				"KUBE_NAMESPACE":      "default",
				"KUBE_CA_CERTIFICATE": "",
				"KUBE_CREDENTIALS_ID": "credentialsId",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertKubeCtl(tt.input, tt.input.ParameterMap)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertNunit() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

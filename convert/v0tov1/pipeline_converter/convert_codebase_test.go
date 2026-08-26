package pipelineconverter

import (
	"encoding/json"
	"strings"
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

func TestConvertPipeline_CodebaseClonePreFetchAndSparseCheckout(t *testing.T) {
	converter := NewPipelineConverter()

	pipeline := &v0.Pipeline{
		ID:   "git_clone_advanced",
		Name: "git clone advanced",
		Props: v0.Properties{
			CI: v0.CI{
				Codebase: &v0.Codebase{
					Conn: "account.test",
					Build: flexible.Field[v0.Build]{
						Value: v0.Build{
							Type: "branch",
							Spec: v0.BuildSpec{Branch: "dev"},
						},
					},
					PreFetchCommand: "mkdir preFetchDir",
					SparseCheckout: &flexible.Field[[]string]{
						Value: []string{"folder2", "folder with space"},
					},
					Lfs:               &flexible.Field[bool]{Value: true},
					FetchTags:         &flexible.Field[bool]{Value: true},
					SubmoduleStrategy: &flexible.Field[bool]{Value: "recursive"},
				},
			},
		},
	}

	result := converter.ConvertPipeline(pipeline)
	if result == nil || result.Clone == nil {
		t.Fatal("expected pipeline clone config")
	}
	if result.Clone.PreFetchCommand != "mkdir preFetchDir" {
		t.Fatalf("pre-fetch command = %q", result.Clone.PreFetchCommand)
	}
	paths, ok := result.Clone.SparseCheckout.AsStruct()
	if !ok {
		t.Fatal("expected sparse-checkout list")
	}
	if len(paths) != 2 || paths[0] != "folder2" || paths[1] != "folder with space" {
		t.Fatalf("sparse-checkout = %#v", paths)
	}

	raw, err := json.Marshal(result.Clone)
	if err != nil {
		t.Fatalf("marshal clone: %v", err)
	}
	out := string(raw)
	if strings.Contains(out, "pre-fetch-command") {
		t.Fatalf("must emit pre-fetch, not pre-fetch-command: %s", out)
	}
	if strings.Contains(out, `"folder2,folder with space"`) {
		t.Fatalf("sparse-checkout must stay an array: %s", out)
	}
}

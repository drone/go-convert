package yaml

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/drone/go-convert/internal/flexible"
)

func TestCloneJSONMarshalPreFetchAndSparseCheckout(t *testing.T) {
	clone := Clone{
		Enabled:        true,
		PreFetchCommand: "mkdir preFetchDir",
		SparseCheckout: &flexible.Field[[]string]{
			Value: []string{"folder2", "folder with space"},
		},
	}

	raw, err := json.Marshal(clone)
	if err != nil {
		t.Fatalf("marshal clone: %v", err)
	}

	out := string(raw)
	if strings.Contains(out, "pre-fetch-command") {
		t.Fatalf("expected pre-fetch key, got legacy pre-fetch-command in %s", out)
	}
	if !strings.Contains(out, `"pre-fetch":"mkdir preFetchDir"`) {
		t.Fatalf("expected pre-fetch value in %s", out)
	}
	if strings.Contains(out, `"folder2,folder with space"`) {
		t.Fatalf("sparse-checkout must remain a list, not comma string: %s", out)
	}
	if !strings.Contains(out, `"folder2"`) || !strings.Contains(out, `"folder with space"`) {
		t.Fatalf("expected sparse-checkout array entries in %s", out)
	}
}

func TestCloneJSONUnmarshalPreFetchLegacyKey(t *testing.T) {
	raw := []byte(`{"enabled":true,"pre-fetch-command":"legacy command"}`)

	var clone Clone
	if err := json.Unmarshal(raw, &clone); err != nil {
		t.Fatalf("unmarshal clone: %v", err)
	}
	if clone.PreFetchCommand != "legacy command" {
		t.Fatalf("expected legacy pre-fetch-command to deserialize, got %q", clone.PreFetchCommand)
	}
}

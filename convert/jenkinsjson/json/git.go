package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

func ConvertClone(node Node) (*harness.CloneStage, *harness.Repository) {
	clone := &harness.CloneStage{
		Strategy: "http",
	}
	repo := &harness.Repository{
		Name:      node.SpanName,
		Connector: fmt.Sprintf("%v", node.ParameterMap["url"]),
	}
	return clone, repo
}

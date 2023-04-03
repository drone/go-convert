package converter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/drone/go-convert/convert/circle/commons"
	"github.com/drone/go-convert/convert/circle/converter/circleci"
	"gopkg.in/yaml.v3"
)

func Convert(opts commons.Opts, inFile, outFilePrefix string) error {
	data, err := os.ReadFile(inFile)
	if err != nil {
		return err
	}

	pipelines, err := circleci.Convert(opts, data)
	if err != nil {
		return err
	}

	for i, p := range pipelines {
		b, err := json.Marshal(p)
		if err != nil {
			return err
		}

		y, err := JSONToYAML(b)
		if err != nil {
			return err
		}

		fname := fmt.Sprintf("%s_%s_%d.yml", outFilePrefix, p.Name, i)
		if err := os.WriteFile(fname, y, 0644); err != nil {
			return err
		}
	}

	return nil
}

// JSONToYAML Converts JSON to YAML.
func JSONToYAML(j []byte) ([]byte, error) {
	// Convert the JSON to an object.
	var jsonObj interface{}
	// We are using yaml.Unmarshal here (instead of json.Unmarshal) because the
	// Go JSON library doesn't try to pick the right number type (int, float,
	// etc.) when unmarshalling to interface{}, it just picks float64
	// universally. go-yaml does go through the effort of picking the right
	// number type, so we can preserve number type throughout this process.
	err := yaml.Unmarshal(j, &jsonObj)
	if err != nil {
		return nil, err
	}
	// Marshal this object into YAML.
	return yaml.Marshal(jsonObj)
}

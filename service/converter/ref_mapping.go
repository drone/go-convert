package converter

import (
	"strings"

	pipelineconverter "github.com/drone/go-convert/convert/v0tov1/pipeline_converter"
	"gopkg.in/yaml.v3"
)

// ApplyRefMappings traverses a v1 YAML document and applies two independent
// rewrite maps to ref-bearing scalar values:
//
//   - templateRefs is applied to template references (the "ref" half of
//     "ref@version" in `template.uses` / `approval.uses`).
//   - pipelineRefs is applied to pipeline identifiers (`pipeline.id`, the
//     pipeline segment of `chain.uses` "org/project/pipeline", and
//     `pipelineIdentifier` on triggers).
//
// The trigger-embedded `inputYaml` scalar is itself a v1 YAML fragment; this
// function recurses into it so a single top-level call rewrites both the
// trigger wrapper and its nested pipeline inputs.
//
// When both maps are empty the input is returned unchanged with no parse.
func ApplyRefMappings(yamlBytes []byte, templateRefs, pipelineRefs map[string]string) ([]byte, error) {
	if len(templateRefs) == 0 && len(pipelineRefs) == 0 {
		return yamlBytes, nil
	}
	var root yaml.Node
	if err := yaml.Unmarshal(yamlBytes, &root); err != nil {
		return nil, err
	}
	applyRefsWalk(&root, "", "", templateRefs, pipelineRefs)
	return yaml.Marshal(&root)
}

// applyRefsWalk recursively walks the YAML tree. parentKey and
// grandparentKey carry the two enclosing mapping keys so rewrite rules
// can fire only in the right structural context (e.g. "inputs" → "overlay"
// → id).
func applyRefsWalk(node *yaml.Node, parentKey, grandparentKey string, templateRefs, pipelineRefs map[string]string) {
	if node == nil {
		return
	}
	switch node.Kind {
	case yaml.DocumentNode:
		for _, child := range node.Content {
			applyRefsWalk(child, parentKey, grandparentKey, templateRefs, pipelineRefs)
		}

	case yaml.MappingNode:
		for i := 0; i+1 < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]

			if keyNode.Kind != yaml.ScalarNode {
				applyRefsWalk(valueNode, parentKey, grandparentKey, templateRefs, pipelineRefs)
				continue
			}
			key := keyNode.Value

			if valueNode.Kind == yaml.ScalarNode {
				switch {
				case key == "uses" && parentKey == "template":
					valueNode.Value = remapTemplateUses(valueNode.Value, templateRefs)

				case key == "uses" && parentKey == "chain":
					valueNode.Value = remapChainUses(valueNode.Value, pipelineRefs)

				case key == "id" && (parentKey == "pipeline" || (parentKey == "overlay" && grandparentKey == "inputs")):
					if n, ok := pipelineRefs[valueNode.Value]; ok {
						valueNode.Value = n
					}

				case key == "pipelineIdentifier":
					if n, ok := pipelineRefs[valueNode.Value]; ok {
						valueNode.Value = n
					}

				case key == "inputYaml" && strings.TrimSpace(valueNode.Value) != "":
					inner, err := ApplyRefMappings([]byte(valueNode.Value), templateRefs, pipelineRefs)
					if err != nil {
						pipelineconverter.GetMessageLogger().LogWarning(
							"INPUT_YAML_REF_MAPPING_FAILED",
							"failed to apply ref mappings to embedded inputYaml",
							pipelineconverter.WithContext(map[string]string{"error": err.Error()}),
						)
					} else {
						valueNode.Value = string(inner)
					}
				}
			}

			applyRefsWalk(valueNode, key, parentKey, templateRefs, pipelineRefs)
		}

	case yaml.SequenceNode:
		for _, child := range node.Content {
			applyRefsWalk(child, parentKey, grandparentKey, templateRefs, pipelineRefs)
		}
	}
}

// remapTemplateUses rewrites a `uses:` value of the form "ref" or
// "ref@version". Only the ref portion is looked up in the mapping; the
// "@version" suffix (if present) is preserved verbatim.
func remapTemplateUses(value string, templateRefs map[string]string) string {
	if len(templateRefs) == 0 || value == "" {
		return value
	}
	ref := value
	suffix := ""
	if idx := strings.Index(value, "@"); idx >= 0 {
		ref = value[:idx]
		suffix = value[idx:]
	}
	if n, ok := templateRefs[ref]; ok {
		return n + suffix
	}
	return value
}

// remapChainUses rewrites a chain `uses:` value of the form
// "org/project/pipeline". A full-value match is tried first (so callers can
// remap by fully-qualified ref); if that fails, the third segment (pipeline
// identifier) is looked up independently.
func remapChainUses(value string, pipelineRefs map[string]string) string {
	if len(pipelineRefs) == 0 || value == "" {
		return value
	}
	if n, ok := pipelineRefs[value]; ok {
		return n
	}
	parts := strings.Split(value, "/")
	if len(parts) == 3 {
		if n, ok := pipelineRefs[parts[2]]; ok {
			parts[2] = n
			return strings.Join(parts, "/")
		}
	}
	return value
}

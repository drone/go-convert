package converter

import (
	"strings"

	pipelineconverter "github.com/drone/go-convert/convert/v0tov1/pipeline_converter"
	"gopkg.in/yaml.v3"
)

// ApplyRefMappings traverses a v1 YAML document and applies four independent
// rewrite maps to ref-bearing scalar values:
//
//   - templateRefs is applied to template references (the "ref" half of
//     "ref@version" in `template.uses` / `approval.uses`).
//   - pipelineRefs is applied to pipeline identifiers (`pipeline.id`, the
//     pipeline segment of `chain.uses` "org/project/pipeline", and
//     `pipelineIdentifier` on triggers).
//   - templateRefBranches and pipelineRefBranches are keyed like templateRefs /
//     pipelineRefs, but map to a destination branch instead of a new
//     identifier. When an identifier is rewritten via `uses` AND has an
//     entry in the matching branch map, the sibling `ref:` key on the
//     template or chain block is overwritten with that branch; otherwise
//     it's left untouched.
//
// The trigger-embedded `inputYaml` scalar is itself a v1 YAML fragment; this
// function recurses into it so a single top-level call rewrites both the
// trigger wrapper and its nested pipeline inputs.
//
// When all four maps are empty the input is returned unchanged with no parse.
func ApplyRefMappings(yamlBytes []byte, templateRefs, pipelineRefs, templateRefBranches, pipelineRefBranches map[string]string) ([]byte, error) {
	if len(templateRefs) == 0 && len(pipelineRefs) == 0 && len(templateRefBranches) == 0 && len(pipelineRefBranches) == 0 {
		return yamlBytes, nil
	}
	var root yaml.Node
	if err := yaml.Unmarshal(yamlBytes, &root); err != nil {
		return nil, err
	}
	applyRefsWalk(&root, "", "", templateRefs, pipelineRefs, templateRefBranches, pipelineRefBranches)
	return yaml.Marshal(&root)
}

// applyRefsWalk recursively walks the YAML tree. parentKey and
// grandparentKey carry the two enclosing mapping keys so rewrite rules
// can fire only in the right structural context (e.g. "inputs" → "overlay"
// → id).
func applyRefsWalk(node *yaml.Node, parentKey, grandparentKey string, templateRefs, pipelineRefs, templateRefBranches, pipelineRefBranches map[string]string) {
	if node == nil {
		return
	}
	switch node.Kind {
	case yaml.DocumentNode:
		for _, child := range node.Content {
			applyRefsWalk(child, parentKey, grandparentKey, templateRefs, pipelineRefs, templateRefBranches, pipelineRefBranches)
		}

	case yaml.MappingNode:
		// Resolve the destination branch (if any) for this node's "uses"
		// entry up front, so the ref rewrite below doesn't depend on key
		// order within the node.
		resolvedBranch := resolveBranchForUses(node, parentKey, templateRefs, pipelineRefs, templateRefBranches, pipelineRefBranches)
		branchKeySeen := false

		for i := 0; i+1 < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]

			if keyNode.Kind != yaml.ScalarNode {
				applyRefsWalk(valueNode, parentKey, grandparentKey, templateRefs, pipelineRefs, templateRefBranches, pipelineRefBranches)
				continue
			}
			key := keyNode.Value

			if valueNode.Kind == yaml.ScalarNode {
				switch {
				case key == "uses" && parentKey == "template":
					valueNode.Value = remapTemplateUses(valueNode.Value, templateRefs)

				case key == "uses" && parentKey == "chain":
					valueNode.Value = remapChainUses(valueNode.Value, pipelineRefs)

				case key == "ref" && (parentKey == "template" || parentKey == "chain"):
					branchKeySeen = true
					if resolvedBranch != "" {
						valueNode.Value = resolvedBranch
					}

				case key == "id" && (parentKey == "pipeline" || (parentKey == "overlay" && grandparentKey == "inputs")):
					if n, ok := pipelineRefs[valueNode.Value]; ok {
						valueNode.Value = n
					}

				case key == "id" && parentKey == "template":
					if n, ok := templateRefs[valueNode.Value]; ok {
						valueNode.Value = n
					}

				case key == "pipelineIdentifier" && parentKey == "trigger":
					if n, ok := pipelineRefs[valueNode.Value]; ok {
						valueNode.Value = n
					}

				case key == "inputYaml" && strings.TrimSpace(valueNode.Value) != "":
					inner, err := ApplyRefMappings([]byte(valueNode.Value), templateRefs, pipelineRefs, templateRefBranches, pipelineRefBranches)
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

			applyRefsWalk(valueNode, key, parentKey, templateRefs, pipelineRefs, templateRefBranches, pipelineRefBranches)
		}

		// The source step may have had no branch pin at all (ambient, same-branch case), so the
		// ref sibling key never existed for the loop above to rewrite. Append it now so a
		// resolved branch override still lands even when go-convert's base conversion omitted the key.
		if resolvedBranch != "" && !branchKeySeen && (parentKey == "template" || parentKey == "chain") {
			node.Content = append(node.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "ref"},
				&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: resolvedBranch},
			)
		}

	case yaml.SequenceNode:
		for _, child := range node.Content {
			applyRefsWalk(child, parentKey, grandparentKey, templateRefs, pipelineRefs, templateRefBranches, pipelineRefBranches)
		}
	}
}

// resolveBranchForUses returns the destination branch for this node's "uses"
// entry (template or chain), or "" if there's no "uses" entry, no matching
// id-mapping entry, or no branch-mapping entry.
func resolveBranchForUses(node *yaml.Node, parentKey string, templateRefs, pipelineRefs, templateRefBranches, pipelineRefBranches map[string]string) string {
	if parentKey != "template" && parentKey != "chain" {
		return ""
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]
		if keyNode.Kind != yaml.ScalarNode || keyNode.Value != "uses" || valueNode.Kind != yaml.ScalarNode {
			continue
		}

		switch parentKey {
		case "template":
			ref := valueNode.Value
			if idx := strings.Index(ref, "@"); idx >= 0 {
				ref = ref[:idx]
			}
			if _, ok := templateRefs[ref]; ok {
				return templateRefBranches[ref]
			}
		case "chain":
			parts := strings.Split(valueNode.Value, "/")
			if len(parts) == 3 {
				if _, ok := pipelineRefs[parts[2]]; ok {
					return pipelineRefBranches[parts[2]]
				}
			}
		}
		// Found the "uses" entry but no matching id-mapping key: stop
		// looking (there's only one "uses" key per mapping node).
		return ""
	}
	return ""
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
// identifier) is looked up independently. In both cases the mapping value is
// used verbatim as the fully-qualified replacement (it already carries the
// org/project prefix).
func remapChainUses(value string, pipelineRefs map[string]string) string {
	if len(pipelineRefs) == 0 || value == "" {
		return value
	}
	parts := strings.Split(value, "/")
	if len(parts) == 3 {
		if ref, ok := pipelineRefs[parts[2]]; ok {
			// The mapping value is a fully-qualified ref (already carrying the
			// org/project prefix), so replace the whole value rather than just
			// the pipeline segment.
			return ref
		}
	}
	return value
}

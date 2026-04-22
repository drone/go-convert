package converter

import (
	"gopkg.in/yaml.v3"
)

// ReplaceTemplateRefs traverses the YAML and replaces template reference identifiers
// based on the provided mapping. It looks for "templateRef" and "template_ref" fields
// and replaces their values if they exist in the mapping.
func ReplaceTemplateRefs(yamlBytes []byte, mapping map[string]string) ([]byte, error) {
	if len(mapping) == 0 {
		return yamlBytes, nil // no-op if no mapping provided
	}

	var root yaml.Node
	if err := yaml.Unmarshal(yamlBytes, &root); err != nil {
		return nil, err
	}

	// Traverse and replace
	traverseAndReplace(&root, mapping)

	return yaml.Marshal(&root)
}

// traverseAndReplace recursively walks the YAML tree and replaces templateRef values
func traverseAndReplace(node *yaml.Node, mapping map[string]string) {
	if node == nil {
		return
	}

	switch node.Kind {
	case yaml.DocumentNode:
		for _, child := range node.Content {
			traverseAndReplace(child, mapping)
		}

	case yaml.MappingNode:
		// Process key-value pairs
		for i := 0; i < len(node.Content)-1; i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]

			// Check if this is a templateRef field
			if keyNode.Kind == yaml.ScalarNode &&
				(keyNode.Value == "templateRef" || keyNode.Value == "template_ref") {

				// Replace the value if it exists in mapping
				if valueNode.Kind == yaml.ScalarNode {
					if newRef, exists := mapping[valueNode.Value]; exists {
						valueNode.Value = newRef
					}
				}
			}

			// Recurse into nested structures
			traverseAndReplace(valueNode, mapping)
		}

	case yaml.SequenceNode:
		for _, child := range node.Content {
			traverseAndReplace(child, mapping)
		}
	}
}

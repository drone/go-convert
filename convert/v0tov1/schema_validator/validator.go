// Copyright 2022 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schema_validator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
)

// Entity type constants used to select the correct schema.
const (
	EntityPipeline = "pipeline"
	EntityTemplate = "template"
	EntityTrigger  = "trigger"
	EntityInputSet = "inputSet"
)

// schemaFiles maps entity types to their JSON schema file names.
var schemaFiles = map[string]string{
	EntityPipeline: "pipeline.json",
	EntityTemplate: "template.json",
	EntityTrigger:  "trigger.json",
	EntityInputSet: "inputSet.json",
}

// lenientRegexpEngine wraps Go's regexp.Compile but falls back to a
// match-everything regex when the pattern uses Perl-specific syntax
// (e.g. negative lookaheads) that Go's regexp package does not support.
// This is required because the Harness v1 JSON schemas include PCRE
// patterns that are valid in Java/Python but not in Go.
func lenientRegexpEngine(pattern string) (jsonschema.Regexp, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		// Return a regex that matches everything — effectively
		// skipping pattern validation for this field.
		re = regexp.MustCompile(`[\s\S]*`)
	}
	return re, nil
}

// SchemaValidator loads and caches compiled JSON Schemas from a directory,
// then validates converted v1 YAML documents against them.
type SchemaValidator struct {
	schemaDir string
	schemas   map[string]*jsonschema.Schema // entity type → compiled schema
}

// NewSchemaValidator creates a validator by loading and compiling all
// available v1 JSON Schema files from schemaDir. Missing schema files
// produce a log warning rather than a hard error so that partial schema
// directories still work.
func NewSchemaValidator(schemaDir string) (*SchemaValidator, error) {
	info, err := os.Stat(schemaDir)
	if err != nil {
		return nil, fmt.Errorf("schema directory not accessible: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("schema path is not a directory: %s", schemaDir)
	}

	v := &SchemaValidator{
		schemaDir: schemaDir,
		schemas:   make(map[string]*jsonschema.Schema),
	}

	for entityType, fileName := range schemaFiles {
		schemaPath := filepath.Join(schemaDir, fileName)
		if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
			log.Printf("Warning: schema file not found for %s: %s", entityType, schemaPath)
			continue
		}
		c := jsonschema.NewCompiler()
		c.UseRegexpEngine(lenientRegexpEngine)
		sch, err := c.Compile(schemaPath)
		if err != nil {
			log.Printf("Warning: failed to compile schema for %s: %v", entityType, err)
			continue
		}
		v.schemas[entityType] = sch
	}

	if len(v.schemas) == 0 {
		return nil, fmt.Errorf("no schemas could be loaded from %s", schemaDir)
	}
	return v, nil
}

// HasSchema returns true if a schema is loaded for the given entity type.
func (v *SchemaValidator) HasSchema(entityType string) bool {
	_, ok := v.schemas[entityType]
	return ok
}

// Validate parses yamlBytes as YAML, converts to a JSON-compatible
// interface{}, and validates against the schema for entityType.
func (v *SchemaValidator) Validate(yamlBytes []byte, entityType string) *ValidationResult {
	result := &ValidationResult{EntityType: entityType}

	sch, ok := v.schemas[entityType]
	if !ok {
		result.Valid = false
		result.ErrorMessage = fmt.Sprintf("no schema loaded for entity type %q", entityType)
		return result
	}

	// YAML → JSON bytes → interface{}
	jsonBytes, err := yaml.YAMLToJSON(yamlBytes)
	if err != nil {
		result.Valid = false
		result.ErrorMessage = fmt.Sprintf("invalid YAML: %v", err)
		return result
	}

	var doc interface{}
	if err := json.Unmarshal(jsonBytes, &doc); err != nil {
		result.Valid = false
		result.ErrorMessage = fmt.Sprintf("failed to unmarshal to JSON: %v", err)
		return result
	}

	// Run validation
	err = sch.Validate(doc)
	if err == nil {
		result.Valid = true
		return result
	}

	ve, ok := err.(*jsonschema.ValidationError)
	if !ok {
		result.Valid = false
		result.ErrorMessage = err.Error()
		return result
	}

	// Process errors using the BasicOutput (flat list).
	basicOut := ve.BasicOutput()
	leafErrors := collectLeafErrors(basicOut)

	// Filter and deduplicate — ported from validate_v1.py.
	filtered := filterErrors(leafErrors)

	if len(filtered) == 0 {
		result.Valid = true
		return result
	}

	// Build DTOs with stage/step context.
	dtos := make([]SchemaErrorDTO, 0, len(filtered))
	for _, oe := range filtered {
		fqn := instanceLocToFQN(oe.InstanceLocation)
		parts := pathParts(fqn)
		stageInfo := extractNodeInfo(doc, parts, "stage")
		stepInfo := extractNodeInfo(doc, parts, "step")

		errMsg := describeError(oe, doc)

		dtos = append(dtos, SchemaErrorDTO{
			Message:        errMsg,
			MessageWithFQN: fqn + ": " + errMsg,
			FQN:            fqn,
			StageInfo:      stageInfo,
			StepInfo:       stepInfo,
		})
	}

	msgs := make([]string, len(dtos))
	for i, d := range dtos {
		msgs[i] = d.MessageWithFQN
	}

	result.Valid = false
	result.ErrorMessage = strings.Join(msgs, "; ")
	result.SchemaErrors = dtos
	return result
}

// DetectEntityType inspects the raw YAML bytes and returns the entity type.
func DetectEntityType(yamlBytes []byte) string {
	var raw map[string]interface{}
	if err := yaml.Unmarshal(yamlBytes, &raw); err != nil {
		return EntityPipeline // default
	}
	if _, ok := raw["template"]; ok {
		return EntityTemplate
	}
	if k, ok := raw["kind"]; ok {
		switch fmt.Sprintf("%v", k) {
		case "trigger":
			return EntityTrigger
		case "input-set":
			return EntityInputSet
		}
	}
	if _, ok := raw["spec"]; ok {
		if _, hasPipeline := raw["pipeline"]; !hasPipeline {
			// Has spec but no pipeline key — could be trigger or inputSet.
			// Check for version+kind pattern.
			if k, ok := raw["kind"]; ok {
				switch fmt.Sprintf("%v", k) {
				case "trigger":
					return EntityTrigger
				case "input-set":
					return EntityInputSet
				}
			}
		}
	}
	return EntityPipeline
}

// ---------------------------------------------------------------------------
// Error processing helpers — ported from harness-schema validate_v1.py
// ---------------------------------------------------------------------------

// collectLeafErrors extracts all leaf-level OutputUnit entries from
// the BasicOutput flat list (those with a non-nil Error).
func collectLeafErrors(out *jsonschema.OutputUnit) []jsonschema.OutputUnit {
	var leaves []jsonschema.OutputUnit
	if out.Error != nil {
		leaves = append(leaves, *out)
	}
	for _, child := range out.Errors {
		leaves = append(leaves, collectLeafErrors(&child)...)
	}
	return leaves
}

// filterErrors implements the Harness-style error filtering:
//  1. Drop additionalProperties errors.
//  2. Deduplicate: keep deepest error per instance location.
//  3. Drop anyOf/oneOf wrapper errors when children exist at same path.
func filterErrors(errors []jsonschema.OutputUnit) []jsonschema.OutputUnit {
	// 1. Filter out additionalProperties
	var step1 []jsonschema.OutputUnit
	for _, e := range errors {
		if isAdditionalPropertiesError(e) {
			continue
		}
		step1 = append(step1, e)
	}

	// 2. Deduplicate by instance location — keep deepest
	seen := make(map[string]int) // instanceLocation → index in deduped
	deduped := make([]jsonschema.OutputUnit, 0, len(step1))
	for _, e := range step1 {
		loc := e.InstanceLocation
		if idx, exists := seen[loc]; exists {
			// Keep the one with deeper keyword location
			existing := deduped[idx]
			if len(e.KeywordLocation) > len(existing.KeywordLocation) {
				deduped[idx] = e
			}
		} else {
			seen[loc] = len(deduped)
			deduped = append(deduped, e)
		}
	}

	// 3. Drop anyOf/oneOf/if wrappers when more specific errors exist
	childLocs := make(map[string]bool)
	for _, e := range deduped {
		childLocs[e.InstanceLocation] = true
	}

	var final []jsonschema.OutputUnit
	for _, e := range deduped {
		if isCompositionError(e) {
			// Check if a non-composition error shares this instance location
			hasSpecific := false
			for _, other := range deduped {
				if other.InstanceLocation == e.InstanceLocation && !isCompositionError(other) {
					hasSpecific = true
					break
				}
			}
			if hasSpecific {
				continue
			}
		}
		final = append(final, e)
	}
	return final
}

// describeError produces a human-readable error message from an OutputUnit.
// When the library returns a generic "validation failed" (kind.Group), it
// enriches the message with the failing schema keyword extracted from
// KeywordLocation and the actual value at the instance location.
func describeError(oe jsonschema.OutputUnit, doc interface{}) string {
	if oe.Error == nil {
		return "validation failed"
	}

	// Try concrete error kind first — these already produce good messages.
	switch k := oe.Error.Kind.(type) {
	case *kind.Type:
		want := strings.Join(k.Want, " or ")
		return fmt.Sprintf("got %s, want %s", k.Got, want)
	case *kind.Required:
		if len(k.Missing) == 1 {
			return fmt.Sprintf("missing required property '%s'", k.Missing[0])
		}
		return fmt.Sprintf("missing required properties: %s", strings.Join(k.Missing, ", "))
	case *kind.Enum:
		var vals []string
		for _, v := range k.Want {
			vals = append(vals, fmt.Sprintf("%v", v))
		}
		return fmt.Sprintf("value must be one of [%s]", strings.Join(vals, ", "))
	case *kind.Const:
		return fmt.Sprintf("value must be %v", k.Want)
	case *kind.Pattern:
		return fmt.Sprintf("'%s' does not match pattern '%s'", k.Got, k.Want)
	case *kind.MinLength:
		return fmt.Sprintf("string too short: length %d < minimum %d", k.Got, k.Want)
	case *kind.MaxLength:
		return fmt.Sprintf("string too long: length %d > maximum %d", k.Got, k.Want)
	case *kind.Minimum:
		got, _ := k.Got.Float64()
		want, _ := k.Want.Float64()
		return fmt.Sprintf("value %v < minimum %v", got, want)
	case *kind.Maximum:
		got, _ := k.Got.Float64()
		want, _ := k.Want.Float64()
		return fmt.Sprintf("value %v > maximum %v", got, want)
	case *kind.MinItems:
		return fmt.Sprintf("array too short: %d items < minimum %d", k.Got, k.Want)
	case *kind.MaxItems:
		return fmt.Sprintf("array too long: %d items > maximum %d", k.Got, k.Want)
	case *kind.MinProperties:
		return fmt.Sprintf("too few properties: %d < minimum %d", k.Got, k.Want)
	case *kind.MaxProperties:
		return fmt.Sprintf("too many properties: %d > maximum %d", k.Got, k.Want)
	case *kind.FalseSchema:
		return describeFromKeywordLocation(oe.KeywordLocation, oe.InstanceLocation, doc)
	case *kind.AnyOf:
		return describeFromKeywordLocation(oe.KeywordLocation, oe.InstanceLocation, doc)
	case *kind.OneOf:
		if len(k.Subschemas) == 0 {
			return describeFromKeywordLocation(oe.KeywordLocation, oe.InstanceLocation, doc)
		}
		return oe.Error.String()
	case *kind.AllOf:
		return describeFromKeywordLocation(oe.KeywordLocation, oe.InstanceLocation, doc)
	case *kind.Not:
		return describeFromKeywordLocation(oe.KeywordLocation, oe.InstanceLocation, doc)
	case *kind.AdditionalProperties:
		return fmt.Sprintf("additional properties not allowed: %s", strings.Join(k.Properties, ", "))
	case *kind.UniqueItems:
		return fmt.Sprintf("array items at index %d and %d are equal", k.Duplicates[0], k.Duplicates[1])
	case *kind.Group:
		return describeFromKeywordLocation(oe.KeywordLocation, oe.InstanceLocation, doc)
	}

	// Fallback: use the library's own message if it's not just "validation failed"
	msg := oe.Error.String()
	if msg != "validation failed" {
		return msg
	}
	return describeFromKeywordLocation(oe.KeywordLocation, oe.InstanceLocation, doc)
}

// describeFromKeywordLocation builds a descriptive error from the schema
// keyword path and the document value at the failing instance location.
func describeFromKeywordLocation(kwLoc, instLoc string, doc interface{}) string {
	// Extract the failing schema keyword(s) from KeywordLocation.
	// e.g. "/properties/notifications/items/properties/uses/type" → "type"
	// e.g. "/properties/stages/items/anyOf" → "anyOf"
	rule := extractSchemaRule(kwLoc)

	// Try to describe the actual value at the instance location.
	val := resolveJSONPointer(doc, instLoc)
	valDesc := describeValue(val)

	if rule != "" && valDesc != "" {
		return fmt.Sprintf("does not match '%s' rule (%s)", rule, valDesc)
	}
	if rule != "" {
		return fmt.Sprintf("does not match '%s' rule in schema", rule)
	}
	if valDesc != "" {
		return fmt.Sprintf("validation failed (%s)", valDesc)
	}
	return "validation failed"
}

// extractSchemaRule returns the innermost meaningful schema keyword from a
// KeywordLocation path. It skips structural segments like "properties",
// "items", numeric indices, "allOf", "anyOf", "oneOf".
func extractSchemaRule(kwLoc string) string {
	parts := strings.Split(strings.TrimPrefix(kwLoc, "/"), "/")
	if len(parts) == 0 {
		return ""
	}

	// Walk from the end to find the most specific keyword.
	structural := map[string]bool{
		"properties": true, "items": true, "patternProperties": true,
		"$ref": true, "$defs": true, "definitions": true,
	}
	for i := len(parts) - 1; i >= 0; i-- {
		p := parts[i]
		if p == "" {
			continue
		}
		// Skip numeric indices (used in allOf/0, anyOf/1, etc.)
		if _, err := strconv.Atoi(p); err == nil {
			continue
		}
		if structural[p] {
			continue
		}
		// Return the last meaningful keyword with context if possible.
		// e.g. for "properties/uses/type" return "type"
		// e.g. for "properties/stages/items" return the property name before it
		return p
	}
	return ""
}

// resolveJSONPointer follows a JSON Pointer path into a document.
func resolveJSONPointer(doc interface{}, pointer string) interface{} {
	if pointer == "" || pointer == "/" {
		return doc
	}
	parts := strings.Split(strings.TrimPrefix(pointer, "/"), "/")
	current := doc
	for _, p := range parts {
		switch node := current.(type) {
		case map[string]interface{}:
			var ok bool
			current, ok = node[p]
			if !ok {
				return nil
			}
		case []interface{}:
			idx, err := strconv.Atoi(p)
			if err != nil || idx < 0 || idx >= len(node) {
				return nil
			}
			current = node[idx]
		default:
			return nil
		}
	}
	return current
}

// describeValue returns a short description of a JSON value for error context.
func describeValue(val interface{}) string {
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		if len(keys) > 5 {
			keys = append(keys[:5], "...")
		}
		return fmt.Sprintf("got object with keys: [%s]", strings.Join(keys, ", "))
	case []interface{}:
		return fmt.Sprintf("got array with %d items", len(v))
	case string:
		if len(v) > 50 {
			return fmt.Sprintf("got string '%s...'", v[:50])
		}
		return fmt.Sprintf("got string '%s'", v)
	case bool:
		return fmt.Sprintf("got boolean %v", v)
	case float64:
		return fmt.Sprintf("got number %v", v)
	default:
		return fmt.Sprintf("got %T", v)
	}
}

// isAdditionalPropertiesError checks if the error is about additional properties.
func isAdditionalPropertiesError(e jsonschema.OutputUnit) bool {
	return strings.Contains(e.KeywordLocation, "/additionalProperties")
}

// isCompositionError checks if the error is an anyOf/oneOf/if wrapper.
func isCompositionError(e jsonschema.OutputUnit) bool {
	kw := e.KeywordLocation
	return strings.HasSuffix(kw, "/anyOf") ||
		strings.HasSuffix(kw, "/oneOf") ||
		strings.HasSuffix(kw, "/if")
}

// instanceLocToFQN converts a JSON Pointer instance location (e.g.
// "/pipeline/stages/0/steps/1") to a dotted FQN with array indices
// (e.g. "$.pipeline.stages[0].steps[1]").
func instanceLocToFQN(loc string) string {
	if loc == "" {
		return "$"
	}
	parts := strings.Split(strings.TrimPrefix(loc, "/"), "/")
	var sb strings.Builder
	sb.WriteString("$")
	for i := 0; i < len(parts); i++ {
		p := parts[i]
		if _, err := strconv.Atoi(p); err == nil {
			// Array index — attach to previous segment.
			sb.WriteString("[" + p + "]")
		} else {
			sb.WriteString("." + p)
		}
	}
	return sb.String()
}

// pathParts converts "$.pipeline.stages[0].steps[1]" →
// ["pipeline", "stages", "0", "steps", "1"].
var indexRe = regexp.MustCompile(`^([^\[]+)(?:\[(\d+)\])?$`)

func pathParts(fqn string) []string {
	var parts []string
	for _, seg := range strings.Split(strings.TrimPrefix(fqn, "$."), ".") {
		m := indexRe.FindStringSubmatch(seg)
		if m != nil {
			parts = append(parts, m[1])
			if m[2] != "" {
				parts = append(parts, m[2])
			}
		} else {
			parts = append(parts, seg)
		}
	}
	return parts
}

// resolveJSONPath walks a parsed doc by path components.
func resolveJSONPath(doc interface{}, parts []string) interface{} {
	cur := doc
	for _, p := range parts {
		if cur == nil {
			return nil
		}
		switch v := cur.(type) {
		case map[string]interface{}:
			cur = v[p]
		case []interface{}:
			idx, err := strconv.Atoi(p)
			if err != nil || idx < 0 || idx >= len(v) {
				return nil
			}
			cur = v[idx]
		default:
			return nil
		}
	}
	return cur
}

// extractNodeInfo walks up from pathParts looking for the nearest
// stage or step ancestor (mirrors the Python _extract_node_info).
func extractNodeInfo(doc interface{}, parts []string, nodeType string) *NodeErrorInfo {
	parentArrayKey := "stages"
	if nodeType == "step" {
		parentArrayKey = "steps"
	}
	for i := len(parts); i >= 2; i-- {
		if parts[i-2] == parentArrayKey {
			node := resolveJSONPath(doc, parts[:i])
			if m, ok := node.(map[string]interface{}); ok {
				info := &NodeErrorInfo{
					FQN: "$." + strings.Join(parts[:i], "."),
				}
				if id, ok := m["id"].(string); ok {
					info.Identifier = id
				}
				if t, ok := m["type"].(string); ok {
					info.Type = t
				}
				if n, ok := m["name"].(string); ok {
					info.Name = n
				}
				return info
			}
		}
	}
	return nil
}

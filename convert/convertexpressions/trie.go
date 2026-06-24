package convertexpressions

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Precompiled regexes used during trie operations.
var (
	metadataRegex      = regexp.MustCompile(`^([^\{]+)\{([^\}]+)\}$`)
	arrayNotationRegex = regexp.MustCompile(`^([^\[]+)(\[[^\]]+\])$`)
	// Matches parenthesized path with optional array index and child field:
	// Format: (parent.path)[i].childField or (parent.path).childField
	// Group 1: parent path without parens (e.g., "spec.output")
	// Group 2: optional array index (e.g., "[i]" or empty)
	// Group 3: child field after the dot (e.g., "name")
	parenPathRegex = regexp.MustCompile(`^\(([^)]+)\)(\[[^\]]*\])?\.(.+)$`)
	// Matches a standalone parenthesized path (entire pattern is parenthesized)
	// Format: (a.b.c) - collapses multiple segments into one
	standaloneParenRegex = regexp.MustCompile(`^\(([^)]+)\)$`)
)

// TrieNode represents a node in the trie structure
type TrieNode struct {
	children      map[string]*TrieNode // Named children (exact matches)
	wildcardChild *TrieNode            // Single wildcard child for dynamic segments

	// Context-aware children: map of context key (e.g. step type) to sub-trie root
	// These are checked when a ConversionContext is provided during matching
	contextChildren map[string]*TrieNode

	// Node metadata
	alias      string // Alias for relative path entry (e.g., "stage", "step")
	id         string // Unique identifier for this node
	v1Name     string // v1 output name; "-" means skip in output, "" means use key
	isWildcard bool   // True for dynamic IDs (STAGE_ID, STEP_ID)
	isArray    bool   // True if this node represents an array

	// noFQNOverride marks an alias entry node where FQN stage/chain expansion
	// is disabled (the remainder converts structurally). Set via
	// WithNoFQNOverride() on the stage node, its ancestors (pipeline./stages.),
	// and the execution roots, which already carry stage context.
	noFQNOverride bool

	// Array parent path support for rules like: outputVariables[i].name -> (spec.output[i]).alias
	// When set, this indicates the array node maps to a multi-segment v1 path
	arrayParentV1Path string // e.g., "spec.output" for the array node itself

	// Terminal node data
	isEnd  bool   // True if this is a complete conversion rule endpoint
	target string // Replacement pattern if this is an end node
}

// Trie represents the trie structure for conversion rules
type Trie struct {
	root       *TrieNode
	nodeIndex  map[string]*TrieNode   // Quick lookup by node ID
	aliasIndex map[string][]*TrieNode // Quick lookup by alias (single node per alias)
}

// NewTrie creates a new trie instance
func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			children: make(map[string]*TrieNode),
		},
		nodeIndex:  make(map[string]*TrieNode),
		aliasIndex: make(map[string][]*TrieNode),
	}
}

// AttachRulesAt attaches conversion rules to a node identified by ID
// Rules are relative to the node
func (t *Trie) AttachRulesAt(nodeID string, rules []ConversionRule) {
	node := t.nodeIndex[nodeID]
	if node == nil {
		fmt.Printf("Warning: node ID '%s' not found in trie\n", nodeID)
		return
	}

	for _, rule := range rules {
		t.insertFromNode(node, rule.From, rule.To)
	}
}

// AttachRulesWithContextAt attaches context-aware conversion rules to a node.
// These rules are only matched when the ConversionContext has a matching StepType.
// contextKey is typically the step type (e.g., "Run", "Http").
func (t *Trie) AttachRulesWithContextAt(nodeID string, contextKey string, rules []ConversionRule) {
	node := t.nodeIndex[nodeID]
	if node == nil {
		fmt.Printf("Warning: node ID '%s' not found in trie\n", nodeID)
		return
	}

	if node.contextChildren == nil {
		node.contextChildren = make(map[string]*TrieNode)
	}

	// Get or create context-specific sub-trie root
	contextRoot := node.contextChildren[contextKey]
	if contextRoot == nil {
		contextRoot = &TrieNode{
			children: make(map[string]*TrieNode),
		}
		node.contextChildren[contextKey] = contextRoot
	}

	for _, rule := range rules {
		t.insertFromNode(contextRoot, rule.From, rule.To)
	}
}

// parseNodeMetadata extracts node name and metadata from a pattern part
// Format: "nodeName{alias: env, id: node_id}" -> ("nodeName", map[string]string{"alias": "env", "id": "node_id"})
func parseNodeMetadata(part string) (string, map[string]string) {
	metadata := make(map[string]string)

	if matches := metadataRegex.FindStringSubmatch(part); matches != nil {
		nodeName := matches[1]
		metadataStr := matches[2]

		// Parse key-value pairs
		pairs := strings.Split(metadataStr, ",")
		for _, pair := range pairs {
			kv := strings.SplitN(strings.TrimSpace(pair), ":", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				metadata[key] = value
			}
		}

		return nodeName, metadata
	}

	return part, metadata
}

// parsePatternWithParens splits a pattern string into parts, treating parenthesized
// segments as single units. For example:
// - "a.b.c" -> ["a", "b", "c"]
// - "(a.b.c)" -> ["(a.b.c)"]
// - "a.(b.c).d" -> ["a", "(b.c)", "d"]
// - "a.b.(c.d)" -> ["a", "b", "(c.d)"]
func parsePatternWithParens(pattern string) []string {
	var parts []string
	var current strings.Builder
	parenDepth := 0

	for i := 0; i < len(pattern); i++ {
		ch := pattern[i]
		switch ch {
		case '(':
			parenDepth++
			current.WriteByte(ch)
		case ')':
			parenDepth--
			current.WriteByte(ch)
		case '.':
			if parenDepth == 0 {
				// Outside parentheses, this is a separator
				if current.Len() > 0 {
					parts = append(parts, current.String())
					current.Reset()
				}
			} else {
				// Inside parentheses, include the dot
				current.WriteByte(ch)
			}
		default:
			current.WriteByte(ch)
		}
	}

	// Add the last part
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// parenTargetInfo holds parsed information from a parenthesized target pattern
type parenTargetInfo struct {
	parentPath string // The v1 path inside parentheses (e.g., "spec.output")
	hasArray   bool   // Whether the pattern has array notation [i]
	childField string // The child field after the parenthesized part (e.g., "name")
}

// parseParenTarget parses targets like "(spec.output)[i].name" or "(spec.output).name" into:
// - parentPath: "spec.output" (the v1 path inside parentheses)
// - hasArray: true if [i] is present
// - childField: "name" (the v1 name for the child field)
// Returns nil if the target doesn't match the parenthesized format.
func parseParenTarget(target string) *parenTargetInfo {
	match := parenPathRegex.FindStringSubmatch(target)
	if match == nil {
		return nil
	}
	// match[1] = "spec.output", match[2] = "[i]" or "", match[3] = "name"
	return &parenTargetInfo{
		parentPath: match[1],
		hasArray:   match[2] != "",
		childField: match[3],
	}
}

func (t *Trie) insertFromNode(startNode *TrieNode, pattern, target string) {
	// Check if target uses parenthesized syntax: (parent.path)[i].childField or (parent.path).childField
	parenInfo := parseParenTarget(target)
	if parenInfo != nil {
		// Handle parenthesized rule specially
		t.insertParenRule(startNode, pattern, parenInfo)
		return
	}

	// Parse pattern and target with parenthesis support
	patternParts := parsePatternWithParens(pattern)
	targetParts := parsePatternWithParens(target)
	node := startNode

	// Track target index separately since pattern parts may expand
	targetIdx := 0

	for i, part := range patternParts {
		// Check if this pattern part is parenthesized (collapsed multi-segment)
		if match := standaloneParenRegex.FindStringSubmatch(part); match != nil {
			// Parenthesized pattern part like "(a.b.c)" - expand into multiple nodes
			innerPath := match[1]
			innerParts := strings.Split(innerPath, ".")

			// Get the corresponding target part (may also be parenthesized)
			var targetV1Name string
			if targetIdx < len(targetParts) {
				targetPart := targetParts[targetIdx]
				// Strip parentheses from target if present
				if innerMatch := standaloneParenRegex.FindStringSubmatch(targetPart); innerMatch != nil {
					targetV1Name = innerMatch[1]
				} else {
					targetV1Name = targetPart
				}
				targetIdx++
			} else {
				targetV1Name = "-"
			}

			// Create nodes for each inner part
			// All but the last get v1Name="-" (skip in output)
			// The last one gets the target v1Name
			for j, innerPart := range innerParts {
				var v1Name string
				if j == len(innerParts)-1 {
					// Last inner part gets the target name
					v1Name = targetV1Name
				} else {
					// Intermediate parts are skipped in output
					v1Name = "-"
				}

				if node.children[innerPart] == nil {
					node.children[innerPart] = &TrieNode{
						children: make(map[string]*TrieNode),
						v1Name:   v1Name,
					}
				}
				node = node.children[innerPart]
			}
			continue
		}

		// Parse metadata from the part
		nodeName, metadata := parseNodeMetadata(part)

		// Get corresponding target part for v1Name
		var v1Name string
		if targetIdx < len(targetParts) {
			targetPart := targetParts[targetIdx]

			// If this is the last pattern part and there are more target parts,
			// join all remaining target parts as the v1Name
			if i == len(patternParts)-1 && len(targetParts) > targetIdx+1 {
				remainingTarget := targetParts[targetIdx:]
				// Strip parentheses and array notation from first part
				first := remainingTarget[0]
				if innerMatch := standaloneParenRegex.FindStringSubmatch(first); innerMatch != nil {
					remainingTarget[0] = innerMatch[1]
				} else if arrayMatch := arrayNotationRegex.FindStringSubmatch(first); arrayMatch != nil {
					remainingTarget[0] = arrayMatch[1]
				}
				v1Name = strings.Join(remainingTarget, ".")
			} else {
				// Strip parentheses from target if present
				if innerMatch := standaloneParenRegex.FindStringSubmatch(targetPart); innerMatch != nil {
					v1Name = innerMatch[1]
				} else if arrayMatch := arrayNotationRegex.FindStringSubmatch(targetPart); arrayMatch != nil {
					v1Name = arrayMatch[1]
				} else {
					v1Name = targetPart
				}
			}
			targetIdx++
		} else {
			v1Name = "-"
		}

		// Handle array notation in pattern
		arrayMatch := arrayNotationRegex.FindStringSubmatch(nodeName)
		if arrayMatch != nil {
			baseName := arrayMatch[1]
			if node.children[baseName] == nil {
				node.children[baseName] = &TrieNode{
					children: make(map[string]*TrieNode),
					isArray:  true,
					v1Name:   v1Name,
				}
				// Apply metadata
				t.applyMetadata(node.children[baseName], metadata)
			}
			node = node.children[baseName]
		} else if nodeName == "*" {
			if node.wildcardChild == nil {
				node.wildcardChild = &TrieNode{
					children:   make(map[string]*TrieNode),
					isWildcard: true,
					v1Name:     v1Name,
				}
				// Apply metadata
				t.applyMetadata(node.wildcardChild, metadata)
			}
			node = node.wildcardChild
		} else {
			if node.children[nodeName] == nil {
				node.children[nodeName] = &TrieNode{
					children: make(map[string]*TrieNode),
					v1Name:   v1Name,
				}
				// Apply metadata
				t.applyMetadata(node.children[nodeName], metadata)
			}
			node = node.children[nodeName]
		}
	}

	node.isEnd = true
}

// insertParenRule handles rules with parenthesized syntax like:
// pattern: "outputVariables[i].value" -> target: "(spec.output)[i].name"
// This creates:
// - An array node "outputVariables" with arrayParentV1Path="spec.output" (if hasArray)
// - A child node "value" with v1Name="name"
//
// Also handles non-array parenthesized rules like:
// pattern: "a.b" -> target: "(x.y).z"
// This creates:
// - Node "a" with v1Name="x.y"
// - Node "b" with v1Name="z"
func (t *Trie) insertParenRule(startNode *TrieNode, pattern string, info *parenTargetInfo) {
	patternParts := strings.Split(pattern, ".")
	node := startNode

	for i, part := range patternParts {
		nodeName, metadata := parseNodeMetadata(part)

		// Check if this part has array notation
		arrayMatch := arrayNotationRegex.FindStringSubmatch(nodeName)
		if arrayMatch != nil {
			// This is the array node (e.g., "outputVariables[i]")
			baseName := arrayMatch[1]
			if node.children[baseName] == nil {
				node.children[baseName] = &TrieNode{
					children:          make(map[string]*TrieNode),
					isArray:           info.hasArray,
					arrayParentV1Path: info.parentPath,
				}
				t.applyMetadata(node.children[baseName], metadata)
			} else {
				// Update existing node with array parent path
				// This allows both indexed and non-indexed rules to coexist
				if info.hasArray {
					node.children[baseName].arrayParentV1Path = info.parentPath
					node.children[baseName].isArray = true
				}
			}
			node = node.children[baseName]
		} else if nodeName == "*" {
			// Wildcard child
			if node.wildcardChild == nil {
				node.wildcardChild = &TrieNode{
					children:   make(map[string]*TrieNode),
					isWildcard: true,
					v1Name:     info.childField, // Use childField for wildcard after array
				}
				t.applyMetadata(node.wildcardChild, metadata)
			}
			node = node.wildcardChild
		} else {
			// Regular child node
			// For the first part (before array), use parentPath as v1Name if no array in pattern yet
			// For the last part, use childField as v1Name
			var v1Name string
			if i == len(patternParts)-1 {
				v1Name = info.childField
			} else if i == 0 && !info.hasArray {
				// Non-array parenthesized rule: first node maps to parentPath
				v1Name = info.parentPath
			} else {
				v1Name = nodeName
			}
			if node.children[nodeName] == nil {
				node.children[nodeName] = &TrieNode{
					children: make(map[string]*TrieNode),
					v1Name:   v1Name,
				}
				t.applyMetadata(node.children[nodeName], metadata)
			} else if i == len(patternParts)-1 {
				// Update v1Name for existing terminal node
				node.children[nodeName].v1Name = v1Name
			}
			node = node.children[nodeName]
		}
	}

	node.isEnd = true
}

// applyMetadata applies metadata map to a node
func (t *Trie) applyMetadata(node *TrieNode, metadata map[string]string) {
	if alias, ok := metadata["alias"]; ok {
		node.alias = alias
		t.aliasIndex[alias] = append(t.aliasIndex[alias], node)
	}
	if id, ok := metadata["id"]; ok {
		node.id = id
		t.nodeIndex[id] = node
	}
}

// pathPart represents a parsed path segment
type pathPart struct {
	name       string
	arrayIndex string // e.g., "[0]", "[1]"
}

// matchContext holds state during path matching
type matchContext struct {
	arrayIndices []string // Stack of array indices encountered
	v1Path       []string // Built v1 path segments

	// Step context tracking for lazy resolution
	lastStepID  string // Last step ID seen in path (captured from wildcard after "steps")
	inStepsPath bool   // True if we're currently inside a "steps" path segment

	// Scoped context for disambiguating duplicate step IDs
	lastStageID      string // Stage ID captured from wildcard after "stages" node
	lastGroupID      string // Step group ID captured when a step-group wildcard is resolved
	inStepGroupAlias bool   // True if we entered through the "stepGroup" alias root

	// parentHops counts the number of getParentStepGroup edges traversed since
	// entering the "stepGroup" alias.
	parentHops int

	// FQN mode tracking
	fqnAttempted bool // True if we've already attempted FQN conversion (prevents infinite recursion)

	// suppressFQN disables FQN expansion for the current match; set when entry
	// is at a noFQNOverride node so the leaf step stays structural.
	suppressFQN bool
}

func (t *Trie) Match(path string, context *ConversionContext) (string, bool) {
	parts := t.parsePath(path)
	if len(parts) == 0 {
		return path, false
	}

	// Try alias-based matching on the first segment
	firstSegment := parts[0].name
	if aliasedNodes, exists := t.aliasIndex[firstSegment]; exists {
		var bestResult string
		var bestMatched bool
		bestScore := -1

		for _, aliasNode := range aliasedNodes {
			ctx := &matchContext{
				arrayIndices: []string{},
				v1Path:       []string{},
			}

			// Bracketed alias root (stages[X] / steps[X] without the "pipeline."
			// prefix): matchBracketedKeyword emits a combined "keyword[X]" segment
			// and captures the bracket content as the stage/step ID.
			if parts[0].arrayIndex != "" && (firstSegment == "stages" || firstSegment == "steps") && aliasNode.wildcardChild != nil {
				if firstSegment == "steps" {
					ctx.inStepsPath = true
					ctx.inStepGroupAlias = true
				}
				if result, matched := t.matchBracketedKeyword(nil, aliasNode, parts, 0, ctx, context); matched {
					score := t.calculateMatchScore(aliasNode, parts, 1, context)
					if score > bestScore {
						bestScore = score
						bestResult = result
						bestMatched = true
					}
				}
				continue
			}

			// Stage-relative execution roots (execution.* / spec.execution.*) carry
			// no explicit stage, so they later emit the v1 "stage" self-reference
			// keyword. FQN suppression is handled below via noFQNOverride.
			isExecutionRoot := firstSegment == "execution"
			isSpecExecutionRoot := firstSegment == "spec" && len(parts) > 1 && parts[1].name == "execution"

			// noFQNOverride (set at build time on the stage/ancestor and execution
			// root nodes) disables FQN expansion: these already carry stage context.
			if aliasNode.noFQNOverride {
				ctx.suppressFQN = true
			}

			// Seed v1Path by alias root: execution roots emit the "stage"
			// self-reference; a "step." spec/output reference in FQN mode seeds the
			// current step's FQN; otherwise use the alias node's v1Name.
			if isExecutionRoot || isSpecExecutionRoot {
				ctx.v1Path = append(ctx.v1Path, "stage")
			} else if firstSegment == "step" && context != nil && context.UseFQN && context.CurrentFQN != "" &&
				len(parts) > 1 && (parts[1].name == "spec" || parts[1].name == "output") {
				ctx.v1Path = strings.Split(context.CurrentFQN, ".")
			} else if aliasNode.v1Name != "" && aliasNode.v1Name != "-" {
				if aliasNode.v1Name == "*" {
					ctx.v1Path = append(ctx.v1Path, parts[0].name)
				} else {
					ctx.v1Path = append(ctx.v1Path, aliasNode.v1Name)
				}
			}

			// Set inStepsPath if we're matching "steps" alias - this enables step ID capture
			if firstSegment == "steps" {
				ctx.inStepsPath = true
			}

			// "stepGroup"/"steps" are group-relative: resolution scopes them to the
			// call-site's CurrentStepGroupChain.
			if firstSegment == "stepGroup" || firstSegment == "steps" {
				ctx.inStepGroupAlias = true
			}

			if result, matched := t.matchRecursive(aliasNode, parts, 1, ctx, context); matched {
				// A bare stage-spec relative entry ("spec.<stageField>...", e.g.
				// spec.service.identifier or spec.infrastructure.spec.connectorRef)
				// is stage-scoped, so emit the v1 "stage" self-reference prefix,
				// mirroring the equivalent "stage.spec.*" form. The spec.execution.*
				// root already seeds "stage" above, so skip the prefix there. Step
				// spec relatives (spec.command -> spec.script) resolve through
				// step_spec_node and are unaffected.
				if aliasNode.id == "stage_spec_node" && !isSpecExecutionRoot {
					result = "stage." + result
				}
				score := t.calculateMatchScore(aliasNode, parts, 1, context)
				if score > bestScore {
					bestScore = score
					bestResult = result
					bestMatched = true
				}
			}
		}

		if bestMatched {
			return bestResult, true
		}
	}

	// Fallback: try from root
	ctx := &matchContext{
		arrayIndices: []string{},
		v1Path:       []string{},
	}
	if result, matched := t.matchRecursive(t.root, parts, 0, ctx, context); matched {
		return result, true
	}

	return path, false
}

// calculateMatchScore returns the number of non-wildcard matches in a path
func (t *Trie) calculateMatchScore(node *TrieNode, parts []pathPart, index int, context *ConversionContext) int {
	score := 0
	currentNode := node

	for i := index; i < len(parts); i++ {
		part := parts[i]

		// Exact match (non-wildcard) gets a point
		if child, exists := currentNode.children[part.name]; exists {
			score++
			currentNode = child
		} else if currentNode.wildcardChild != nil {
			// Wildcard match, no points
			currentNode = currentNode.wildcardChild
		} else {
			// No match, stop counting
			break
		}
	}

	return score
}

// parsePath splits a path into parts, handling array indices and nested <+...> expressions.
func (t *Trie) parsePath(path string) []pathPart {
	segments := splitPathSegments(path)
	var parts []pathPart

	for _, seg := range segments {
		if arrayMatch := arrayNotationRegex.FindStringSubmatch(seg); arrayMatch != nil {
			parts = append(parts, pathPart{
				name:       arrayMatch[1],
				arrayIndex: arrayMatch[2],
			})
		} else {
			parts = append(parts, pathPart{name: seg})
		}
	}

	return parts
}

// hasResolvableContext returns true if the ConversionContext can resolve a step type.
func hasResolvableContext(convContext *ConversionContext, ctx *matchContext) bool {
	return convContext != nil && (convContext.StepType != "" || convContext.CurrentStepType != "" ||
		len(convContext.CurrentStepTypes) > 0 ||
		(ctx.lastStepID != "" && convContext.StepInfoByFQN != nil))
}

// tryMatchChild attempts to match remaining parts through a child node using both
// context-specific and general rules. The order depends on whether context is available.
// isSkipped indicates the child node suppresses its output (v1Name="-").
// skippedSegment is the original segment name when isSkipped=true (used for passthrough).
func (t *Trie) tryMatchChild(child *TrieNode, parts []pathPart, nextIndex int, ctx *matchContext, convContext *ConversionContext, isSkipped bool, skippedSegment string) (string, bool) {
	hasCtx := hasResolvableContext(convContext, ctx)

	if hasCtx {
		// 1) Candidate-specific rules first (resolved single type or clash-free family).
		if result, matched := t.tryContextMatchCandidates(child, parts, nextIndex, ctx, convContext); matched {
			return result, true
		}
		// 2) General (non-context) rules.
		if result, matched := t.matchRecursive(child, parts, nextIndex, ctx, convContext); matched {
			if !isSkipped || !t.isSkippedNodePassthrough(result, ctx, parts, nextIndex) {
				return result, true
			}
			// Skipped node passthrough: only preserve segment if path ENDS at the skipped node
			// (not when children matched but didn't transform)
			if nextIndex >= len(parts) {
				return result + "." + skippedSegment, true
			}
			// Children matched with passthrough - don't accept this match
		}
		// 3) Deterministic all-rules best-match fallback. Internally gated to a
		// no-op when a single precise type resolved, so the v0-context flow
		// never reaches it; only multi-type v1 families and unmapped uses do.
		if result, matched := t.tryContextMatchAny(child, parts, nextIndex, ctx, convContext); matched {
			return result, true
		}
	} else {
		// No context: try general rules first, then deterministic context fallback
		if result, matched := t.matchRecursive(child, parts, nextIndex, ctx, convContext); matched {
			if !isSkipped || !t.isSkippedNodePassthrough(result, ctx, parts, nextIndex) {
				return result, true
			}
			// Skipped node passthrough: only preserve segment if path ENDS at the skipped node
			if nextIndex >= len(parts) {
				return result + "." + skippedSegment, true
			}
			// Children matched with passthrough - don't accept this match
		}
		if result, matched := t.tryContextMatch(child, parts, nextIndex, ctx, convContext); matched {
			return result, true
		}
	}

	return "", false
}

func (t *Trie) matchRecursive(node *TrieNode, parts []pathPart, index int, ctx *matchContext, convContext *ConversionContext) (string, bool) {
	// Base case: consumed all parts
	if index == len(parts) {
		return strings.Join(ctx.v1Path, "."), true
	}

	currentPart := parts[index]

	// Bracketed keyword forms stages[X] / steps[X]: the bracket content acts as
	// the next wildcard segment (stage/step ID). See matchBracketedKeyword.
	if currentPart.arrayIndex != "" && (currentPart.name == "stages" || currentPart.name == "steps") {
		if child, ok := node.children[currentPart.name]; ok && child.wildcardChild != nil {
			if result, matched := t.matchBracketedKeyword(node, child, parts, index, ctx, convContext); matched {
				return result, true
			}
		}
	}

	if currentPart.arrayIndex != "" {
		ctx.arrayIndices = append(ctx.arrayIndices, currentPart.arrayIndex)
	}

	// Try exact match first
	if child, exists := node.children[currentPart.name]; exists {
		v1Segment := t.buildV1Segment(child, currentPart)
		if v1Segment != "" {
			ctx.v1Path = append(ctx.v1Path, v1Segment)
		}

		// Track "steps" path for step ID capture
		wasInStepsPath := ctx.inStepsPath
		if currentPart.name == "steps" {
			ctx.inStepsPath = true
		}

		// Each parent_step_group_node hop moves one level up from the call-site's
		// group; resolution trims that many trailing elements off the chain.
		savedParentHops := ctx.parentHops
		if child.id == "parent_step_group_node" {
			ctx.parentHops++
		}

		isSkipped := child.v1Name == "-"
		skippedSegment := ""
		if isSkipped {
			skippedSegment = currentPart.name
			if currentPart.arrayIndex != "" {
				skippedSegment += currentPart.arrayIndex
			}
		}

		if result, matched := t.tryMatchChild(child, parts, index+1, ctx, convContext, isSkipped, skippedSegment); matched {
			return result, true
		}

		// Backtrack
		ctx.parentHops = savedParentHops
		ctx.inStepsPath = wasInStepsPath
		if v1Segment != "" {
			ctx.v1Path = ctx.v1Path[:len(ctx.v1Path)-1]
		}
	}

	// Try wildcard match
	if node.wildcardChild != nil {
		v1Segment := t.buildV1Segment(node.wildcardChild, currentPart)
		if v1Segment == "" {
			// Wildcard with no v1Name or v1Name="-", use part name
			v1Segment = currentPart.name
			if currentPart.arrayIndex != "" {
				v1Segment += currentPart.arrayIndex
			}
		}

		ctx.v1Path = append(ctx.v1Path, v1Segment)

		// Capture step ID when matching wildcard after "steps"
		// This enables lazy step type resolution for step.spec expressions
		savedStepID := ctx.lastStepID
		savedStageID := ctx.lastStageID
		savedGroupID := ctx.lastGroupID
		savedV1Path := make([]string, len(ctx.v1Path))
		copy(savedV1Path, ctx.v1Path)

		// Capture stage ID when this wildcard is directly under "stages_node"
		if node.id == "stages_node" {
			ctx.lastStageID = currentPart.name
		}

		if ctx.inStepsPath {
			ctx.lastStepID = currentPart.name
			ctx.inStepsPath = false // Reset after capturing

			// FQN substitution applies only when the step reference targets a spec
			// or output node; other fields (status, name, ...) stay structural.
			targetsSpecOrOutput := stepRefTargetsSpecOrOutput(parts, index)

			// In FQN mode, merge the accumulated path with the step's v1 FQN base.
			if convContext != nil && convContext.UseFQN && targetsSpecOrOutput && !ctx.suppressFQN {
				stepID := currentPart.name
				v1BasePath := resolveStepV1Path(stepID, ctx, convContext)

				if v1BasePath != "" {
					// mergeFQN preserves dynamic segments in the accumulated path.
					rawSeg := currentPart.name
					if currentPart.arrayIndex != "" {
						rawSeg += currentPart.arrayIndex
					}
					ctx.v1Path = mergeFQN(ctx.v1Path, v1BasePath, rawSeg)
				}
			}

			// If the captured ID is a step group, append it to the chain so nested
			// steps resolve within it.
			if convContext != nil && targetsSpecOrOutput {
				promoteGroupChain(ctx, currentPart.name, convContext)
			}
		}

		if result, matched := t.tryMatchChild(node.wildcardChild, parts, index+1, ctx, convContext, false, ""); matched {
			return result, true
		}

		// Backtrack
		ctx.lastStepID = savedStepID
		ctx.lastStageID = savedStageID
		ctx.lastGroupID = savedGroupID
		ctx.v1Path = savedV1Path[:len(savedV1Path)-1]
	}

	// Try context match directly on current node if it has contextChildren
	// This handles cases where we're at a node like step_spec_node that has no direct children
	// but has context-specific rules attached (e.g., ShellScript environmentVariables rules)
	if node.contextChildren != nil && len(node.contextChildren) > 0 {
		if result, matched := t.tryContextMatch(node, parts, index, ctx, convContext); matched {
			return result, true
		}
	}

	// No match — return partial conversion if we have progress
	if len(ctx.v1Path) > 0 {
		remaining := t.buildRemainingPath(parts, index)
		return strings.Join(ctx.v1Path, ".") + "." + remaining, true
	}

	return "", false
}

// tryContextMatch attempts to match remaining parts through a node's context-specific sub-trie.
// Returns the converted path and true if a context match was found.
// When no context is provided, it tries all available context types and returns the first match
// deterministically (by sorting context keys alphabetically).
//
// LAZY STEP TYPE RESOLUTION:
// If convContext.StepType is empty but we have a step ID captured in matchContext,
// we resolve the step type from StepTypeMap before attempting context matching.
func (t *Trie) tryContextMatch(node *TrieNode, parts []pathPart, index int, ctx *matchContext, convContext *ConversionContext) (string, bool) {
	if result, matched := t.tryContextMatchCandidates(node, parts, index, ctx, convContext); matched {
		return result, true
	}
	return t.tryContextMatchAny(node, parts, index, ctx, convContext)
}

// tryContextMatchCandidates matches against only the resolved candidate step
// types (a single confident type, or a clash-free family like run / run-test).
// Returns false when no candidate resolves or none has a matching sub-trie.
func (t *Trie) tryContextMatchCandidates(node *TrieNode, parts []pathPart, index int, ctx *matchContext, convContext *ConversionContext) (string, bool) {
	if node.contextChildren == nil {
		return "", false
	}
	cands := t.resolveStepTypeCandidates(ctx, convContext)
	if len(cands) == 0 {
		return "", false
	}
	var keys []string
	for _, c := range cands {
		if _, ok := node.contextChildren[c]; ok {
			keys = append(keys, c)
		}
	}
	if len(keys) == 0 {
		return "", false
	}
	return t.bestContextMatchAmong(node, keys, parts, index, ctx, convContext)
}

// tryContextMatchAny is the deterministic fallback: try every available context
// sub-trie and keep the one with the most node matches against the path (ties
// broken alphabetically). Used when no step type resolves or candidate matching
// fails; emits a warning so callers can surface the best-effort guess.
func (t *Trie) tryContextMatchAny(node *TrieNode, parts []pathPart, index int, ctx *matchContext, convContext *ConversionContext) (string, bool) {
	if node.contextChildren == nil {
		return "", false
	}
	// Gate: when a single, precise step type resolved, do NOT fall back to the
	// all-rules best-match. This keeps the v0-context flow (which always
	// resolves exactly one precise type) unchanged
	if len(t.resolveStepTypeCandidates(ctx, convContext)) == 1 {
		return "", false
	}
	var contextKeys []string
	for key := range node.contextChildren {
		contextKeys = append(contextKeys, key)
	}
	if len(contextKeys) == 0 {
		return "", false
	}
	result, matched := t.bestContextMatchAmong(node, contextKeys, parts, index, ctx, convContext)
	// Warn only when some step context was expected (FQN map or an explicit
	// step type) but the type couldn't be pinned down — i.e. an unmapped
	// template/approval `uses` or an unresolved candidate. Pure no-context
	// conversions deterministically use best-match and need no warning.
	if matched && contextExpected(convContext) {
		convContext.addWarning("ambiguous step type for path '" + t.buildRemainingPath(parts, index) +
			"'; resolved via best-match across all step rule sets")
	}
	return result, matched
}

// contextExpected reports whether the caller supplied step context (an FQN map
// or an explicit/current step type) such that a best-match fallback is worth
// flagging as a warning.
func contextExpected(cc *ConversionContext) bool {
	return cc != nil && (cc.StepInfoByFQN != nil || cc.StepType != "" ||
		cc.CurrentStepType != "" || len(cc.CurrentStepTypes) > 0)
}

// bestContextMatchAmong tries each given context key's sub-trie and returns the
// match with the highest node-match score; ties are broken alphabetically.
func (t *Trie) bestContextMatchAmong(node *TrieNode, keys []string, parts []pathPart, index int, ctx *matchContext, convContext *ConversionContext) (string, bool) {
	sortStrings(keys)

	var bestResult string
	bestScore := -1
	bestKey := ""

	for _, contextKey := range keys {
		contextRoot := node.contextChildren[contextKey]
		if contextRoot == nil {
			continue
		}
		result, matched := t.tryContextSubtree(contextRoot, parts, index, ctx, convContext)
		if matched {
			score := t.countContextNodeMatches(contextRoot, parts, index)
			if score > bestScore || (score == bestScore && (bestKey == "" || contextKey < bestKey)) {
				bestScore = score
				bestResult = result
				bestKey = contextKey
			}
		}
	}

	if bestScore >= 0 {
		return bestResult, true
	}
	return "", false
}

// isSkippedNodePassthrough checks whether a result from a skipped (v1Name="-") node's
// deeper recursion is just a passthrough — i.e., the skipped node didn't contribute
// any real transformation. This covers two cases:
//   - Base case: path ends at the skipped node, result == v1Path joined (nothing added)
//   - Partial fallback: remaining parts are passed through raw, result == v1Path + "." + raw remaining
func (t *Trie) isSkippedNodePassthrough(result string, ctx *matchContext, parts []pathPart, fromIndex int) bool {
	v1Joined := strings.Join(ctx.v1Path, ".")

	// Base case: path ended at the skipped node
	if fromIndex >= len(parts) {
		return result == v1Joined
	}

	// Partial fallback: remaining parts passed through unchanged
	remaining := t.buildRemainingPath(parts, fromIndex)
	if remaining == "" {
		return result == v1Joined
	}
	return result == v1Joined+"."+remaining
}

// countContextNodeMatches counts the number of exact (non-wildcard) node matches
// in a context subtree against the remaining input path parts.
// This is used to determine the best-match context when no context is provided.
func (t *Trie) countContextNodeMatches(contextRoot *TrieNode, parts []pathPart, startIndex int) int {
	score := 0
	currentNode := contextRoot

	for i := startIndex; i < len(parts); i++ {
		part := parts[i]

		if child, exists := currentNode.children[part.name]; exists {
			score++
			currentNode = child
		} else if currentNode.wildcardChild != nil {
			currentNode = currentNode.wildcardChild
		} else {
			break
		}
	}

	return score
}

// tryContextSubtree attempts to match through a specific context sub-trie
func (t *Trie) tryContextSubtree(contextRoot *TrieNode, parts []pathPart, index int, ctx *matchContext, convContext *ConversionContext) (string, bool) {
	// Try matching remaining parts through the context sub-trie
	contextCtx := &matchContext{
		arrayIndices: ctx.arrayIndices,
		v1Path:       make([]string, len(ctx.v1Path)),
	}
	copy(contextCtx.v1Path, ctx.v1Path)

	result, matched := t.matchRecursive(contextRoot, parts, index, contextCtx, convContext)
	if matched {
		// Verify this is a complete match, not a partial fallback
		// The partial match fallback in matchRecursive can return true even when
		// it didn't properly match through the context sub-trie, so we need to validate.
		// A partial fallback produces: v1Path + "." + remainingRaw (raw parts unchanged)
		// A real context match produces: v1Path + transformed parts from context sub-trie
		remainingRaw := t.buildRemainingPath(parts, index)
		if remainingRaw != "" {
			expectedPartialFallback := strings.Join(contextCtx.v1Path, ".") + "." + remainingRaw
			if result == expectedPartialFallback {
				// This is exactly a partial fallback, not a real context match
				return "", false
			}
		}

		return result, true
	}

	return "", false
}

// sortStrings sorts a slice of strings in place.
func sortStrings(strs []string) {
	sort.Strings(strs)
}

// matchBracketedKeyword handles stages[X] / steps[X], where X is an index or an
// (optionally quoted/dynamic) identifier. It advances through the keyword's
// wildcard child in one step, emits a combined "keyword[X]" segment, and uses
// the bracket content for ID context (lastStageID/lastStepID, group promotion,
// FQN merging) like a normal wildcard. Numeric indices are emitted verbatim
// without FQN lookup.
func (t *Trie) matchBracketedKeyword(node, keywordNode *TrieNode, parts []pathPart, index int, ctx *matchContext, convContext *ConversionContext) (string, bool) {
	currentPart := parts[index]
	wildcardChild := keywordNode.wildcardChild

	// Save state for backtracking
	savedV1Path := make([]string, len(ctx.v1Path))
	copy(savedV1Path, ctx.v1Path)
	savedStepID := ctx.lastStepID
	savedStageID := ctx.lastStageID
	savedGroupID := ctx.lastGroupID
	savedInStepsPath := ctx.inStepsPath

	// Emit the combined keyword[X] segment. Keyword v1Name comes from the
	// keywordNode (e.g. "stages" / "steps"); fall back to part name when blank.
	keywordV1 := keywordNode.v1Name
	if keywordV1 == "" || keywordV1 == "-" {
		keywordV1 = currentPart.name
	}
	combinedSeg := keywordV1 + currentPart.arrayIndex
	ctx.v1Path = append(ctx.v1Path, combinedSeg)

	// Extract and classify the bracket content.
	bracketID := extractBracketContent(currentPart.arrayIndex)
	isIndex := isNumericIndex(bracketID)

	// Mirror the wildcard-capture context updates from matchRecursive.
	if keywordNode.id == "stages_node" && !isIndex {
		ctx.lastStageID = bracketID
	}

	if currentPart.name == "steps" && !isIndex {
		ctx.lastStepID = bracketID
		ctx.inStepsPath = false

		// FQN substitution applies only for spec/output references.
		targetsSpecOrOutput := stepRefTargetsSpecOrOutput(parts, index)

		// In FQN mode, merge the accumulated path with the step's v1 FQN base.
		if convContext != nil && convContext.UseFQN && targetsSpecOrOutput && !ctx.suppressFQN {
			v1Base := resolveStepV1Path(bracketID, ctx, convContext)
			if v1Base != "" {
				// Drop the just-appended combinedSeg before merging, then re-add.
				accPath := ctx.v1Path[:len(ctx.v1Path)-1]
				ctx.v1Path = mergeFQN(accPath, v1Base, combinedSeg)
			}
		}

		// Promote lastGroupID when the captured ID resolves to a step group.
		if convContext != nil && targetsSpecOrOutput {
			promoteGroupChain(ctx, bracketID, convContext)
		}
	}

	if result, matched := t.tryMatchChild(wildcardChild, parts, index+1, ctx, convContext, false, ""); matched {
		return result, true
	}

	// Backtrack on failure so the caller can fall back to normal matching.
	ctx.v1Path = savedV1Path
	ctx.lastStepID = savedStepID
	ctx.lastStageID = savedStageID
	ctx.lastGroupID = savedGroupID
	ctx.inStepsPath = savedInStepsPath
	return "", false
}

// extractBracketContent returns the content between the outermost square
// brackets, stripping one surrounding pair of quotes if present. Examples:
//
//	"[0]"                          -> "0"
//	"['approval_dev']"             -> "approval_dev"
//	"[\"approval_dev\"]"           -> "approval_dev"
//	"[<+stage.var.X>]"             -> "<+stage.var.X>"
//	"['approval_<+stage.env.name>']" -> "approval_<+stage.env.name>"
func extractBracketContent(bracketStr string) string {
	s := strings.TrimPrefix(bracketStr, "[")
	s = strings.TrimSuffix(s, "]")
	if len(s) >= 2 {
		first, last := s[0], s[len(s)-1]
		if (first == '\'' && last == '\'') || (first == '"' && last == '"') {
			s = s[1 : len(s)-1]
		}
	}
	return s
}

// appendGroupChain appends id to a '>'-separated step-group ancestry chain
// (e.g. "Outer>Inner"), matching the registration key format in postprocess.go.
func appendGroupChain(chain, id string) string {
	if chain == "" {
		return id
	}
	return chain + ">" + id
}

// isNumericIndex reports whether s is a non-empty string of ASCII digits.
func isNumericIndex(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// buildV1Segment builds the v1 output segment for a node
func (t *Trie) buildV1Segment(node *TrieNode, part pathPart) string {
	// Handle array nodes with arrayParentV1Path (from parenthesized rules)
	// e.g., outputVariables[1] with arrayParentV1Path="spec.output" -> "spec.output[1]"
	// Only use arrayParentV1Path when there's an actual array index in the path
	if node.arrayParentV1Path != "" && node.isArray && part.arrayIndex != "" {
		return node.arrayParentV1Path + part.arrayIndex
	}

	if node.v1Name == "-" {
		// Skip this segment in output
		return ""
	}

	v1Name := node.v1Name
	if v1Name == "" {
		v1Name = part.name
	}

	// For wildcard nodes, replace "*" in v1Name with the actual matched value
	if node.isWildcard && strings.Contains(v1Name, "*") {
		v1Name = strings.Replace(v1Name, "*", part.name, 1)
	}

	// Append array index if present
	if part.arrayIndex != "" && node.isArray {
		return v1Name + part.arrayIndex
	}

	return v1Name
}

// buildRemainingPath constructs path from remaining parts
func (t *Trie) buildRemainingPath(parts []pathPart, startIndex int) string {
	var segments []string
	for i := startIndex; i < len(parts); i++ {
		segment := parts[i].name
		if parts[i].arrayIndex != "" {
			segment += parts[i].arrayIndex
		}
		segments = append(segments, segment)
	}
	return strings.Join(segments, ".")
}

// resolveStepTypeCandidates mirrors resolveStepType but returns candidate
// types: a single confident type, or a clash-free family (e.g. run /
// run-test) when one v1 field maps to several v0 types. Empty when nothing
// resolves (the caller then uses the all-rules best-match fallback).
func (t *Trie) resolveStepTypeCandidates(ctx *matchContext, convContext *ConversionContext) []string {
	if convContext == nil {
		return nil
	}

	// Priority 1: explicitly pre-set step type.
	if convContext.StepType != "" {
		return []string{convContext.StepType}
	}

	// Priority 2: captured step ID via the FQN resolver.
	if ctx.lastStepID != "" && convContext.StepInfoByFQN != nil {
		stageRef := ctx.lastStageID
		if stageRef == "" {
			stageRef = convContext.CurrentStageID
		}
		chainRef := splitChain(ctx.lastGroupID)
		aliasKind := AliasAbsolute
		if ctx.inStepGroupAlias {
			aliasKind = AliasGroupRelative
		}
		if _, info, outcome := convContext.resolveStepFQN(stageRef, ctx.lastStepID, chainRef, aliasKind, false, ctx.parentHops); info != nil &&
			(outcome == OutcomeExactChain || outcome == OutcomeSubsequenceWithGap || outcome == OutcomeFlatUniqueFallback) {
			if len(info.Types) > 0 {
				return info.Types
			}
			if info.Type != "" {
				return []string{info.Type}
			}
		}
	}

	// Priority 3: the step we're physically inside.
	if len(convContext.CurrentStepTypes) > 0 {
		return convContext.CurrentStepTypes
	}
	if convContext.CurrentStepType != "" {
		return []string{convContext.CurrentStepType}
	}
	return nil
}

// mergeFQN overlays the trie-accumulated path (accPath) onto the resolved FQN
// base (fqnBase), preserving accPath's dynamic <+...> segments in their ID
// positions. ID positions are matched per keyword ("stages"/"steps") right-to-
// left, so a dynamic step-group ID is never placed in a stage slot; the last
// segment is then replaced with rawSeg to keep its raw captured form.
//
// Example: accPath ["steps","<+stepGroup.identifier>","steps","Sonar_Scan"] over
// fqnBase "pipeline.stages.CI.steps.Python_Sonar.steps.Sonar_Scan" yields
// "pipeline.stages.CI.steps.<+stepGroup.identifier>.steps.Sonar_Scan". When
// accPath is empty, returns fqnBase with its last segment replaced by rawSeg.
func mergeFQN(accPath []string, fqnBase string, rawSeg string) []string {
	if fqnBase == "" {
		return accPath
	}
	fqnParts := strings.Split(fqnBase, ".")

	if len(accPath) == 0 {
		// No accumulated path; use FQN directly, replacing last with rawSeg
		if rawSeg != "" && len(fqnParts) > 0 {
			fqnParts[len(fqnParts)-1] = rawSeg
		}
		return fqnParts
	}

	// Start with fqnParts as the base result
	merged := make([]string, len(fqnParts))
	copy(merged, fqnParts)

	// Partition ID positions by preceding keyword. Stage IDs and step IDs must
	// be matched independently so a step-group dynamic segment never lands in
	// a stage slot (and vice versa).
	fqnStageIDs := extractIDPositionsByKind(fqnParts, "stages")
	fqnStepIDs := extractStepIDPositions(fqnParts)
	accStageIDs := extractIDPositionsByKind(accPath, "stages")
	accStepIDs := extractStepIDPositions(accPath)

	substituteDynamicByKind(merged, fqnStageIDs, accPath, accStageIDs)
	substituteDynamicByKind(merged, fqnStepIDs, accPath, accStepIDs)

	// Always apply rawSeg to the last segment (preserves the raw captured form)
	if rawSeg != "" && len(merged) > 0 {
		merged[len(merged)-1] = rawSeg
	}

	// Restore bracketed forms so a dynamic bracketed reference (e.g.
	// stages[<+stage.identifier>]) isn't rewritten to the FQN's static ID.
	merged = collapseBracketedKeywords(merged, accPath)

	return merged
}

// collapseBracketedKeywords rewrites "kind"+"ID" pairs in merged back into
// combined "kind[X]" segments where accPath used the bracketed form, matched
// left-to-right per kind (stages or steps).
func collapseBracketedKeywords(merged, accPath []string) []string {
	var stageBrackets, stepBrackets []string
	for _, seg := range accPath {
		if !strings.HasSuffix(seg, "]") {
			continue
		}
		if strings.HasPrefix(seg, "stages[") {
			stageBrackets = append(stageBrackets, seg)
		} else if strings.HasPrefix(seg, "steps[") {
			stepBrackets = append(stepBrackets, seg)
		}
	}
	if len(stageBrackets) == 0 && len(stepBrackets) == 0 {
		return merged
	}
	result := make([]string, 0, len(merged))
	si, ti := 0, 0
	for i := 0; i < len(merged); i++ {
		if i+1 < len(merged) && merged[i] == "stages" && si < len(stageBrackets) {
			result = append(result, stageBrackets[si])
			si++
			i++ // skip following ID slot
			continue
		}
		if i+1 < len(merged) && merged[i] == "steps" && ti < len(stepBrackets) {
			result = append(result, stepBrackets[ti])
			ti++
			i++ // skip following ID slot
			continue
		}
		result = append(result, merged[i])
	}
	return result
}

// substituteDynamicByKind matches accIDs to fqnIDs right-to-left (so the most
// recent / closest entries align first) and writes any dynamic accPath segment
// into the corresponding merged slot.
func substituteDynamicByKind(merged []string, fqnIDs []int, accPath []string, accIDs []int) {
	for i := 0; i < len(accIDs) && i < len(fqnIDs); i++ {
		accIdx := accIDs[len(accIDs)-1-i]
		fqnIdx := fqnIDs[len(fqnIDs)-1-i]
		accSeg := accPath[accIdx]
		if strings.Contains(accSeg, "<+") {
			merged[fqnIdx] = accSeg
		}
	}
}

// extractIDPositions returns the indices of segments that are identifiers
// (stage IDs, step IDs, group IDs) — i.e., segments immediately after
// "stages" or "steps".
func extractIDPositions(parts []string) []int {
	var positions []int
	for i := range parts {
		if i > 0 && (parts[i-1] == "stages" || parts[i-1] == "steps") {
			positions = append(positions, i)
		}
	}
	return positions
}

// extractIDPositionsByKind returns the indices of segments whose immediately
// preceding segment equals kind (e.g. "stages" or "steps").
func extractIDPositionsByKind(parts []string, kind string) []int {
	var positions []int
	for i := range parts {
		if i > 0 && parts[i-1] == kind {
			positions = append(positions, i)
		}
	}
	return positions
}

// extractStepIDPositions returns the indices of step ID segments, treating
// "steps", "rollbackSteps", and "rollback" as the same step-container keyword
// so they align during mergeFQN (v0, trie output, and registry FQNs differ).
func extractStepIDPositions(parts []string) []int {
	var positions []int
	for i := range parts {
		if i == 0 {
			continue
		}
		prev := parts[i-1]
		if prev == "steps" || prev == "rollbackSteps" || prev == "rollback" {
			positions = append(positions, i)
		}
	}
	return positions
}

// staticPrefix extracts the static (non-expression) prefix of a path segment.
// e.g., "Linux_Connectivity_<+expr>" → "Linux_Connectivity_"
//
//	"Test_Helm_Chart<+stepGroup.id>" → "Test_Helm_Chart"
//	"plainID" → "plainID"
func staticPrefix(seg string) string {
	idx := strings.Index(seg, "<+")
	if idx < 0 {
		return seg
	}
	return seg[:idx]
}

// resolveStepV1Path resolves the v1 FQN base path for stepID via ResolveStepFQN.
// Unknown/ambiguous references return "" so the caller leaves the segment as-is.
func resolveStepV1Path(stepID string, ctx *matchContext, convContext *ConversionContext) string {
	if convContext == nil || convContext.StepInfoByFQN == nil {
		return ""
	}

	stageRef := ctx.lastStageID
	if stageRef == "" {
		stageRef = convContext.CurrentStageID
	}
	chainRef := splitChain(ctx.lastGroupID)
	aliasKind := AliasAbsolute
	if ctx.inStepGroupAlias {
		aliasKind = AliasGroupRelative
	}
	fqn, info, outcome := convContext.resolveStepFQN(stageRef, stepID, chainRef, aliasKind, true, ctx.parentHops)
	switch outcome {
	case OutcomeExactChain, OutcomeSubsequenceWithGap, OutcomeFlatUniqueFallback:
		if fqn != "" && info != nil {
			return fqn
		}
	}
	return ""
}

// splitChain converts the legacy ">"-joined chain string into a slice. Empty
// input yields an empty slice (not nil-of-one-empty-string).
func splitChain(chain string) []string {
	if chain == "" {
		return nil
	}
	return strings.Split(chain, ">")
}

// stepRefTargetsSpecOrOutput reports whether the step reference at stepIndex
// ultimately targets a "spec" or "output" node (walking past nested steps.<id>
// pairs). Only such targets get FQN substitution; other fields stay structural.
func stepRefTargetsSpecOrOutput(parts []pathPart, stepIndex int) bool {
	pos := stepIndex + 1
	for pos+1 < len(parts) && parts[pos].name == "steps" {
		pos += 2
	}
	if pos < len(parts) {
		f := parts[pos].name
		return f == "spec" || f == "output"
	}
	return false
}

// promoteGroupChain appends capturedID to ctx.lastGroupID if it resolves to a
// step group (resolver called with skipStepGroups=false).
func promoteGroupChain(ctx *matchContext, capturedID string, convContext *ConversionContext) {
	if convContext == nil || convContext.StepInfoByFQN == nil {
		return
	}
	stageRef := ctx.lastStageID
	if stageRef == "" {
		stageRef = convContext.CurrentStageID
	}
	chainRef := splitChain(ctx.lastGroupID)
	aliasKind := AliasAbsolute
	if ctx.inStepGroupAlias {
		aliasKind = AliasGroupRelative
	}
	_, info, outcome := convContext.resolveStepFQN(stageRef, capturedID, chainRef, aliasKind, false, ctx.parentHops)
	if info != nil && info.Type == "StepGroup" &&
		(outcome == OutcomeExactChain || outcome == OutcomeSubsequenceWithGap || outcome == OutcomeFlatUniqueFallback) {
		ctx.lastGroupID = appendGroupChain(ctx.lastGroupID, info.StepID)
	}
}

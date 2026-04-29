package convertexpressions

import (
	"fmt"
)

// NodeBuilder provides a fluent API for building trie paths
type NodeBuilder struct {
	trie *Trie
	node *TrieNode
}

// AddPath starts building a new path from the root
func (t *Trie) AddPath() *NodeBuilder {
	return &NodeBuilder{
		trie: t,
		node: t.root,
	}
}

// AddPathFromID starts building a new path from a node identified by ID
func (t *Trie) AddPathFromID(id string) *NodeBuilder {
	node := t.nodeIndex[id]
	if node == nil {
		fmt.Printf("Warning: node ID '%s' not found in trie\n", id)
		return &NodeBuilder{
			trie: t,
			node: t.root, // Fallback to root
		}
	}
	return &NodeBuilder{
		trie: t,
		node: node,
	}
}

// Node adds or navigates to a child node
// For wildcards, use "*" as the v0Name
func (nb *NodeBuilder) Node(v0Name string) *NodeBuilder {
	if v0Name == "*" {
		// Handle wildcard node
		if nb.node.wildcardChild == nil {
			nb.node.wildcardChild = &TrieNode{
				children:   make(map[string]*TrieNode),
				isWildcard: true,
				v1Name:     "*", // Will be replaced with actual value during matching
			}
		}
		nb.node = nb.node.wildcardChild
	} else {
		// Handle named node
		if nb.node.children[v0Name] == nil {
			nb.node.children[v0Name] = &TrieNode{
				children: make(map[string]*TrieNode),
				v1Name:   v0Name, // Default: same as v0
			}
		}
		nb.node = nb.node.children[v0Name]
	}
	return nb
}

// WithAlias sets an alias for this node (for relative path matching)
func (nb *NodeBuilder) WithAlias(alias string) *NodeBuilder {
	nb.node.alias = alias
	nb.trie.aliasIndex[alias] = append(nb.trie.aliasIndex[alias], nb.node)
	return nb
}

// WithID sets a unique identifier for this node
func (nb *NodeBuilder) WithID(id string) *NodeBuilder {
	nb.node.id = id
	nb.trie.nodeIndex[id] = nb.node
	return nb
}

// WithV1Name sets the v1 output name
// Use "-" to skip this segment in output
func (nb *NodeBuilder) WithV1Name(v1Name string) *NodeBuilder {
	nb.node.v1Name = v1Name
	return nb
}

// AsWildcard marks this node as a wildcard (already set by Node("*"))
func (nb *NodeBuilder) AsWildcard() *NodeBuilder {
	nb.node.isWildcard = true
	return nb
}

// AsArray marks this node as an array node
func (nb *NodeBuilder) AsArray() *NodeBuilder {
	nb.node.isArray = true
	return nb
}

// End marks this node as a terminal rule with a target pattern
func (nb *NodeBuilder) End(target string) *NodeBuilder {
	nb.node.isEnd = true
	nb.node.target = target
	return nb
}

// LinkToNodeByID creates an edge (named or wildcard) to an existing node identified by ID
// If edgeName is "*", creates a wildcard edge; otherwise creates a named edge
func (nb *NodeBuilder) LinkToNodeByID(edgeName string, targetNodeID string) *NodeBuilder {
	targetNode := nb.trie.nodeIndex[targetNodeID]
	if targetNode == nil {
		fmt.Printf("Warning: target node ID '%s' not found in trie\n", targetNodeID)
		return nb
	}
	
	if edgeName == "*" {
		// Create wildcard edge
		nb.node.wildcardChild = targetNode
	} else {
		// Create named edge
		nb.node.children[edgeName] = targetNode
	}
	
	return nb
}
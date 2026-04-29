package convertexpressions

import "strings"

// FindHarnessExprs finds all top-level <+...> expressions in a string,
// handling arbitrary nesting depth via balanced bracket scanning.
// Returns a list of {start, end} index pairs.
func FindHarnessExprs(s string) [][2]int {
	var results [][2]int
	n := len(s)
	for i := 0; i < n-1; i++ {
		if s[i] == '<' && s[i+1] == '+' {
			start := i
			depth := 1
			j := i + 2
			for j < n && depth > 0 {
				if j+1 < n && s[j] == '<' && s[j+1] == '+' {
					depth++
					j += 2
				} else if s[j] == '>' {
					depth--
					j++
				} else {
					j++
				}
			}
			if depth == 0 {
				results = append(results, [2]int{start, j})
				i = j - 1 // skip past this expression
			}
		}
	}
	return results
}

// replaceHarnessExprs replaces all top-level <+...> expressions in s
// using the provided converter function on the inner content (without <+ and >).
func replaceHarnessExprs(s string, converter func(inner string) string) string {
	spans := FindHarnessExprs(s)
	if len(spans) == 0 {
		return s
	}
	var b strings.Builder
	prev := 0
	for _, span := range spans {
		b.WriteString(s[prev:span[0]])
		innerContent := s[span[0]+2 : span[1]-1] // strip <+ and >
		b.WriteString("<+" + converter(innerContent) + ">")
		prev = span[1]
	}
	b.WriteString(s[prev:])
	return b.String()
}

// splitPathSegments splits a dotted path into segments while treating
// nested <+...> expressions (including multi-level nesting) as single opaque segments.
func splitPathSegments(path string) []string {
	var segments []string
	depth := 0
	start := 0

	for i := 0; i < len(path); i++ {
		if i+1 < len(path) && path[i] == '<' && path[i+1] == '+' {
			depth++
			i++ // skip '+'
		} else if path[i] == '>' && depth > 0 {
			depth--
		} else if path[i] == '.' && depth == 0 {
			segments = append(segments, path[start:i])
			start = i + 1
		}
	}

	// Add the last segment
	if start < len(path) {
		segments = append(segments, path[start:])
	}

	return segments
}

func isPathChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

// consumePathSegment consumes a single path segment starting at pos.
// A segment can be a mix of words ([a-zA-Z0-9_]+) and nested <+...> expressions.
// For example: deployment_<+expr1>_<+expr2> is treated as one segment.
// The segment ends at a dot (.) or when no more path chars or expressions follow.
func consumePathSegment(expr string, pos int) int {
	n := len(expr)
	if pos >= n {
		return pos
	}

	startPos := pos

	// Keep consuming as long as we see path chars or nested expressions
	for pos < n {
		if pos+1 < n && expr[pos] == '<' && expr[pos+1] == '+' {
			// Consume nested expression
			depth := 1
			pos += 2
			for pos < n && depth > 0 {
				if pos+1 < n && expr[pos] == '<' && expr[pos+1] == '+' {
					depth++
					pos += 2
				} else if expr[pos] == '>' {
					depth--
					pos++
				} else {
					pos++
				}
			}
		} else if isPathChar(expr[pos]) {
			// Consume word characters
			pos++
		} else if expr[pos] == '[' {
			// Consume array index
			pos++
			for pos < n && expr[pos] != ']' {
				pos++
			}
			if pos < n && expr[pos] == ']' {
				pos++
			}
		} else {
			// Stop at any other character (like '.')
			break
		}
	}

	// If we didn't consume anything, advance at least one position to avoid infinite loop
	if pos == startPos {
		pos++
	}

	return pos
}

// replaceDottedPaths finds dotted path-like segments in expr (e.g. word.word.word)
// while treating nested <+...> expressions as opaque segments, and replaces each
// matched path using the provided replacer function.
func replaceDottedPaths(expr string, replacer func(string) string) string {
	var result strings.Builder
	i := 0
	n := len(expr)

	for i < n {
		// Check if we're at the start of a path-like segment (word char)
		if isPathChar(expr[i]) || (i+1 < n && expr[i] == '<' && expr[i+1] == '+') {
			// Try to consume a dotted path
			pathStart := i
			hasDot := false
			i = consumePathSegment(expr, i)

			for i < n && expr[i] == '.' {
				next := i + 1
				if next < n && (isPathChar(expr[next]) || (next+1 < n && expr[next] == '<' && expr[next+1] == '+')) {
					hasDot = true
					i = next
					i = consumePathSegment(expr, i)
				} else {
					break
				}
			}

			path := expr[pathStart:i]
			if hasDot {
				replaced := replacer(path)
				result.WriteString(replaced)
			} else {
				result.WriteString(path)
			}
		} else {
			result.WriteByte(expr[i])
			i++
		}
	}

	return result.String()
}

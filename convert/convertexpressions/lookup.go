package convertexpressions

import (
	"sort"
	"strings"
)

// AliasKind controls chain scoping: group-relative aliases prepend
// CurrentStepGroupChain to the captured chain; absolute aliases do not.
type AliasKind int

const (
	// AliasAbsolute — stages./pipeline.stages./stage./execution. or none; the
	// chain is a stage-anchored absolute path.
	AliasAbsolute AliasKind = iota
	// AliasGroupRelative — steps. or stepGroup.; the chain is scoped to the
	// call-site's enclosing step group.
	AliasGroupRelative
)

// prefixMatchWithBoundary returns true if registered is a prefix of static
// AND the character following the prefix in static is a delimiter (or static
// equals registered). Prevents over-matching like "prod" → "production".
func prefixMatchWithBoundary(static, registered string) bool {
	if !strings.HasPrefix(static, registered) {
		return false
	}
	if len(static) == len(registered) {
		return true
	}
	switch static[len(registered)] {
	case '_', '-', '.', '<':
		return true
	}
	return false
}

// fuzzyIDMatch reports whether captured matches registered, tolerating
// matrix-expansion suffixes (captured="prod_0" vs "prod") and dynamic <+...>
// substitutions (captured="LZ_<+x>" vs "LZ_Page"). It accepts exact equality, a
// boundary-respecting prefix in either direction, or a pure-dynamic wildcard.
// Boundaries are '_', '-', '.', '<', or end-of-string.
func fuzzyIDMatch(captured, registered string) bool {
	if captured == registered {
		return true
	}
	sp := staticPrefix(captured)
	if sp == "" {
		// Pure-dynamic captured: wildcard match against any registered.
		return true
	}
	// (a) Registered shorter: captured begins with registered + boundary.
	if prefixMatchWithBoundary(sp, registered) {
		return true
	}
	// (b) Captured-static shorter: registered begins with captured-static, with a
	// boundary either at the end of sp or in registered just after it.
	if strings.HasPrefix(registered, sp) && len(registered) > len(sp) {
		// Inner boundary at the end of sp.
		switch sp[len(sp)-1] {
		case '_', '-', '.', '>':
			return true
		}
		// Outer boundary in registered just after sp.
		switch registered[len(sp)] {
		case '_', '-', '.', '<':
			return true
		}
	}
	return false
}

// isLiteralEqual reports whether captured and registered are byte-equal — the
// "best" possible match, used by the literal-equal-beats-fuzzy tie-breaker.
func isLiteralEqual(captured, registered string) bool {
	return captured == registered
}

// ResolveStepFQN finds the v1 FQN of the step (or step group) referenced by
// (stageRef, chainRef, stepRef). aliasKind controls how chainRef composes with
// the call-site's CurrentStepGroupChain; skipStepGroups drops StepGroup
// candidates. Returns the FQN, its info, and an outcome (the pick reason on
// success or the failure mode) that callers use to decide logging/emission.
func (c *ConversionContext) ResolveStepFQN(
	stageRef, stepRef string,
	chainRef []string,
	aliasKind AliasKind,
	skipStepGroups bool,
) (fqn string, info *StepInfoFQN, outcome ResolutionOutcome) {
	return c.resolveStepFQN(stageRef, stepRef, chainRef, aliasKind, skipStepGroups, 0)
}

// resolveStepFQN is the parent-hop-aware implementation behind ResolveStepFQN.
// For group-relative aliases, parentHops trailing elements are dropped from
// CurrentStepGroupChain (one per getParentStepGroup edge) before prepending.
func (c *ConversionContext) resolveStepFQN(
	stageRef, stepRef string,
	chainRef []string,
	aliasKind AliasKind,
	skipStepGroups bool,
	parentHops int,
) (fqn string, info *StepInfoFQN, outcome ResolutionOutcome) {
	if c == nil || c.StepInfoByFQN == nil {
		return "", nil, OutcomeUnknown
	}
	c.EnsureIndexes()

	// Step 0 — build effective chain. For group-relative aliases, prepend the
	// call-site's enclosing step-group chain, after trimming parentHops trailing
	// elements to honour each getParentStepGroup navigation.
	effectiveChain := chainRef
	if aliasKind == AliasGroupRelative && len(c.CurrentStepGroupChain) > 0 {
		base := c.CurrentStepGroupChain
		if parentHops > 0 {
			if parentHops >= len(base) {
				base = nil
			} else {
				base = base[:len(base)-parentHops]
			}
		}
		prepended := make([]string, 0, len(base)+len(chainRef))
		prepended = append(prepended, base...)
		prepended = append(prepended, chainRef...)
		effectiveChain = prepended
	}

	// Step 1 — resolve stage. When neither expression nor context supplies a
	// stage, we still allow the flat-ID fallback below to recover globally
	// unique step IDs (typical for pipeline-level expressions written before
	// the stage was known).
	resolvedStage := c.resolveStageID(stageRef)
	if resolvedStage == "" {
		if len(effectiveChain) == 0 {
			if c.flatStepIDCount[stepRef] == 1 {
				fqn := c.flatStepIDFQN[stepRef]
				info := c.StepInfoByFQN[fqn]
				if !(skipStepGroups && info != nil && info.Type == "StepGroup") {
					return fqn, info, OutcomeFlatUniqueFallback
				}
			}
			if c.flatStepIDCount[stepRef] > 1 {
				return "", nil, OutcomeAmbiguousFlatFallback
			}
		}
		return "", nil, OutcomePureDynamic
	}

	// Step 2 — gather candidates by (stage, stepID). Refuse a pure-dynamic LEAF
	// step ID: wildcard-matching + tie-breaking would invent an arbitrary leaf.
	// (Dynamic chain elements are fine — they only narrow candidates.)
	if staticPrefix(stepRef) == "" {
		return "", nil, OutcomePureDynamic
	}
	candidates := c.gatherStepCandidates(resolvedStage, stepRef, skipStepGroups)
	if len(candidates) == 0 {
		// Step 5 — flat-ID fallback when chainRef is empty AND bare ID is
		// globally unique. This recovers the "step ID alone, no chain" case
		// without making silent picks on duplicates.
		if len(effectiveChain) == 0 && c.flatStepIDCount[stepRef] == 1 {
			fqn := c.flatStepIDFQN[stepRef]
			info := c.StepInfoByFQN[fqn]
			if !(skipStepGroups && info != nil && info.Type == "StepGroup") {
				return fqn, info, OutcomeFlatUniqueFallback
			}
		}
		if len(effectiveChain) == 0 && c.flatStepIDCount[stepRef] > 1 {
			return "", nil, OutcomeAmbiguousFlatFallback
		}
		return "", nil, OutcomeUnknown
	}

	// Step 3 — exact-chain phase.
	exactMatches := exactChainMatches(c.StepInfoByFQN, candidates, effectiveChain)
	if len(exactMatches) == 1 {
		return exactMatches[0], c.StepInfoByFQN[exactMatches[0]], OutcomeExactChain
	}
	if len(exactMatches) > 1 {
		// Tie-breaker 3 (literal-equal beats fuzzy). Tertiary alphabetical.
		picked, tied := tieBreakLiteralThenAlpha(exactMatches, c.StepInfoByFQN, effectiveChain)
		if !tied {
			return picked, c.StepInfoByFQN[picked], OutcomeExactChain
		}
		return "", nil, OutcomeAmbiguousExactChain
	}

	// Step 4 — subsequence fallback.
	subMatches := subsequenceMatches(c.StepInfoByFQN, candidates, effectiveChain)
	if len(subMatches) == 0 {
		// Try flat-ID fallback even when chainRef is non-empty — the chain
		// may simply be wrong. Same hybrid guard: only if globally unique.
		if c.flatStepIDCount[stepRef] == 1 {
			fqn := c.flatStepIDFQN[stepRef]
			info := c.StepInfoByFQN[fqn]
			if !(skipStepGroups && info != nil && info.Type == "StepGroup") {
				return fqn, info, OutcomeFlatUniqueFallback
			}
		}
		return "", nil, OutcomeUnknown
	}
	// Tie-breakers: 1 shallowest, 2 most-compact, 3 literal-equal-beats-fuzzy,
	// 4 alphabetical.
	winners := subMatches
	if len(winners) > 1 {
		winners = filterShallowest(winners, c.StepInfoByFQN)
	}
	if len(winners) > 1 {
		winners = filterMostCompact(winners, c.StepInfoByFQN, effectiveChain)
	}
	if len(winners) > 1 {
		picked, tied := tieBreakLiteralThenAlpha(winners, c.StepInfoByFQN, effectiveChain)
		if tied {
			return "", nil, OutcomeAmbiguousSubsequence
		}
		return picked, c.StepInfoByFQN[picked], OutcomeSubsequenceWithGap
	}
	return winners[0], c.StepInfoByFQN[winners[0]], OutcomeSubsequenceWithGap
}

// resolveStageID picks the stage ID to use. Preference: the expression's
// stageRef (with fuzzyIDMatch against known stage IDs); otherwise fall back to
// CurrentStageID. Returns "" if no usable stage was found (pure-dynamic with
// no context).
func (c *ConversionContext) resolveStageID(stageRef string) string {
	if stageRef == "" {
		return c.CurrentStageID
	}
	// Exact match wins.
	for _, s := range c.stageIDs {
		if s == stageRef {
			return s
		}
	}
	// Fuzzy match: pick the longest registered stage ID matching the captured
	// reference bidirectionally. Longest wins to favour the most specific
	// stage (e.g. "prod_2" should match registered "prod_2" not bare "prod").
	best := ""
	for _, s := range c.stageIDs {
		if fuzzyIDMatch(stageRef, s) {
			if len(s) > len(best) {
				best = s
			}
		}
	}
	if best != "" {
		return best
	}
	// Last resort: fall back to call-site context.
	return c.CurrentStageID
}

// gatherStepCandidates returns the FQNs whose stageID matches resolvedStage and
// whose leaf step ID fuzzy-matches stepRef, optionally filtering step groups.
func (c *ConversionContext) gatherStepCandidates(resolvedStage, stepRef string, skipStepGroups bool) []string {
	byStep := c.stepsByStageStep[resolvedStage]
	if byStep == nil {
		return nil
	}
	var out []string
	// Try exact first to keep the hot path fast.
	if exact, ok := byStep[stepRef]; ok {
		for _, fqn := range exact {
			info := c.StepInfoByFQN[fqn]
			if skipStepGroups && info != nil && info.Type == "StepGroup" {
				continue
			}
			out = append(out, fqn)
		}
		if len(out) > 0 {
			return out
		}
	}
	// Fuzzy step-ID fallback.
	for sid, fqns := range byStep {
		if !fuzzyIDMatch(stepRef, sid) {
			continue
		}
		for _, fqn := range fqns {
			info := c.StepInfoByFQN[fqn]
			if skipStepGroups && info != nil && info.Type == "StepGroup" {
				continue
			}
			out = append(out, fqn)
		}
	}
	return out
}

// exactChainMatches filters candidates to those whose chain length matches
// len(effectiveChain) and whose elements all fuzzy-match in order.
func exactChainMatches(byFQN map[string]*StepInfoFQN, candidates []string, effectiveChain []string) []string {
	var out []string
	for _, fqn := range candidates {
		info := byFQN[fqn]
		if info == nil {
			continue
		}
		if len(info.Chain) != len(effectiveChain) {
			continue
		}
		ok := true
		for i := range effectiveChain {
			if !fuzzyIDMatch(effectiveChain[i], info.Chain[i]) {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, fqn)
		}
	}
	return out
}

// subsequenceMatches filters candidates to those whose chain contains
// effectiveChain as an in-order subsequence (each element of effectiveChain
// fuzzy-matches some chain element, in order). Returns the FQNs only —
// per-candidate match positions are recomputed by tie-breakers as needed.
func subsequenceMatches(byFQN map[string]*StepInfoFQN, candidates []string, effectiveChain []string) []string {
	var out []string
	for _, fqn := range candidates {
		info := byFQN[fqn]
		if info == nil {
			continue
		}
		if isSubsequenceFuzzy(effectiveChain, info.Chain) {
			out = append(out, fqn)
		}
	}
	return out
}

// isSubsequenceFuzzy reports whether short appears as an in-order subsequence
// of long using fuzzyIDMatch element-by-element. Greedy from the left — the
// first viable position consumes the short[i] element.
func isSubsequenceFuzzy(short, long []string) bool {
	if len(short) == 0 {
		return true
	}
	j := 0
	for i := 0; i < len(long) && j < len(short); i++ {
		if fuzzyIDMatch(short[j], long[i]) {
			j++
		}
	}
	return j == len(short)
}

// matchPositions returns the indices in long where short[i] matched, using
// the same greedy strategy as isSubsequenceFuzzy. Returns nil if not a
// subsequence at all.
func matchPositions(short, long []string) []int {
	if len(short) == 0 {
		return []int{}
	}
	out := make([]int, 0, len(short))
	j := 0
	for i := 0; i < len(long) && j < len(short); i++ {
		if fuzzyIDMatch(short[j], long[i]) {
			out = append(out, i)
			j++
		}
	}
	if j == len(short) {
		return out
	}
	return nil
}

// filterShallowest returns the subset of candidates whose chain length equals
// the minimum among candidates.
func filterShallowest(candidates []string, byFQN map[string]*StepInfoFQN) []string {
	minLen := -1
	for _, fqn := range candidates {
		info := byFQN[fqn]
		if info == nil {
			continue
		}
		if minLen < 0 || len(info.Chain) < minLen {
			minLen = len(info.Chain)
		}
	}
	var out []string
	for _, fqn := range candidates {
		info := byFQN[fqn]
		if info == nil || len(info.Chain) != minLen {
			continue
		}
		out = append(out, fqn)
	}
	return out
}

// filterMostCompact selects candidates whose subsequence match has the
// smallest sum of "gap positions" (i.e., the positions in the candidate's
// chain not matched by the expression).
func filterMostCompact(candidates []string, byFQN map[string]*StepInfoFQN, effectiveChain []string) []string {
	type scored struct {
		fqn string
		sum int
	}
	scoredList := make([]scored, 0, len(candidates))
	for _, fqn := range candidates {
		info := byFQN[fqn]
		if info == nil {
			continue
		}
		positions := matchPositions(effectiveChain, info.Chain)
		matched := map[int]struct{}{}
		for _, p := range positions {
			matched[p] = struct{}{}
		}
		sum := 0
		for i := range info.Chain {
			if _, ok := matched[i]; !ok {
				sum += i
			}
		}
		scoredList = append(scoredList, scored{fqn, sum})
	}
	if len(scoredList) == 0 {
		return nil
	}
	minSum := scoredList[0].sum
	for _, s := range scoredList[1:] {
		if s.sum < minSum {
			minSum = s.sum
		}
	}
	var out []string
	for _, s := range scoredList {
		if s.sum == minSum {
			out = append(out, s.fqn)
		}
	}
	return out
}

// tieBreakLiteralThenAlpha picks among candidates by most literal-equal chain
// matches, then alphabetical FQN. The tied flag is reserved for future strict
// policies; the alphabetical fallback currently always disambiguates.
func tieBreakLiteralThenAlpha(candidates []string, byFQN map[string]*StepInfoFQN, effectiveChain []string) (string, bool) {
	if len(candidates) == 0 {
		return "", true
	}
	if len(candidates) == 1 {
		return candidates[0], false
	}
	// Score each candidate by number of literal-equal matches.
	type scored struct {
		fqn   string
		score int
	}
	scoredList := make([]scored, 0, len(candidates))
	for _, fqn := range candidates {
		info := byFQN[fqn]
		if info == nil {
			continue
		}
		positions := matchPositions(effectiveChain, info.Chain)
		if positions == nil {
			// Exact-chain phase guarantees same length; positions may be nil
			// if alignment failed — score 0 in that case.
			scoredList = append(scoredList, scored{fqn, 0})
			continue
		}
		s := 0
		for i, p := range positions {
			if i < len(effectiveChain) && isLiteralEqual(effectiveChain[i], info.Chain[p]) {
				s++
			}
		}
		scoredList = append(scoredList, scored{fqn, s})
	}
	maxScore := -1
	for _, s := range scoredList {
		if s.score > maxScore {
			maxScore = s.score
		}
	}
	var top []string
	for _, s := range scoredList {
		if s.score == maxScore {
			top = append(top, s.fqn)
		}
	}
	if len(top) == 1 {
		return top[0], false
	}
	// Alphabetical FQN as the final, deterministic tie-breaker.
	sort.Strings(top)
	return top[0], false
}

// ParseFQN splits a full step FQN into its stage, ancestor-group chain, and
// leaf step ID, e.g. "pipeline.stages.prod.steps.G1.steps.G2.steps.X" →
// stage="prod", chain=["G1","G2"], step="X".
func ParseFQN(fqn string) (stage string, chain []string, step string, ok bool) {
	parts := strings.Split(fqn, ".")
	// Expect at least: pipeline.stages.<stage>.steps.<step>
	if len(parts) < 5 {
		return "", nil, "", false
	}
	if parts[0] != "pipeline" || parts[1] != "stages" {
		return "", nil, "", false
	}
	stage = parts[2]
	// Remaining must be alternating "steps" / "<id>".
	rest := parts[3:]
	if len(rest)%2 != 0 {
		return "", nil, "", false
	}
	ids := make([]string, 0, len(rest)/2)
	for i := 0; i < len(rest); i += 2 {
		if rest[i] != "steps" {
			return "", nil, "", false
		}
		ids = append(ids, rest[i+1])
	}
	if len(ids) == 0 {
		return "", nil, "", false
	}
	step = ids[len(ids)-1]
	if len(ids) > 1 {
		chain = ids[:len(ids)-1]
	}
	return stage, chain, step, true
}

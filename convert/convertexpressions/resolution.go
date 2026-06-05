package convertexpressions

// ResolutionOutcome categorises the result of ResolveStepFQN. Callers use it
// to choose the emitted path and whether to log a warning.
type ResolutionOutcome int

const (
	OutcomeNone ResolutionOutcome = iota
	// OutcomeExactChain — one candidate's chain matched effectiveChain exactly.
	OutcomeExactChain
	// OutcomeSubsequenceWithGap — matched via in-order subsequence (unmentioned
	// groups); caller warns EXPRESSION_SUBSEQUENCE_GAP_FILLED.
	OutcomeSubsequenceWithGap
	// OutcomeAmbiguousExactChain — multiple exact-chain matches, unresolved.
	OutcomeAmbiguousExactChain
	// OutcomeAmbiguousSubsequence — multiple subsequence matches, unresolved.
	OutcomeAmbiguousSubsequence
	// OutcomeFlatUniqueFallback — empty chainRef, globally unique bare step ID.
	OutcomeFlatUniqueFallback
	// OutcomeAmbiguousFlatFallback — empty chainRef, bare step ID not unique.
	OutcomeAmbiguousFlatFallback
	// OutcomeUnknown — no candidate found.
	OutcomeUnknown
	// OutcomePureDynamic — fully dynamic stage/step ID with no usable context.
	OutcomePureDynamic
)

// String returns a stable identifier for the outcome, suitable for log codes.
func (o ResolutionOutcome) String() string {
	switch o {
	case OutcomeExactChain:
		return "EXACT_CHAIN"
	case OutcomeSubsequenceWithGap:
		return "SUBSEQUENCE_GAP_FILLED"
	case OutcomeAmbiguousExactChain:
		return "AMBIGUOUS_EXACT_CHAIN"
	case OutcomeAmbiguousSubsequence:
		return "AMBIGUOUS_SUBSEQUENCE"
	case OutcomeFlatUniqueFallback:
		return "FLAT_UNIQUE_FALLBACK"
	case OutcomeAmbiguousFlatFallback:
		return "AMBIGUOUS_FLAT_FALLBACK"
	case OutcomeUnknown:
		return "UNKNOWN_STEP"
	case OutcomePureDynamic:
		return "PURELY_DYNAMIC"
	default:
		return "NONE"
	}
}

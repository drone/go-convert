package converthelpers

import (
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// ConvertDelegate converts v0 FlexibleField[[]string] and IncludeInfraSelectors to v1 FlexibleField[Delegate]
func ConvertDelegate(src *flexible.Field[[]string], includeInfraSelectors *flexible.Field[bool]) *flexible.Field[*v1.Delegate] {
	// If source is nil/empty, return nil
	if src == nil {
		return nil
	}

	delegate := &flexible.Field[*v1.Delegate]{}

	// If delegate selectors is an expression, pass it through as expression
	if val, ok := src.AsString(); ok {
		delegate.SetString(val)
		return delegate
	}

	selectors, ok := src.AsStruct()
	if !ok {
		return nil
	}
	// We have concrete delegate selectors ([]string)
	// Now handle includeInfraSelectors
	if includeInfraSelectors != nil {
		if val, ok := includeInfraSelectors.AsStruct(); ok {
			// Boolean value
			if val {
				delegateValue := &v1.Delegate{
					Filter:  selectors,
					Inherit: &flexible.Field[bool]{},
				}
				delegateValue.Inherit.Set(true)
				delegate.Set(delegateValue)
			} else {
				delegate.Set(&v1.Delegate{Filter: selectors})
			}
			return delegate
		}
		if expr, ok := includeInfraSelectors.AsString(); ok {
			// Expression value
			if expr == "<+input>" {
				// <+input> means inherit is also <+input>
				delegateValue := &v1.Delegate{
					Filter:  selectors,
					Inherit: &flexible.Field[bool]{},
				}
				delegateValue.Inherit.SetString("<+input>")
				delegate.Set(delegateValue)
			} else {
				// Other expression: build ternary
				// <+ <+originalExpr> ? "inherit-from-infrastructure" : ["delegate1","delegate2",...] >
				ternary := v1.FormatDelegateExpression(expr, selectors)
				delegate.SetString(ternary)
			}
			return delegate
		}
	}
	// No includeInfraSelectors, just set filter
	delegate.Set(&v1.Delegate{Filter: selectors})
	return delegate
}

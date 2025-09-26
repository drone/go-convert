package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertDelegate converts v0 FlexibleField[[]string] to v1 FlexibleField[Delegate]
func ConvertDelegate(src v0.FlexibleField[[]string]) *v1.FlexibleField[*v1.Delegate] {
	// If source is nil/empty, return nil
	if src.IsNil() {
		return nil
	}

	delegate := &v1.FlexibleField[*v1.Delegate]{}

	// If it's an expression, pass it through as expression
	if src.IsExpression() {
		delegate.SetExpression(src.AsString())
	} else if selectors, ok := src.AsStruct(); ok {
		delegateValue := &v1.Delegate{
			Filter: selectors,
		}
		delegate.Set(delegateValue)
	}

	return delegate
}

package converthelpers

import (
	// v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// ConvertDelegate converts v0 FlexibleField[[]string] to v1 FlexibleField[Delegate]
func ConvertDelegate(src *flexible.Field[[]string]) *flexible.Field[*v1.Delegate] {
	// If source is nil/empty, return nil
	if src == nil {
		return nil
	}

	delegate := &flexible.Field[*v1.Delegate]{}

	// If it's an expression, pass it through as expression
	if val, ok := src.AsString(); ok {
		delegate.SetString(val)
	} else if selectors, ok := src.AsStruct(); ok {
		delegateValue := &v1.Delegate{
			Filter: selectors,
		}
		delegate.Set(delegateValue)
	}

	return delegate
}

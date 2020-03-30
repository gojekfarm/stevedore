package stevedore

// Matchers represents list of conditions
type Matchers []Conditions

// Contains returns true if any of the conditions are met with the context
func (matchers Matchers) Contains(ctx Context) bool {
	predicate := NewPredicateFromContext(ctx)
	for _, conditions := range matchers {
		if predicate.Contains(conditions) {
			return true
		}
	}
	return false
}

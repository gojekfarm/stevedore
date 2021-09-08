package stevedore

// Predicate represents the set of condition to match
type Predicate struct {
	conditions Conditions
}

func (predicate Predicate) add(key, value string) {
	predicate.conditions[key] = value
}

// Contains checks if all the given conditions are present in the predicate
func (predicate Predicate) Contains(conditions Conditions) bool {
	if len(conditions) == 0 {
		return false
	}
	for key, value := range conditions {
		if predicate.conditions[key] != value {
			return false
		}
	}
	return true
}

// NewPredicate returns a predicate based on releaseSpecification and context
func NewPredicate(releaseSpecification ReleaseSpecification, context Context) Predicate {
	predicate := NewPredicateFromContext(context)
	predicate.add(ConditionApplicationName, releaseSpecification.Release.Name)
	return predicate
}

// NewPredicateFromContext returns a predicate based on context
func NewPredicateFromContext(context Context) Predicate {
	return Predicate{conditions: context.Conditions()}
}

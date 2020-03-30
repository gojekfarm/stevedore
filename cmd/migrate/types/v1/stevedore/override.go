package stevedore

import "github.com/gojek/stevedore/pkg/stevedore"

// Overrides represents array of override
type Overrides stevedore.OverrideSpecifications

// Override represents a single override
type Override struct {
	Matches stevedore.Conditions `validate:"criteria"`
	Values  stevedore.Values
}

// Convert converts the override to stevedore.Overrides
func (overrides Overrides) Convert() stevedore.Overrides {
	return stevedore.Overrides{
		Kind:    stevedore.KindStevedoreOverride,
		Version: stevedore.OverrideCurrentVersion,
		Spec:    stevedore.OverrideSpecifications(overrides),
	}
}

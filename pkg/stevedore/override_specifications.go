package stevedore

import (
	"sort"

	"github.com/gojek/stevedore/pkg/merger"
)

// OverrideSpecifications represents array of OverrideSpecification
type OverrideSpecifications []OverrideSpecification

// FilterBy takes a list of overrides and a matcher and gives a list of applicable overrides
func (specs OverrideSpecifications) filterBy(predicate Predicate) OverrideSpecifications {
	matchedOverrides := OverrideSpecifications{}
	for _, override := range specs {
		if predicate.Contains(override.Matches) {
			matchedOverrides = append(matchedOverrides, override)
		}
	}
	return matchedOverrides
}

// Sort returns Overrides based on pre determined order
func (specs OverrideSpecifications) sort(labels Labels) {
	sort.SliceStable(specs, func(i, j int) bool {
		return specs[i].weight(labels) < specs[j].weight(labels)
	})
}

// CollateBy filters overrides by predicate and sort it by its weight
func (specs OverrideSpecifications) CollateBy(predicate Predicate, labels Labels) OverrideSpecifications {
	filteredOverrides := specs.filterBy(predicate)
	filteredOverrides.sort(labels)
	return filteredOverrides
}

// MergeValuesInto merges the values from overrides into the base values
func (specs OverrideSpecifications) MergeValuesInto(base Values) Values {
	values := []map[string]interface{}{base}
	for _, override := range specs {
		values = append(values, override.Values)
	}
	result, err := merger.Merge(values...)
	if err != nil {
		panic(err)
	}
	return result
}

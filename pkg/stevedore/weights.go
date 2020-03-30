package stevedore

import "math"

// Weights represents logical weight
type Weights struct {
	values map[string]int
}

// Sum computes sum of weights for the given knownCriteria
func (weights Weights) Sum(criteria []string) int {
	result := 0
	for _, cr := range criteria {
		if weight, ok := weights.values[cr]; ok {
			result += weight
		}
	}
	return result
}

// NewWeights assigns weight and returns Weight
func NewWeights(criteria []string) Weights {
	values := map[string]int{}
	for index, cr := range criteria {
		values[cr] = int(math.Pow(2, float64(index)))
	}
	return Weights{values: values}
}

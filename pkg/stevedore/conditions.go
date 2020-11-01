package stevedore

import "fmt"

const (
	// ConditionContextName represents condition for context name
	ConditionContextName = "contextName"
	// ConditionApplicationName represents condition for application name
	ConditionApplicationName = "applicationName"
	// ConditionEnvironment represents condition for environment name
	ConditionEnvironment = "environment"
)

var (
	knownCriteria = []string{
		"environmentType",
		ConditionEnvironment,
		"contextType",
		ConditionContextName,
		ConditionApplicationName,
	}
)

// Conditions represents knownCriteria and its corresponding value
type Conditions map[string]string

// Weight returns sum of weight of conditions
func (conditions Conditions) Weight(labels Labels) int {
	var criteria []string
	for key := range conditions {
		criteria = append(criteria, key)
	}
	return labels.Weights().Sum(criteria)
}

// Convert converts given condition to another based on the context
func (conditions Conditions) Convert(using Context) Conditions {
	result := Conditions{}
	target := using.Conditions()
	for key, value := range conditions {
		if targetValue, ok := target[key]; ok {
			result[key] = targetValue
		} else {
			result[key] = value
		}
	}
	return result
}

// Format implements fmt.Formatter. It accepts the formats
// 'y' (yaml)
// 'j' (json)
// '#j' (prettier json).
func (conditions Conditions) Format(f fmt.State, c rune) {
	formatAsJSONOrYaml(f, c, conditions)
}

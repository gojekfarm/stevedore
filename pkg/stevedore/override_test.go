package stevedore_test

import (
	"fmt"
	"testing"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/stretchr/testify/assert"
)

func TestOverrideSpecificationIsValid(t *testing.T) {
	type scenario struct {
		name       string
		valid      bool
		conditions stevedore.Conditions
		error      string
	}

	scenarios := []scenario{
		{name: "should be valid if conditions are empty", valid: true, conditions: stevedore.Conditions{}},
		{name: "should be valid for environmentType", valid: true, conditions: stevedore.Conditions{"environmentType": "staging"}},
		{name: "should be valid for environment", valid: true, conditions: stevedore.Conditions{"environment": "staging"}},
		{name: "should be valid for contextType", valid: true, conditions: stevedore.Conditions{"contextType": "components"}},
		{name: "should be valid for contextName", valid: true, conditions: stevedore.Conditions{"contextName": "components"}},
		{name: "should be valid for applicationName", valid: true, conditions: stevedore.Conditions{"applicationName": "x-service"}},
		{name: "should be invalid for unknown", valid: false, conditions: stevedore.Conditions{"unknown": "invalid"}, error: "Key: 'OverrideSpecification.Matches' Error:Field validation for 'Matches' failed on the 'criteria' tag"},
		{name: "should be invalid even if one field is not valid", valid: false, conditions: stevedore.Conditions{"applicationName": "x-service", "unknown": "invalid"}, error: "Key: 'OverrideSpecification.Matches' Error:Field validation for 'Matches' failed on the 'criteria' tag"},
	}

	for _, scenario := range scenarios {
		t.Run(fmt.Sprintf("%v", scenario.name), func(t *testing.T) {
			override := stevedore.OverrideSpecification{Matches: scenario.conditions}

			err := override.IsValid()

			if scenario.valid {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Equal(t, scenario.error, err.Error())
				}
			}
		})
	}
}

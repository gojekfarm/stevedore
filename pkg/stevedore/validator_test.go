package stevedore_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/stretchr/testify/assert"
)

func TestValidateAndGenerate(t *testing.T) {
	type test struct {
		Name string `yaml:"name" validate:"required"`
		Age  string `yaml:"age" validate:"required"`
	}

	t.Run("should validate input and return error", func(t *testing.T) {
		input := `
name:
age: 21`
		actual := test{}
		err := stevedore.ValidateAndGenerate(strings.NewReader(input), &actual)

		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "Key: 'test.Name' Error:Field validation for 'Name' failed on the 'required' tag")
		}
	})

	t.Run("should validate input and populate object", func(t *testing.T) {
		input := `
name: tom
age: 21`
		actual := test{}
		err := stevedore.ValidateAndGenerate(strings.NewReader(input), &actual)

		assert.NoError(t, err)

		assert.Equal(t, "tom", actual.Name)
		assert.Equal(t, "21", actual.Age)
	})

}

func TestValidate(t *testing.T) {
	type test struct {
		Name string `yaml:"name" validate:"required"`
		Age  string `yaml:"age" validate:"required"`
	}

	t.Run("should validate input and return error", func(t *testing.T) {
		actual := test{Age: "21"}

		err := stevedore.Validate(actual)

		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "Key: 'test.Name' Error:Field validation for 'Name' failed on the 'required' tag")
		}
	})

	t.Run("should validate input and populate object", func(t *testing.T) {
		actual := test{Name: "tom", Age: "21"}
		err := stevedore.Validate(actual)

		assert.NoError(t, err)
	})
}

func TestValidateCriteria(t *testing.T) {
	type test struct {
		Matches stevedore.Conditions `validate:"criteria"`
	}

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
		{name: "should be valid for contextType", valid: true, conditions: stevedore.Conditions{"contextType": "staging"}},
		{name: "should be valid for contextName", valid: true, conditions: stevedore.Conditions{"contextName": "staging"}},
		{name: "should be valid for applicationName", valid: true, conditions: stevedore.Conditions{"applicationName": "staging"}},
		{name: "should be invalid for unknown", valid: false, conditions: stevedore.Conditions{"unknown": "invalid"}, error: "Key: 'test.Matches' Error:Field validation for 'Matches' failed on the 'criteria' tag"},
		{name: "should be invalid even if one of the key is invalid", valid: false, conditions: stevedore.Conditions{"applicationName": "stevedore", "unknown": "invalid"}, error: "Key: 'test.Matches' Error:Field validation for 'Matches' failed on the 'criteria' tag"},
	}

	for _, scenario := range scenarios {
		t.Run(fmt.Sprintf("%v", scenario.name), func(t *testing.T) {
			x := test{Matches: scenario.conditions}

			err := stevedore.Validate(x)

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

func TestValidateAny(t *testing.T) {
	t.Run("should return true if validator matches with any of the provided value", func(t *testing.T) {
		type test struct {
			Color string `yaml:"name" validate:"any=red/green/blue"`
		}

		actual := test{Color: "red"}

		err := stevedore.Validate(actual)

		assert.NoError(t, err)
	})

	t.Run("should return false if validator doesn't matches with any of the provided value", func(t *testing.T) {
		type test struct {
			Color string `yaml:"name" validate:"any=red/green/blue"`
		}

		actual := test{Color: "yellow"}

		err := stevedore.Validate(actual)

		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "Key: 'test.Color' Error:Field validation for 'Color' failed on the 'any' tag")
		}
	})
}

func TestValidateComponent(t *testing.T) {
	t.Run("should be valid if chart name is provided", func(t *testing.T) {
		release := stevedore.Release{Name: "x-service", Namespace: "default", Chart: "chart/example"}

		err := stevedore.Validate(release)
		assert.NoError(t, err)
	})

	t.Run("should be valid if dependencies of chart spec is provided", func(t *testing.T) {
		release := stevedore.Release{Name: "x-service", Namespace: "default", ChartSpec: stevedore.ChartSpec{Name: "x-service-dependencies", Dependencies: stevedore.Dependencies{{Name: "dependencyA"}}}}

		err := stevedore.Validate(release)
		assert.NoError(t, err)
	})

	t.Run("should be invalid if chart name and chart spec is provided", func(t *testing.T) {
		release := stevedore.Release{Name: "x-service", Namespace: "default", Chart: "chart/example", ChartSpec: stevedore.ChartSpec{Name: "x-service-dependencies", Dependencies: stevedore.Dependencies{{Name: "dependencyA"}}}}

		err := stevedore.Validate(release)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "Key: 'Release.Chart' Error:Field validation for 'Chart' failed on the 'EitherChartOrChartSpec' tag")
			assert.Contains(t, err.Error(), "Release.ChartSpec' Error:Field validation for 'ChartSpec' failed on the 'EitherChartOrChartSpec' tag")
		}
	})

	t.Run("should be invalid if chart name and chart spec is not provided", func(t *testing.T) {
		release := stevedore.Release{Name: "x-service", Namespace: "default"}

		err := stevedore.Validate(release)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "Key: 'Release.Chart' Error:Field validation for 'Chart' failed on the 'EitherChartOrChartSpec' tag")
			assert.Contains(t, err.Error(), "Release.ChartSpec' Error:Field validation for 'ChartSpec' failed on the 'EitherChartOrChartSpec' tag")
		}
	})
}

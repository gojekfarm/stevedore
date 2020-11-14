package stevedore_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/stretchr/testify/assert"
)

func TestValuesFlatMap(t *testing.T) {
	values := stevedore.Values{
		"student": map[string]interface{}{
			"name": map[string]interface{}{
				"first": "tom",
				"last":  "jerry",
			},
			"marks": map[interface{}]interface{}{
				"english": 10,
				"tamil":   20,
				"maths":   35,
			},
			"countriesVisited": []string{"india", "japan", "usa"},
			"studiedSchools": []map[interface{}]interface{}{
				{"name": "public school", "place": "chennai"},
				{"name": "private school", "place": "delhi"},
			},
			"scholarshipsReceived": []map[string]interface{}{
				{"kind": "Full tuition fees waiver", "for": "First rank in University"},
				{"kind": "Support for Higher Studies", "for": "Impacting Research paper on Climate Change"},
			},
		},
	}

	expected := map[string]string{
		"student.name.first":                  "tom",
		"student.name.last":                   "jerry",
		"student.marks.english":               "10",
		"student.marks.tamil":                 "20",
		"student.marks.maths":                 "35",
		"student.countriesVisited":            "[india japan usa]",
		"student.studiedSchools.count":        "2",
		"student.studiedSchools.0.name":       "public school",
		"student.studiedSchools.0.place":      "chennai",
		"student.studiedSchools.1.name":       "private school",
		"student.studiedSchools.1.place":      "delhi",
		"student.scholarshipsReceived.count":  "2",
		"student.scholarshipsReceived.0.kind": "Full tuition fees waiver",
		"student.scholarshipsReceived.0.for":  "First rank in University",
		"student.scholarshipsReceived.1.kind": "Support for Higher Studies",
		"student.scholarshipsReceived.1.for":  "Impacting Research paper on Climate Change",
	}

	actual := values.FlatMap()

	if !cmp.Equal(expected, actual) {
		assert.Fail(t, cmp.Diff(expected, actual))
	}
}

func TestValuesMergeWith(t *testing.T) {
	t.Run("should return base values merged with overrides", func(t *testing.T) {
		baseValues := stevedore.Values{
			"key1": "baseValueForK1",
			"key2": "baseValueForK2",
		}

		overrides := stevedore.Overrides{
			Spec: stevedore.OverrideSpecifications{
				{Values: stevedore.Values{"key2": "firstOverrideForK2"}},
				{Values: stevedore.Values{"key1": "secondOverrideForK1"}},
				{Values: stevedore.Values{"key2": "thirdOverrideForK2"}},
				{Values: stevedore.Values{"key3": "fourthOverrideForK3"}},
				{Values: stevedore.Values{"key1": "fifthOverrideForK1"}},
			},
		}

		expected := stevedore.Values{
			"key1": "fifthOverrideForK1",
			"key2": "thirdOverrideForK2",
			"key3": "fourthOverrideForK3",
		}

		actual := baseValues.MergeWith(overrides)

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Actual: %#v did not match \nExpected: %#v", actual, expected)
		}
	})
}

func TestValuesReplace(t *testing.T) {
	t.Run("should replace all placeholders", func(t *testing.T) {
		values := stevedore.Values{
			"name":        "${NAME}",
			"environment": map[interface{}]interface{}{"name": "${NAME}", "type": "stg-${TYPE}"},
		}

		expected := stevedore.Values{
			"name":        "x-service",
			"environment": map[interface{}]interface{}{"name": "x-service", "type": "stg-type"},
		}
		expectedUsedSubstitute := stevedore.Substitute{"NAME": "x-service", "TYPE": "type"}

		actual, usedSubstitute, err := values.Replace(stevedore.Substitute{"NAME": "x-service", "TYPE": "type", "DB_NAME": "database"})

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, expected, actual)
		assert.Equal(t, expectedUsedSubstitute, usedSubstitute)
	})

	t.Run("should not replace placeholders", func(t *testing.T) {
		values := stevedore.Values{
			"name": "${NAME}",
			"type": "stg-${TYPE}",
			"url":  "${URL}",
			"env":  "staging",
		}
		actual, usedSubstitute, err := values.Replace(stevedore.Substitute{"NAME": "x-service"})

		if assert.NotNil(t, err) {
			assert.Equal(t, "Unable to replace 2 variable(s):\n\t1. ${TYPE}\n\t2. ${URL}", err.Error())
		}
		assert.NotNil(t, actual)
		assert.Empty(t, usedSubstitute)
		assert.Equal(t, values, actual)
	})

	t.Run("should not replace if provided value is empty", func(t *testing.T) {
		values := stevedore.Values{
			"name": "${NAME}",
			"type": "stg-${TYPE}",
			"url":  "${URL}",
			"env":  "staging",
		}
		actual, usedSubstitute, err := values.Replace(stevedore.Substitute{"NAME": "x-service", "TYPE": "type", "URL": ""})

		if assert.NotNil(t, err) {
			assert.Equal(t, "Unable to replace 1 variable(s):\n\t1. ${URL}", err.Error())
		}
		assert.NotNil(t, actual)
		assert.Empty(t, usedSubstitute)
		assert.Equal(t, values, actual)
	})
}

func TestValuesToYAML(t *testing.T) {
	t.Run("should return YAML as string", func(t *testing.T) {
		values := stevedore.Values{
			"name": "SomeName",
			"type": map[string]string{"nestedType": "nestedValue"},
			"url":  "someURL",
		}
		expectedYAMLContent := `name: SomeName
type:
  nestedType: nestedValue
url: someURL
`

		yamlContent, err := values.ToYAML()

		assert.NoError(t, err)
		assert.Equal(t, expectedYAMLContent, yamlContent)
	})
}

func TestValuesVariables(t *testing.T) {
	t.Run("should get all the env variables", func(t *testing.T) {
		values := stevedore.Values{
			"name": "SomeName",
			"type": map[string]string{"nestedType": "${NESTED_VALUE}"},
			"url":  "${someUrl}",
			"host": "${HoSt}",
		}

		vars, err := values.Variables()

		assert.Equal(t, []string{"HoSt", "NESTED_VALUE", "someUrl"}, vars)
		assert.Nil(t, err)
	})
}

package stevedore

import (
	"bytes"
	"fmt"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"regexp"
)

var variablePattern *regexp.Regexp

func init() {
	pattern, err := regexp.Compile(`\${(\w*)}`)
	if err != nil {
		panic(fmt.Errorf("[values init] %v", err))
	}
	variablePattern = pattern
}

// Values to wrap all helm values that are passed to deploy the release specification
type Values map[string]interface{}

// MergeWith merges values with the given overrides
func (values Values) MergeWith(overrides Overrides) Values {
	return overrides.MergeValuesInto(values)
}

// Replace replace all placeholder text with given value
func (values Values) Replace(substitute Substitute) (Values, Substitute, error) {
	substitutes := Substitute{}
	valueStr, err := values.toString()

	if err != nil {
		return values, substitutes, err
	}

	resultStr, err := substitute.Perform(valueStr)
	if err != nil {
		return values, substitutes, err
	}

	var result = Values{}
	err = yaml.Unmarshal([]byte(resultStr), &result)
	if err != nil {
		return values, substitutes, err
	}

	variables, err := values.Variables()

	if err != nil {
		return result, substitutes, nil
	}

	errors := SubstituteError{}
	for _, variable := range variables {
		if value, ok := substitute[variable]; ok {
			if value == "" {
				errors = append(errors, fmt.Sprintf("${%s}", variable))
			}
			substitutes[variable] = value
		}
	}

	if len(errors) != 0 {
		return values, Substitute{}, errors
	}

	return result, substitutes, nil
}

func (values Values) toString() (string, error) {
	data, err := yaml.Marshal(values)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToYAML serializes values into a YAML string
func (values Values) ToYAML() (string, error) {
	var buffer bytes.Buffer

	if err := yaml.NewEncoder(&buffer).Encode(values); err != nil {
		return "", fmt.Errorf("error serializing values: %v", err)
	}

	return buffer.String(), nil
}

// Variables returns all placeholders
func (values Values) Variables() ([]string, error) {
	valueStr, err := values.toString()

	if err != nil {
		return nil, err
	}

	placeHolders := placeholderPattern.FindAllString(valueStr, -1)
	var result []string
	for _, placeHolder := range placeHolders {
		matchGroups := variablePattern.FindStringSubmatch(placeHolder)
		if len(matchGroups) == 2 {
			result = append(result, matchGroups[1])
		}
	}
	return result, nil
}

// FlatMap return the values as map[string]string,
// nested key will be converted to dotted key and leaf node values are converted to string
// if leaf node is not a primitive type, then value will be sha256sum
func (values Values) FlatMap() map[string]string {
	return merge(map[string]string{}, flatten("", values))
}

func flatten(parentKey string, source map[string]interface{}) map[string]string {
	result := map[string]string{}

	for key, value := range source {
		resultKey := formatKey(parentKey, key)
		switch value := value.(type) {
		case string:
			result[resultKey] = value
		case map[string]interface{}:
			flattenedResult := flatten(resultKey, value)
			merge(result, flattenedResult)
		case map[interface{}]interface{}:
			convertedMap := convert(value)
			merge(result, flatten(resultKey, convertedMap))
		case []map[string]interface{}:
			items := value
			addTo(result, items, resultKey)
		case []map[interface{}]interface{}:
			items := convertGenericArray(value)
			addTo(result, items, resultKey)
		default:
			result[resultKey] = fmt.Sprintf("%v", value)
		}
	}
	return result
}

func convertGenericArray(source []map[interface{}]interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for _, item := range source {
		result = append(result, convert(item))
	}
	return result
}

func addTo(result map[string]string, items []map[string]interface{}, key string) {
	result[fmt.Sprintf("%s.count", key)] = fmt.Sprintf("%v", len(items))
	for index, item := range items {
		merge(result, flatten(formatKey(key, fmt.Sprintf("%v", index)), item))
	}
}

func convert(source map[interface{}]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	for key, value := range source {
		resultKey := fmt.Sprintf("%v", key)
		result[resultKey] = value
	}
	return result
}

func formatKey(parentKey, key string) string {
	if parentKey != "" {
		return fmt.Sprintf("%s.%s", parentKey, key)
	}
	return key
}

func merge(source, another map[string]string) map[string]string {
	err := mergo.Merge(&source, another, mergo.WithAppendSlice, mergo.WithOverride)
	if err != nil {
		source["error"] = err.Error()
	}
	return source
}

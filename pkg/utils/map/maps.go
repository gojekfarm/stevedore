package maputils

import "reflect"

// ExtractCommon returns common map between two maps
func ExtractCommon(a, b map[string]interface{}) map[string]interface{} {
	common := make(map[string]interface{})
	for key, aValue := range a {
		bValue, ok := b[key]
		if !ok || reflect.TypeOf(aValue).Kind() != reflect.TypeOf(bValue).Kind() {
			continue
		}

		if reflect.TypeOf(aValue).Kind() == reflect.Map {
			extractedCommon := ExtractCommon(aValue.(map[string]interface{}), bValue.(map[string]interface{}))
			if len(extractedCommon) > 0 {
				common[key] = extractedCommon
			}
		} else if bValue == aValue {
			common[key] = aValue
		}
	}

	return common
}

// Keys gives keys of a given map
func Keys(data map[string]interface{}) []string {
	keys := make([]string, 0, len(data))

	for key := range data {
		keys = append(keys, key)
	}

	return keys
}

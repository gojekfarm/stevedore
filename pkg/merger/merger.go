package merger

import (
	"github.com/imdario/mergo"
)

// Merge array of hashes
func Merge(hashes ...map[string]interface{}) (map[string]interface{}, error) {

	finalMap := make(map[string]interface{})

	for _, hash := range hashes {
		if err := mergo.Merge(&finalMap, hash, mergo.WithOverride); err != nil {
			return nil, err
		}
	}
	return finalMap, nil
}

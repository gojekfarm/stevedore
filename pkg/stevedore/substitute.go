package stevedore

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	stringutils "github.com/gojek/stevedore/pkg/utils/string"
	"github.com/imdario/mergo"
)

var placeholderPattern *regexp.Regexp

func init() {
	pattern, err := regexp.Compile(`(\${\w*})`)
	if err != nil {
		panic(fmt.Errorf("[values init] %v", err))
	}
	placeholderPattern = pattern
}

// Substitute holds key value pair for substitution
type Substitute map[string]interface{}

// SubstituteError represents substitution errors
type SubstituteError []string

// Error returns substitution errors as string
func (err SubstituteError) Error() string {
	buff := bytes.NewBufferString(fmt.Sprintf("Unable to replace %d variable(s):", len(err)))
	for index, err := range err {
		buff.WriteString(fmt.Sprintf("\n\t%d. %s", index+1, err))
	}
	return strings.TrimSpace(buff.String())
}

func (sub Substitute) get(key string, isInterpolated bool) string {
	if value, ok := sub[key]; ok {
		if _, ok = value.(string); ok && !isInterpolated {
			return fmt.Sprintf("'%s'", value)
		}
		return fmt.Sprintf("%v", value)
	}
	return fmt.Sprintf("${%s}", key)
}

// Perform will validate and substitute variables
func (sub Substitute) Perform(str string) (string, error) {
	resultStr := stringutils.Expand(str, sub.get)
	placeHolders := placeholderPattern.FindAllString(resultStr, -1)
	if len(placeHolders) != 0 {
		return str, SubstituteError(placeHolders)
	}
	return resultStr, nil
}

// Merge merges the substitutes and returns the result
func (sub Substitute) Merge(dest ...Substitute) (Substitute, error) {
	intermediate := []Substitute{sub}
	result := Substitute{}
	for _, hash := range append(intermediate, dest...) {
		if err := mergo.Merge(&result, hash, mergo.WithAppendSlice, mergo.WithOverride); err != nil {
			return nil, err
		}
	}
	return result, nil
}

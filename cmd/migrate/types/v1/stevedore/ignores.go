package stevedore

import "github.com/gojek/stevedore/pkg/stevedore"

// Ignores is a list of Ignore
type Ignores []Ignore

// Convert convert ignores to newer stevedore ignores format
func (ignores Ignores) Convert() stevedore.Ignores {
	result := stevedore.Ignores{}
	for _, ignore := range ignores {
		result = append(result, ignore.Convert())
	}
	return result
}

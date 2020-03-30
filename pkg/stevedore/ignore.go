package stevedore

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

// Ignore represents a single Ignore
type Ignore struct {
	Matches  Conditions      `yaml:"matches" json:"matches" validate:"criteria"`
	Releases IgnoredReleases `yaml:"releases" json:"releases"`
}

// Ignores is a list of Ignore
type Ignores []Ignore

// IsValid validates the context and returns error if any
func (ignore Ignore) IsValid() error {
	return validate.Struct(ignore)
}

// NewIgnores to Validate the Stevedore manifest configuration
func NewIgnores(reader io.Reader) (Ignores, error) {
	ignores := Ignores{}
	err := yaml.NewDecoder(reader).Decode(&ignores)
	if err != nil {
		return nil, fmt.Errorf("[NewIgnores] error when validating from file:\n%v", err)
	}

	for _, ignore := range ignores {
		if err := ignore.IsValid(); err != nil {
			return nil, err
		}
	}

	return ignores, nil
}

// Filter takes a list of Ignores and a matcher and gives a list of applicable Ignores
func (ignores Ignores) Filter(predicate Predicate) IgnoredReleases {
	matchedIgnores := Ignores{}
	for _, ignore := range ignores {
		if predicate.Contains(ignore.Matches) {
			matchedIgnores = append(matchedIgnores, ignore)
		}
	}
	return matchedIgnores.components()
}

func (ignores Ignores) components() IgnoredReleases {
	var result IgnoredReleases

	for _, ignore := range ignores {
		result = append(result, ignore.Releases...)
	}
	return result
}

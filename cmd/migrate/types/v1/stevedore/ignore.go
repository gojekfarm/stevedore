package stevedore

import "github.com/gojek/stevedore/pkg/stevedore"

// Ignore represents a single Ignore
type Ignore struct {
	Matches    stevedore.Conditions      `yaml:"matches" json:"matches" validate:"criteria"`
	Components stevedore.IgnoredReleases `yaml:"components" json:"components"`
}

// Convert convert ignore to newer stevedore ignore format
func (ignore Ignore) Convert() stevedore.Ignore {
	return stevedore.Ignore{
		Matches:  ignore.Matches,
		Releases: ignore.Components,
	}
}

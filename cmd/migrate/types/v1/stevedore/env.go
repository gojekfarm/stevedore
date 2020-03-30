package stevedore

import "github.com/gojek/stevedore/pkg/stevedore"

// Env represents all the env file
type Env struct {
	Matches stevedore.Conditions `validate:"criteria"`
	Values  stevedore.Substitute `yaml:"env"`
}

// Convert convert env to newer stevedore env format
func (env Env) Convert() stevedore.EnvSpecification {
	return stevedore.EnvSpecification{
		Matches: env.Matches,
		Values:  env.Values,
	}
}

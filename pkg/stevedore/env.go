package stevedore

import (
	"fmt"
	"io"
	"sort"

	"gopkg.in/yaml.v2"
)

// Env represents all the env file
type Env struct {
	Kind    string            `yaml:"kind" json:"kind" validate:"required"`
	Version string            `yaml:"version" json:"version" validate:"required"`
	Spec    EnvSpecifications `yaml:"spec" json:"spec" validate:"required,dive"`
}

// Envs is collection of Env
type Envs []Env

// IsValid validates the context and returns error if any
func (env Env) IsValid() error {
	return validate.Struct(env)
}

// EnvSpecification represents specification for env
type EnvSpecification struct {
	Matches Conditions `validate:"criteria" yaml:"matches" json:"matches"`
	Values  Substitute `yaml:"env" json:"env"`
}

// IsApplicableFor returns true if Manifest is applicable for the given environment
func (env EnvSpecification) IsApplicableFor(context Context) bool {
	predicate := NewPredicateFromContext(context)
	return predicate.Contains(env.Matches)
}

// EnvSpecifications is collection of Env
type EnvSpecifications []EnvSpecification

// weight returns the weight of the envs
func (env EnvSpecification) weight(labels Labels) int {
	return env.Matches.Weight(labels)
}

// Sort will sort the envs based on the pre-determined order
func (envs EnvSpecifications) Sort(labels Labels) {
	sort.SliceStable(envs, func(i, j int) bool {
		return envs[i].weight(labels) < envs[j].weight(labels)
	})
}

// Format implements fmt.Formatter. It accepts the formats
// 'y' (yaml)
// 'j' (json)
// '#j' (prettier json).
func (env Env) Format(f fmt.State, c rune) {
	formatAsJSONOrYaml(f, c, env)
}

// NewEnv to Validate the Stevedore env configuration
func NewEnv(reader io.Reader) (Env, error) {
	env := Env{}
	err := yaml.NewDecoder(reader).Decode(&env)
	if err != nil {
		return Env{}, fmt.Errorf("[NewEnv] error when validating from file:\n%v", err)
	}

	err = env.IsValid()
	if err != nil {
		return Env{}, fmt.Errorf("[NewEnv] error when validating from file:\n%v", err)
	}

	return env, nil
}

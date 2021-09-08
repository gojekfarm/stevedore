package stevedore

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

// Overrides is a list of override
type Overrides struct {
	Kind    string                 `yaml:"kind" json:"kind" validate:"required"`
	Version string                 `yaml:"version" json:"version" validate:"required"`
	Spec    OverrideSpecifications `yaml:"spec" json:"spec" validate:"required"`
}

// EmptyOverrides returns an empty Override with Kind & Version populated
func EmptyOverrides() Overrides {
	return Overrides{
		Kind:    KindStevedoreOverride,
		Version: OverrideCurrentVersion,
	}
}

// Format implements fmt.Formatter. It accepts the formats
// 'y' (yaml)
// 'j' (json)
// '#j' (prettier json).
func (overrides Overrides) Format(f fmt.State, c rune) {
	formatAsJSONOrYaml(f, c, overrides)
}

// NewOverrides to Validate the Stevedore manifest configuration
func NewOverrides(reader io.Reader) (Overrides, error) {
	overrides := Overrides{}
	err := yaml.NewDecoder(reader).Decode(&overrides)
	if err != nil {
		return Overrides{}, fmt.Errorf("[NewOverrides] error when validating from file:\n%v", err)
	}

	err = validate.Struct(overrides)
	if err != nil {
		return Overrides{}, err
	}

	for _, override := range overrides.Spec {
		if err := override.IsValid(); err != nil {
			return Overrides{}, err
		}
	}
	return overrides, nil
}

// CollateBy filters overrides by predicate and sort it by its weight
func (overrides Overrides) CollateBy(predicate Predicate, labels Labels) Overrides {
	overrides.Spec = overrides.Spec.CollateBy(predicate, labels)
	return overrides
}

// MergeValuesInto merges the values from overrides into the base values
func (overrides Overrides) MergeValuesInto(base Values) Values {
	return overrides.Spec.MergeValuesInto(base)
}

package stevedore

// OverrideSpecification represents a single override
type OverrideSpecification struct {
	FileName string     `yaml:"-" json:"-"`
	Matches  Conditions `validate:"criteria" yaml:"matches" json:"matches"`
	Values   Values     `yaml:"values" json:"values"`
}

// IsValid validates the context and returns error if any
func (spec OverrideSpecification) IsValid() error {
	return validate.Struct(spec)
}

func (spec OverrideSpecification) weight(labels Labels) int {
	return spec.Matches.Weight(labels)
}

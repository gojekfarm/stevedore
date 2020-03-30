package stevedore

import (
	"io"
	"strings"

	"gopkg.in/yaml.v2"

	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

// ValidateAndGenerate validates if the input conforms to requirement and returns error if any
func ValidateAndGenerate(reader io.Reader, out interface{}) error {
	if err := yaml.NewDecoder(reader).Decode(out); err != nil {
		return err
	}

	return Validate(out)
}

// Validate validates if the type conforms to the requirement and returns error if any
func Validate(v interface{}) error {
	return validate.Struct(v)
}

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("any", validateAny)
	_ = validate.RegisterValidation("criteria", validateCriteria)
	validate.RegisterStructValidation(releaseValidation, Release{})
}

func validateAny(fl validator.FieldLevel) bool {
	matchWith := fl.Param()
	value := fl.Field().String()
	for _, c := range strings.Split(matchWith, "/") {
		if c == value {
			return true
		}
	}
	return false
}

func contains(matchers []string, value string) bool {
	for _, matcher := range matchers {
		if matcher == value {
			return true
		}
	}
	return false
}

func validateCriteria(fl validator.FieldLevel) bool {
	values := fl.Field().MapKeys()

	if len(values) == 0 {
		return true
	}

	for _, value := range values {
		if !contains(knownCriteria, value.String()) {
			return false
		}
	}

	return true
}

func releaseValidation(sl validator.StructLevel) {
	release, ok := sl.Current().Interface().(Release)
	if !ok {
		return
	}

	if (release.Chart == "" && len(release.ChartSpec.Dependencies) == 0) || (release.Chart != "" && len(release.ChartSpec.Dependencies) > 0) {
		sl.ReportError(release.Chart, "Chart", "chart", "EitherChartOrChartSpec", "")
		sl.ReportError(release.ChartSpec, "ChartSpec", "chartSpec", "EitherChartOrChartSpec", "")
	}
	if release.Chart == "" {
		if release.ChartSpec.Name == "" {
			sl.ReportError(release.ChartSpec, "ChartSpec", "Name", "Name is required in chartSpec", "")
		}
		if len(release.ChartSpec.Dependencies) == 0 {
			sl.ReportError(release.ChartSpec, "ChartSpec", "Dependencies", "Dependencies in chartSpec cannot be empty", "")
		}
	}
}

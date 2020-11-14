package lint

import (
	"fmt"
	"strings"

	"github.com/gojek/stevedore/pkg/stevedore"
)

// Error represents lint error
type Error struct {
	Description string
	Matches     stevedore.Conditions
}

// Error returns the underlying error
func (err Error) Error() string {
	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("found %s in overrides:", err.Description))
	msg.WriteString(fmt.Sprintf("\n%y", err.Matches))
	return msg.String()
}

// Errors represents collection of lint Error
type Errors []Error

// Error returns the underlying error
func (errs Errors) Error() string {
	msg := strings.Builder{}

	msg.WriteString(fmt.Sprintf("found %d issue(s) in overrides:\n", len(errs)))
	for index, err := range errs {
		msg.WriteString(fmt.Sprintf("\t%d. %s", index+1, err.Error()))
	}
	return msg.String()
}

// Lint is used to lint overrides
func Lint(overrides stevedore.Overrides) error {
	errors := Errors{}
	matchesMap := make(map[string]struct{}, len(overrides.Spec))
	for _, overrideSpec := range overrides.Spec {
		matches := overrideSpec.Matches
		matchesStr := fmt.Sprintf("%y", matches)
		_, matchesExistsAlready := matchesMap[matchesStr]
		if matchesExistsAlready {
			errors = append(errors, Error{Description: "duplicate matches", Matches: matches})
			continue
		}
		matchesMap[matchesStr] = struct{}{}
	}

	if len(errors) != 0 {
		return errors
	}
	return nil
}

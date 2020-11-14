package lint_test

import (
	"testing"

	"github.com/gojek/stevedore/pkg/lint"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/stretchr/testify/assert"
)

func TestLint(t *testing.T) {
	appName := "abc-service"
	t.Run("lint should fail when overrides has duplicate matches", func(t *testing.T) {
		contextName := "integration-cluster"
		overridesSpecs := stevedore.OverrideSpecifications{
			{
				Matches: stevedore.Conditions{
					stevedore.ConditionApplicationName: appName,
					stevedore.ConditionContextName:     contextName,
				},
				Values: stevedore.Values{
					"image": "abc-service:1.0.0",
				},
			},
			{
				Matches: stevedore.Conditions{
					stevedore.ConditionContextName:     contextName,
					stevedore.ConditionApplicationName: appName,
				},
				Values: stevedore.Values{
					"imagePullPolicy": "Always",
				},
			},
			{
				Matches: stevedore.Conditions{
					stevedore.ConditionApplicationName: appName,
					stevedore.ConditionContextName:     contextName,
				},
				Values: stevedore.Values{
					"replicas": "2",
				},
			},
		}
		overrides := stevedore.EmptyOverrides()
		overrides.Spec = overridesSpecs
		expectedError := `found 2 issue(s) in overrides:
	1. found duplicate matches in overrides:
applicationName: abc-service
contextName: integration-cluster
	2. found duplicate matches in overrides:
applicationName: abc-service
contextName: integration-cluster
`
		err := lint.Lint(overrides)

		assert.EqualError(t, err, expectedError)
	})

	t.Run("lint should pass when overrides has no duplicate matches", func(t *testing.T) {
		overridesSpecs := stevedore.OverrideSpecifications{
			{
				Matches: stevedore.Conditions{
					stevedore.ConditionApplicationName: appName,
				},
				Values: stevedore.Values{
					"image": "abc-service:1.0.0",
				},
			},
			{
				Matches: stevedore.Conditions{
					stevedore.ConditionApplicationName: appName,
					stevedore.ConditionContextName:     "production-cluster-context",
				},
				Values: stevedore.Values{
					"imagePullPolicy": "Always",
				},
			},
		}
		overrides := stevedore.EmptyOverrides()
		overrides.Spec = overridesSpecs

		errorMessages := lint.Lint(overrides)

		assert.Nil(t, errorMessages)
	})
}

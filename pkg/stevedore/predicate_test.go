package stevedore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPredicate(t *testing.T) {
	t.Run("should create predicate given release specification and context", func(t *testing.T) {
		app := ReleaseSpecification{
			Release: Release{
				Name:      "x-stevedore",
				Namespace: "ns",
				Chart:     "chart/x-stevedore-dependencies",
			},
		}
		ctx := Context{
			Name:              "components",
			KubernetesContext: "gke://components",
			Type:              "components",
			EnvironmentType:   "staging",
			Environment:       "staging",
		}

		expectedPredicate := Predicate{
			conditions: map[string]string{
				"environment":     "staging",
				"environmentType": "staging",
				"contextName":     "components",
				"contextType":     "components",
				"applicationName": "x-stevedore",
			},
		}
		predicate := NewPredicate(app, ctx)

		assert.NotNil(t, predicate)
		assert.Equal(t, expectedPredicate, predicate)
	})
}

func TestNewPredicateFromContext(t *testing.T) {
	t.Run("should create predicate given context", func(t *testing.T) {
		ctx := Context{
			Name:              "components",
			KubernetesContext: "gke://components",
			Type:              "components",
			EnvironmentType:   "staging",
			Environment:       "staging",
		}

		expectedPredicate := Predicate{
			conditions: map[string]string{
				"environment":     "staging",
				"environmentType": "staging",
				"contextName":     "components",
				"contextType":     "components",
			},
		}
		predicate := NewPredicateFromContext(ctx)

		assert.NotNil(t, predicate)
		assert.Equal(t, expectedPredicate, predicate)
	})
}

func TestPredicateContains(t *testing.T) {
	app := ReleaseSpecification{
		Release: Release{
			Name:      "x-stevedore",
			Namespace: "ns",
			Chart:     "chart/x-stevedore-dependencies",
		},
	}

	ctx := Context{
		Name:              "components",
		KubernetesContext: "gke://components",
		Type:              "components",
		EnvironmentType:   "staging",
		Environment:       "staging",
	}

	predicate := NewPredicate(app, ctx)

	t.Run("should return true if given conditions is a subset of the predicate", func(t *testing.T) {
		conditions := map[string]string{
			"contextName": "components",
		}

		assert.True(t, predicate.Contains(conditions))
	})

	t.Run("should return false if given conditions is not a subset of the predicate", func(t *testing.T) {
		conditions := map[string]string{
			"contextName":     "components",
			"contextType":     "components",
			"environmentType": "production",
		}

		assert.False(t, predicate.Contains(conditions))
	})

	t.Run("should return false if given conditions is empty", func(t *testing.T) {
		conditions := map[string]string{}

		assert.False(t, predicate.Contains(conditions))
	})

	t.Run("should return false if given conditions is nil", func(t *testing.T) {
		assert.False(t, predicate.Contains(nil))
	})
}

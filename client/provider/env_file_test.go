package provider_test

import (
	"testing"

	"github.com/gojek/stevedore/client/provider"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/stretchr/testify/assert"
)

func TestEnvsFilesFilter(t *testing.T) {
	t.Run("should return all the applicable envs for the given context", func(t *testing.T) {

		envsFiles := provider.EnvsFiles{
			provider.EnvsFile{
				Name: "x-env",
				EnvSpecifications: stevedore.EnvSpecifications{
					stevedore.EnvSpecification{
						Matches: stevedore.Conditions{
							"environment": "staging",
						},
						Values: stevedore.Substitute{"size": "10Gi"},
					},
				},
			},
			provider.EnvsFile{
				Name: "y-env",
				EnvSpecifications: stevedore.EnvSpecifications{
					stevedore.EnvSpecification{
						Matches: stevedore.Conditions{
							"environment": "production",
						},
						Values: stevedore.Substitute{"persistence": "true"},
					},
				},
			},
		}

		expected := provider.EnvsFiles{
			provider.EnvsFile{
				Name: "x-env",
				EnvSpecifications: stevedore.EnvSpecifications{
					stevedore.EnvSpecification{
						Matches: stevedore.Conditions{
							"environment": "staging",
						},
						Values: stevedore.Substitute{"size": "10Gi"},
					},
				},
			},
		}

		filteredEnvsFiles := envsFiles.Filter(stevedore.Context{Labels: stevedore.Conditions{"environment": "staging"}})

		assert.NotNil(t, filteredEnvsFiles)
		assert.Equal(t, expected, filteredEnvsFiles)
	})

	t.Run("should return empty envs if none of the envs are applicable for the given context", func(t *testing.T) {
		envsFiles := provider.EnvsFiles{
			provider.EnvsFile{
				Name: "x-env",
				EnvSpecifications: stevedore.EnvSpecifications{
					stevedore.EnvSpecification{
						Matches: stevedore.Conditions{
							"environment": "staging",
						},
						Values: stevedore.Substitute{"size": "10Gi"},
					},
				},
			},
			provider.EnvsFile{
				Name: "y-env",
				EnvSpecifications: stevedore.EnvSpecifications{
					stevedore.EnvSpecification{
						Matches: stevedore.Conditions{
							"environment": "production",
						},
						Values: stevedore.Substitute{"persistence": "true"},
					},
				},
			},
		}

		expected := provider.EnvsFiles{}

		filteredEnvsFiles := envsFiles.Filter(stevedore.Context{Labels: stevedore.Conditions{"environment": "uat"}})

		assert.NotNil(t, filteredEnvsFiles)
		assert.Equal(t, expected, filteredEnvsFiles)
	})
}

func TestEnvsFilesSortAndMerge(t *testing.T) {
	t.Run("should merge all the values", func(t *testing.T) {
		envsFiles := provider.EnvsFiles{
			provider.EnvsFile{
				Name: "x-env",
				EnvSpecifications: stevedore.EnvSpecifications{
					stevedore.EnvSpecification{
						Matches: stevedore.Conditions{
							"contextName": "cluster",
						},
						Values: stevedore.Substitute{"name": "x-env", "size": "8Gi"},
					},
				},
			},
			provider.EnvsFile{
				Name: "y-env",
				EnvSpecifications: stevedore.EnvSpecifications{
					stevedore.EnvSpecification{
						Matches: stevedore.Conditions{
							"contextType": "env",
						},
						Values: stevedore.Substitute{"size": "4Gi", "persistence": "true"},
					},
				},
			},
			provider.EnvsFile{
				Name: "z-env",
				EnvSpecifications: stevedore.EnvSpecifications{
					stevedore.EnvSpecification{
						Matches: stevedore.Conditions{
							"contextType": "env",
							"environment": "staging",
						},
						Values: stevedore.Substitute{"size": "6Gi", "readonly": "true", "primary_slot_name": "readonly_cluster"},
					},
				},
			},
		}

		labels := stevedore.Labels{
			{Name: "environmentType"},
			{Name: "environment"},
			{Name: "contextType"},
			{Name: "contextName"},
			{Name: "applicationName"},
		}
		actual, err := envsFiles.SortAndMerge(stevedore.Substitute{"readonly": "false"}, labels)

		assert.NoError(t, err)
		assert.Equal(t, stevedore.Substitute{"name": "x-env", "size": "8Gi", "persistence": "true", "readonly": "false", "primary_slot_name": "readonly_cluster"}, actual)
	})
}

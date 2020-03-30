package stevedore_test

import (
	"fmt"
	"github.com/gojek/stevedore/pkg/helm"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponses_GroupByFile(t *testing.T) {
	t.Run("should group the responses by filename", func(t *testing.T) {
		nginxResponses := stevedore.Responses{
			{
				File:            "nginx.yaml",
				ReleaseName:     "nginx-1",
				UpstallResponse: helm.UpstallResponse{},
				Err:             nil,
			},
			{
				File:            "nginx.yaml",
				ReleaseName:     "nginx-2",
				UpstallResponse: helm.UpstallResponse{},
				Err:             fmt.Errorf("nginx-2 error"),
			},
		}
		postgresResponses := stevedore.Responses{
			{
				File:            "postgres.yaml",
				ReleaseName:     "postgres",
				UpstallResponse: helm.UpstallResponse{},
				Err:             nil,
			},
		}
		responses := stevedore.Responses{nginxResponses[0], nginxResponses[1], postgresResponses[0]}

		expected := stevedore.GroupedResponses{
			"nginx.yaml":    nginxResponses,
			"postgres.yaml": postgresResponses,
		}
		assert.Equal(t, expected, responses.GroupByFile())
	})
}

func TestResponses_SortByManifest(t *testing.T) {
	t.Run("should sort by manifest name", func(t *testing.T) {
		responses := stevedore.Responses{
			{
				File:            "nginx.yaml",
				ReleaseName:     "a-nginx",
				UpstallResponse: helm.UpstallResponse{},
				Err:             nil,
			},
			{
				File:            "nginx.yaml",
				ReleaseName:     "z-nginx",
				UpstallResponse: helm.UpstallResponse{},
				Err:             fmt.Errorf("nginx-2 error"),
			},
			{
				File:            "postgres.yaml",
				ReleaseName:     "postgres",
				UpstallResponse: helm.UpstallResponse{},
				Err:             nil,
			},
		}
		responses.SortByReleaseName()

		expected := stevedore.Responses{
			{
				File:            "nginx.yaml",
				ReleaseName:     "a-nginx",
				UpstallResponse: helm.UpstallResponse{},
				Err:             nil,
			},
			{
				File:            "postgres.yaml",
				ReleaseName:     "postgres",
				UpstallResponse: helm.UpstallResponse{},
				Err:             nil,
			},
			{
				File:            "nginx.yaml",
				ReleaseName:     "z-nginx",
				UpstallResponse: helm.UpstallResponse{},
				Err:             fmt.Errorf("nginx-2 error"),
			},
		}
		assert.Equal(t, expected, responses)
	})
}

func TestGroupedResponses_SortByFileName(t *testing.T) {
	nginxResponse := stevedore.Responses{}
	postgresResponse := stevedore.Responses{}
	redisResponse := stevedore.Responses{}
	expected := []string{
		"nginx.yaml",
		"postgres.yaml",
		"redis.yaml",
	}

	groupedResponses := stevedore.GroupedResponses{
		"redis.yaml":    redisResponse,
		"nginx.yaml":    nginxResponse,
		"postgres.yaml": postgresResponse,
	}

	assert.Equal(t, expected, groupedResponses.SortedFileNames())
}

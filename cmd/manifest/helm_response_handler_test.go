package manifest

import (
	"fmt"
	"testing"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/stretchr/testify/assert"
)

func TestPlan_render(t *testing.T) {
	t.Run("should return nil when no error in response", func(t *testing.T) {
		r := newResponseHandler(true)
		group := stevedore.GroupedResponses{
			"nginx.yaml": stevedore.Responses{
				{
					File:        "nginx.yaml",
					ReleaseName: "nginx-a",
					Err:         nil,
				},
				{
					File:        "nginx.yaml",
					ReleaseName: "nginx-b",
					Err:         nil,
				},
			},
		}
		err := r.render(group)

		assert.NoError(t, err)
	})

	t.Run("should return nil when response has error", func(t *testing.T) {
		r := newResponseHandler(true)
		group := stevedore.GroupedResponses{
			"nginx.yaml": stevedore.Responses{
				{
					File:        "nginx.yaml",
					ReleaseName: "nginx-a",
					Err:         fmt.Errorf("error installing nginx"),
				},
				{
					File:        "nginx.yaml",
					ReleaseName: "nginx-b",
					Err:         nil,
				},
			},
		}
		err := r.render(group)

		assert.NoError(t, err)
	})
}

func TestApply_render(t *testing.T) {
	t.Run("should return nil when no error in response", func(t *testing.T) {
		r := newResponseHandler(false)
		group := stevedore.GroupedResponses{
			"nginx.yaml": stevedore.Responses{
				{
					File:        "nginx.yaml",
					ReleaseName: "nginx-a",
					Err:         nil,
				},
				{
					File:        "nginx.yaml",
					ReleaseName: "nginx-b",
					Err:         nil,
				},
			},
		}
		err := r.render(group)

		assert.NoError(t, err)
	})

	t.Run("should return error when response has error", func(t *testing.T) {
		r := newResponseHandler(false)
		group := stevedore.GroupedResponses{
			"nginx.yaml": stevedore.Responses{
				{
					File:        "nginx.yaml",
					ReleaseName: "nginx-a",
					Err:         fmt.Errorf("error installing nginx"),
				},
				{
					File:        "nginx.yaml",
					ReleaseName: "nginx-b",
					Err:         nil,
				},
			},
		}
		err := r.render(group)

		assert.Error(t, err, "error installing nginx")
	})
}

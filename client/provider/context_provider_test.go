package provider_test

import (
	"github.com/gojek/stevedore/client/provider"
	"testing"

	"github.com/gojek/stevedore/client/internal/mocks"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestDefaultContextProvider(t *testing.T) {
	t.Run("should return new instance of provider", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		contextString := `
current: components
contexts:
  - name: components
    environment: env
    kubernetesContext: components
    environmentType: staging
    type: components
  - name: services
    environment: env
    kubernetesContext: services
    environmentType: production
    type: services`

		contextFile := "/mock/contextFile"
		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		memFs := afero.NewMemMapFs()

		mockEnvironment.EXPECT().Fetch().Return(map[string]interface{}{})

		_ = afero.WriteFile(memFs, contextFile, []byte(contextString), 0644)
		expected := stevedore.Context{
			Name:              "components",
			Type:              "components",
			Environment:       "env",
			KubernetesContext: "components",
			EnvironmentType:   "staging",
		}

		contextProvider := provider.NewContextProvider(memFs, contextFile, mockEnvironment)
		context, err := contextProvider.Context()

		assert.NoError(t, err)
		assert.Equal(t, expected, context)
	})

	t.Run("should return error if context file is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		contextFile := "/mock/contextFile"
		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		memFs := afero.NewMemMapFs()

		contextProvider := provider.NewContextProvider(memFs, contextFile, mockEnvironment)
		context, err := contextProvider.Context()

		if assert.Error(t, err) {
			assert.Equal(t, "current context is not set", err.Error())
		}
		assert.Equal(t, stevedore.Context{}, context)
	})
}

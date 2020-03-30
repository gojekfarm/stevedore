package stevedore_test

import (
	"fmt"
	"github.com/gojek/stevedore/client/provider"
	"github.com/gojek/stevedore/pkg/internal/mocks/plugin"
	pkgPlugin "github.com/gojek/stevedore/pkg/plugin"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigsFetch(t *testing.T) {
	t.Run("should return values fetched from store", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var pluginAConfigs []map[string]interface{}
		var pluginBConfigs []map[string]interface{}

		pluginAResponse := map[string]interface{}{
			"name": "x-service",
			"type": "server",
		}
		pluginBResponse := map[string]interface{}{
			"app-name": "y-service",
			"app-env":  "staging",
		}

		stevedoreContext := stevedore.Context{Environment: "staging"}

		pluginAConfigProvider := mockPlugin.NewMockConfigInterface(ctrl)
		pluginAConfigProvider.EXPECT().Type().Return(pkgPlugin.TypeConfig, nil)
		pluginAConfigProvider.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(pluginAResponse, nil)
		pluginBConfigProvider := mockPlugin.NewMockConfigInterface(ctrl)
		pluginBConfigProvider.EXPECT().Type().Return(pkgPlugin.TypeConfig, nil)
		pluginBConfigProvider.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(pluginBResponse, nil)

		plugins := provider.Plugins{"pluginAStore": provider.ClientPlugin{PluginImpl: pluginAConfigProvider}, "pluginBStore": provider.ClientPlugin{PluginImpl: pluginBConfigProvider}}
		configProviders, _ := plugins.ConfigProviders()

		expected := stevedore.Substitute{
			"app-env":  "staging",
			"app-name": "y-service",
			"name":     "x-service",
			"type":     "server",
		}
		configs := stevedore.Configs{
			"pluginAStore": pluginAConfigs,
			"pluginBStore": pluginBConfigs,
		}

		substitutes, err := configs.Fetch(configProviders, stevedoreContext)

		assert.NoError(t, err)
		if !cmp.Equal(expected, substitutes) {
			assert.Fail(t, cmp.Diff(expected, substitutes))
		}
	})

	t.Run("should fail to return values from store when any plugin call fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		err := fmt.Errorf("some error")

		pluginConfigs := []map[string]interface{}{{"name": "plugin"}}

		stevedoreContext := stevedore.Context{Environment: "staging"}
		contextAsMap, _ := stevedoreContext.Map()

		pluginConfigProvider := mockPlugin.NewMockConfigInterface(ctrl)
		pluginConfigProvider.EXPECT().Type().Return(pkgPlugin.TypeConfig, nil)
		pluginConfigProvider.EXPECT().Fetch(contextAsMap, gomock.Any()).Return(nil, err)

		plugins := provider.Plugins{"pluginStore": provider.ClientPlugin{PluginImpl: pluginConfigProvider}}
		configProviders, _ := plugins.ConfigProviders()

		configs := stevedore.Configs{
			"pluginStore": pluginConfigs,
		}

		substitutes, err := configs.Fetch(configProviders, stevedoreContext)

		assert.Error(t, err)
		assert.Equal(t, "error in fetching from provider: some error", err.Error())
		assert.Nil(t, substitutes)
	})
}

package cmd

import (
	"os"
	"testing"

	"github.com/gojek/stevedore/client/provider"
	"github.com/gojek/stevedore/cmd/internal/mocks/mockPlugin"
	pkgPlugin "github.com/gojek/stevedore/pkg/plugin"
	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestGetPlugins(t *testing.T) {
	t.Run("should return the list of plugins", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expected := []PluginInfo{
			{name: "pluginA", pluginType: pkgPlugin.TypeConfig, version: "1.0.0"},
			{name: "pluginB", pluginType: pkgPlugin.TypeConfig, version: "2.0.0"},
		}

		pluginA := mockPlugin.NewMockConfigInterface(ctrl)
		pluginB := mockPlugin.NewMockConfigInterface(ctrl)

		pluginA.EXPECT().Version().Return("1.0.0", nil)
		pluginA.EXPECT().Type().Return(pkgPlugin.TypeConfig, nil)
		pluginB.EXPECT().Version().Return("2.0.0", nil)
		pluginB.EXPECT().Type().Return(pkgPlugin.TypeConfig, nil)

		plugins := provider.Plugins{"pluginA": provider.ClientPlugin{PluginImpl: pluginA}, "pluginB": provider.ClientPlugin{PluginImpl: pluginB}}

		pluginInfos, err := getPluginsInfo(plugins)

		assert.NoError(t, err)
		assert.ElementsMatch(t, expected, pluginInfos)
	})
}

func TestDownloadFile(t *testing.T) {
	t.Run("should download URL", func(t *testing.T) {
		defer gock.Off()

		bodyString := "this is a downloaded file"
		gock.New("http://some-url.com").
			Get("/plugin").
			Reply(200).
			BodyString(bodyString)

		actual, err := downloadFile("http://some-url.com/plugin")
		assert.NoError(t, err)
		assert.Equal(t, []byte(bodyString), actual)
	})
}

func TestFindSource(t *testing.T) {
	t.Run("should find remote URL", func(t *testing.T) {
		src, err := findSource("http://some-url")

		assert.NoError(t, err)
		assert.Equal(t, remote, src)
	})

	t.Run("should find local file", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		localFilename := "some-local-file"
		err := afero.WriteFile(fs, localFilename, []byte("content"), 0644)
		assert.NoError(t, err)

		src, err := findSource(localFilename)

		assert.NoError(t, err)
		assert.Equal(t, local, src)
	})
}

func TestMain(m *testing.M) {
	ret := m.Run()
	closePlugins()
	os.Exit(ret)
}

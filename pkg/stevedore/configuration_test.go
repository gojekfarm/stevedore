package stevedore

import (
	"fmt"
	"testing"

	"github.com/gojek/stevedore/pkg/internal/mocks"
	"github.com/golang/mock/gomock"
	"gopkg.in/yaml.v2"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const ConfigFileName = "/mock/.stevedore/config"

func saveConfig(config string) afero.Fs {
	memFs := afero.NewMemMapFs()
	_ = afero.WriteFile(memFs, ConfigFileName, []byte(config), 0644)
	return memFs
}

func readConfigurationFromFs(t *testing.T, fs afero.Fs, filename string) Configuration {
	var savedConfiguration Configuration

	data, err := afero.ReadFile(fs, filename)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	err = yaml.Unmarshal(data, &savedConfiguration)
	assert.NoError(t, err)

	return savedConfiguration
}

func TestOverriddenContext(t *testing.T) {
	t.Run("should return true and the overridden context name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		mockEnvironment.EXPECT().Fetch().Return(map[string]interface{}{"STEVEDORE_CONTEXT": "dev"})

		current, ok := OverriddenContext(mockEnvironment)

		assert.True(t, ok)
		assert.Equal(t, "dev", current)
	})

	t.Run("should return false and the empty context name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		mockEnvironment.EXPECT().Fetch().Return(map[string]interface{}{})

		current, ok := OverriddenContext(mockEnvironment)

		assert.False(t, ok)
		assert.Equal(t, "", current)
	})
}

func TestNewConfigurationFromFile(t *testing.T) {
	t.Run("should successfully load", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		mockEnvironment.EXPECT().Fetch().Return(map[string]interface{}{})
		configString := `
current: services
contexts:
  - name: components
    kubernetesContext: components
    labels:
      environment: env
      environmentType: staging
      type: components
  - name: services
    kubernetesContext: services
    labels:
      environment: env
      environmentType: production
      type: services`

		memFs := saveConfig(configString)

		expectedConfigurationConfig := Configuration{
			Current: "services",
			Contexts: []Context{
				{
					Name:              "components",
					KubernetesContext: "components",
					Labels: Conditions{
						"environment":     "env",
						"type":            "components",
						"environmentType": "staging",
					},
				},
				{
					Name:              "services",
					KubernetesContext: "services",
					Labels: Conditions{
						"environment":     "env",
						"type":            "services",
						"environmentType": "production",
					},
				},
			},
			filename: ConfigFileName,
			fs:       memFs,
		}

		actualConfigurationConfig, err := NewConfigurationFromFile(memFs, ConfigFileName, mockEnvironment)

		assert.NoError(t, err)
		if assert.NotNil(t, actualConfigurationConfig) {
			assert.Equal(t, expectedConfigurationConfig, *actualConfigurationConfig)
		}
	})

	t.Run("should use current stevedore context from env variable", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEnvironment := mocks.NewMockEnvironment(ctrl)
		mockEnvironment.EXPECT().Fetch().Return(map[string]interface{}{"STEVEDORE_CONTEXT": "components"})

		configString := `
current: services
contexts:
  - name: components
    kubernetesContext: components
    labels:
      environment: env
  - name: services
    kubernetesContext: services
    labels:
      environment: env`

		memFs := saveConfig(configString)

		expectedConfigurationConfig := Configuration{
			Current: "components",
			Contexts: []Context{
				{
					Name:              "components",
					KubernetesContext: "components",
					Labels: Conditions{
						"environment": "env",
					},
				},
				{
					Name:              "services",
					KubernetesContext: "services",
					Labels: Conditions{
						"environment": "env",
					},
				},
			},
			filename: ConfigFileName,
			fs:       memFs,
		}

		actualConfigurationConfig, err := NewConfigurationFromFile(memFs, ConfigFileName, mockEnvironment)

		assert.NoError(t, err)
		if assert.NotNil(t, actualConfigurationConfig) {
			assert.Equal(t, expectedConfigurationConfig, *actualConfigurationConfig)
		}
	})

	t.Run("should throw error if read file fails at fs layer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFs := mocks.NewMockFs(ctrl)
		mockFs.EXPECT().Stat(ConfigFileName)
		mockFs.EXPECT().Open(ConfigFileName).Return(nil, fmt.Errorf("error when reading file"))

		mockEnvironment := mocks.NewMockEnvironment(ctrl)

		actualConfigurationConfig, err := NewConfigurationFromFile(mockFs, ConfigFileName, mockEnvironment)
		assert.Error(t, err)

		assert.Nil(t, actualConfigurationConfig)
	})

	t.Run("should fail if not valid yaml", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEnvironment := mocks.NewMockEnvironment(ctrl)

		configString := `
store
  host: host-url`

		memFs := saveConfig(configString)

		actualConfigurationConfig, err := NewConfigurationFromFile(memFs, ConfigFileName, mockEnvironment)
		assert.Error(t, err)

		assert.Nil(t, actualConfigurationConfig)
	})
}

func TestConfigurationUse(t *testing.T) {
	t.Run("should set current context", func(t *testing.T) {
		stevedore := &Configuration{
			Current: "components",
			Contexts: Contexts{
				Context{Name: "services"},
				Context{Name: "components"},
			},
			filename: ConfigFileName,
		}
		fileContent, _ := yaml.Marshal(stevedore)
		memFs := saveConfig(string(fileContent))
		stevedore.fs = memFs

		err := stevedore.Use("services")

		assert.NoError(t, err)

		data, err := afero.ReadFile(memFs, ConfigFileName)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		var savedConfiguration Configuration

		err = yaml.Unmarshal(data, &savedConfiguration)
		assert.NoError(t, err)

		assert.Equal(t, "services", savedConfiguration.Current)
	})

	t.Run("should fail if context is not valid", func(t *testing.T) {
		stevedore := &Configuration{
			Current: "components",
			Contexts: Contexts{
				Context{Name: "dev-007"},
				Context{Name: "components"},
			},
			filename: ConfigFileName,
		}
		fileContent, _ := yaml.Marshal(stevedore)
		memFs := saveConfig(string(fileContent))
		stevedore.fs = memFs

		err := stevedore.Use("not-a-valid-context")

		assert.Error(t, err)
		assert.Equal(t, "components", stevedore.Current)

		data, err := afero.ReadFile(memFs, ConfigFileName)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		var currentConfiguration Configuration

		err = yaml.Unmarshal(data, &currentConfiguration)
		assert.NoError(t, err)

		assert.Equal(t, "components", currentConfiguration.Current)
	})
}

func TestConfigurationAdd(t *testing.T) {
	t.Run("should add context to stevedore config", func(t *testing.T) {
		stevedore := &Configuration{
			Current: "components",
			Contexts: Contexts{
				Context{Name: "components"},
			},
			filename: ConfigFileName,
		}
		fileContent, _ := yaml.Marshal(stevedore)
		memFs := saveConfig(string(fileContent))
		stevedore.fs = memFs

		ctx := Context{
			Name:              "services",
			KubernetesContext: "services",
			Labels: Conditions{
				"environment":     "env",
				"type":            "services",
				"environmentType": "staging",
			},
		}

		err := stevedore.Add(ctx)

		assert.NoError(t, err)

		data, err := afero.ReadFile(memFs, ConfigFileName)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		var savedConfiguration Configuration

		err = yaml.Unmarshal(data, &savedConfiguration)
		assert.NoError(t, err)

		expectedContexts := Contexts{
			Context{Name: "components"},
			Context{
				Name:              "services",
				KubernetesContext: "services",
				Labels: Conditions{
					"environment":     "env",
					"type":            "services",
					"environmentType": "staging",
				},
			},
		}

		assert.Equal(t, expectedContexts, savedConfiguration.Contexts)
	})

	t.Run("should fail if context already exists", func(t *testing.T) {
		stevedore := &Configuration{
			Current: "components",
			Contexts: Contexts{
				Context{Name: "components"},
			},
			filename: ConfigFileName,
		}
		fileContent, _ := yaml.Marshal(stevedore)
		memFs := saveConfig(string(fileContent))
		stevedore.fs = memFs

		ctx := Context{
			Name:              "components",
			KubernetesContext: "components",
			Labels: Conditions{
				"environment": "env",
			},
		}

		err := stevedore.Add(ctx)

		assert.Error(t, err)

		data, err := afero.ReadFile(memFs, ConfigFileName)
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		var savedConfiguration Configuration

		err = yaml.Unmarshal(data, &savedConfiguration)
		assert.NoError(t, err)

		expectedContexts := Contexts{
			Context{Name: "components"},
		}

		assert.Equal(t, expectedContexts, savedConfiguration.Contexts)
	})
}

func TestConfigurationDelete(t *testing.T) {
	t.Run("should delete context from stevedore config", func(t *testing.T) {
		stevedore := &Configuration{
			Current: "components",
			Contexts: Contexts{
				Context{Name: "components"},
				Context{Name: "services"},
				Context{Name: "env-services"},
			},
			filename: ConfigFileName,
		}
		fileContent, _ := yaml.Marshal(stevedore)
		memFs := saveConfig(string(fileContent))
		stevedore.fs = memFs

		err := stevedore.Delete("services")
		assert.NoError(t, err)

		savedConfiguration := readConfigurationFromFs(t, memFs, ConfigFileName)
		expectedContexts := Contexts{
			Context{Name: "components"},
			Context{Name: "env-services"},
		}

		assert.Equal(t, expectedContexts, savedConfiguration.Contexts)
		assert.Equal(t, "components", savedConfiguration.Current)
	})

	t.Run("should reset current context if it is getting deleted", func(t *testing.T) {
		stevedore := &Configuration{
			Current: "services",
			Contexts: Contexts{
				Context{Name: "components"},
				Context{Name: "services"},
			},
			filename: ConfigFileName,
		}
		fileContent, _ := yaml.Marshal(stevedore)
		memFs := saveConfig(string(fileContent))
		stevedore.fs = memFs

		err := stevedore.Delete("services")
		assert.NoError(t, err)

		savedConfiguration := readConfigurationFromFs(t, memFs, ConfigFileName)
		assert.Equal(t, "", savedConfiguration.Current)
	})

	t.Run("should error out if context is not valid", func(t *testing.T) {
		stevedore := &Configuration{
			Current: "services",
			Contexts: Contexts{
				Context{Name: "components"},
				Context{Name: "services"},
			},
			filename: ConfigFileName,
		}
		fileContent, _ := yaml.Marshal(stevedore)
		memFs := saveConfig(string(fileContent))
		stevedore.fs = memFs

		err := stevedore.Delete("Some-random-context")
		assert.Error(t, err)

		savedConfiguration := readConfigurationFromFs(t, memFs, ConfigFileName)
		expectedContexts := Contexts{
			Context{Name: "components"},
			Context{Name: "services"},
		}

		assert.Equal(t, expectedContexts, savedConfiguration.Contexts)
	})
}

func TestConfigurationRename(t *testing.T) {
	t.Run("should rename context from stevedore config", func(t *testing.T) {
		stevedore := &Configuration{
			Current: "components",
			Contexts: Contexts{
				Context{Name: "components"},
				Context{Name: "services"},
			},
			filename: ConfigFileName,
		}
		fileContent, _ := yaml.Marshal(stevedore)
		memFs := saveConfig(string(fileContent))
		stevedore.fs = memFs

		err := stevedore.Rename("services", "dev-services")
		assert.NoError(t, err)

		savedConfiguration := readConfigurationFromFs(t, memFs, ConfigFileName)
		expectedContexts := Contexts{
			Context{Name: "components"},
			Context{Name: "dev-services"},
		}

		assert.Equal(t, expectedContexts, savedConfiguration.Contexts)
		assert.Equal(t, "components", savedConfiguration.Current)
	})

	t.Run("should fail if source context does not exists", func(t *testing.T) {
		stevedore := &Configuration{
			Current: "components",
			Contexts: Contexts{
				Context{Name: "components"},
				Context{Name: "services"},
			},
			filename: ConfigFileName,
		}
		fileContent, _ := yaml.Marshal(stevedore)
		memFs := saveConfig(string(fileContent))
		stevedore.fs = memFs

		err := stevedore.Rename("some-random-context", "dev-services")
		assert.Error(t, err)

		savedConfiguration := readConfigurationFromFs(t, memFs, ConfigFileName)
		expectedContexts := Contexts{
			Context{Name: "components"},
			Context{Name: "services"},
		}

		assert.Equal(t, expectedContexts, savedConfiguration.Contexts)
	})

	t.Run("should fail if target context already exists", func(t *testing.T) {
		stevedore := &Configuration{
			Current: "components",
			Contexts: Contexts{
				Context{Name: "components"},
				Context{Name: "services"},
			},
			filename: ConfigFileName,
		}
		fileContent, _ := yaml.Marshal(stevedore)
		memFs := saveConfig(string(fileContent))
		stevedore.fs = memFs

		err := stevedore.Rename("services", "components")
		assert.Error(t, err)

		savedConfiguration := readConfigurationFromFs(t, memFs, ConfigFileName)
		expectedContexts := Contexts{
			Context{Name: "components"},
			Context{Name: "services"},
		}

		assert.Equal(t, expectedContexts, savedConfiguration.Contexts)
	})

	t.Run("should change current context if renamed", func(t *testing.T) {
		stevedore := &Configuration{
			Current: "components",
			Contexts: Contexts{
				Context{Name: "components"},
				Context{Name: "services"},
			},
			filename: ConfigFileName,
		}
		fileContent, _ := yaml.Marshal(stevedore)
		memFs := saveConfig(string(fileContent))
		stevedore.fs = memFs

		err := stevedore.Rename("components", "dev-components")
		assert.NoError(t, err)

		savedConfiguration := readConfigurationFromFs(t, memFs, ConfigFileName)
		expectedContexts := Contexts{
			Context{Name: "dev-components"},
			Context{Name: "services"},
		}

		assert.Equal(t, expectedContexts, savedConfiguration.Contexts)
		assert.Equal(t, "dev-components", savedConfiguration.Current)
	})
}

func TestConfigurationCurrentContext(t *testing.T) {
	t.Run("should return error if current context is not set", func(t *testing.T) {
		configuration := Configuration{}

		ctx, err := configuration.CurrentContext()

		assert.NotNil(t, err)
		assert.Equal(t, "current context is not set", err.Error())
		assert.Equal(t, Context{}, ctx)
	})

	t.Run("should return error if current context is not valid", func(t *testing.T) {
		configuration := Configuration{Current: "services", Contexts: Contexts{{Name: "components"}}}

		ctx, err := configuration.CurrentContext()

		assert.NotNil(t, err)
		assert.Equal(t, "unable to find current context services", err.Error())
		assert.Equal(t, Context{}, ctx)
	})

	t.Run("should return current context", func(t *testing.T) {
		configuration := Configuration{Current: "services", Contexts: Contexts{{Name: "services"}}}

		ctx, err := configuration.CurrentContext()

		assert.Nil(t, err)
		assert.Equal(t, Context{Name: "services"}, ctx)
	})
}

func TestConfiguration_Labels(t *testing.T) {
	t.Run("should add context, environment, applicationName", func(t *testing.T) {
		configuration := Configuration{
			UserLabels: Labels{
				{Name: "one"},
				{Name: "two"},
				{Name: "three"},
			},
		}

		expected := Labels{
			{Name: "one"},
			{Name: "two"},
			{Name: "three"},
			{Name: ConditionEnvironment},
			{Name: ConditionContextName},
			{Name: ConditionApplicationName},
		}

		actual := configuration.Labels()

		assert.Equal(t, expected, actual)
	})

	t.Run("should add context, environment, applicationName and assign weights accordingly", func(t *testing.T) {
		configuration := Configuration{
			UserLabels: Labels{
				{Name: "one", Weight: 2},
				{Name: "two", Weight: 3},
				{Name: "three", Weight: 1},
			},
		}

		expected := Labels{
			{Name: "one", Weight: 2},
			{Name: "two", Weight: 3},
			{Name: "three", Weight: 1},
			{Name: ConditionEnvironment, Weight: 4},
			{Name: ConditionContextName, Weight: 5},
			{Name: ConditionApplicationName, Weight: 6},
		}

		actual := configuration.Labels()

		assert.Equal(t, expected, actual)
	})
}

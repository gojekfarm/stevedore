package kubeconfig_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gojek/stevedore/cmd/kubeconfig"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestMockHomePathResolver(t *testing.T) {
	t.Run("should return path to default kubeconfig file", func(t *testing.T) {
		homeDir := func() (string, error) { return "/Users/someone", nil }

		f, err := kubeconfig.DefaultFile(homeDir)

		assert.NoError(t, err)
		assert.Equal(t, f, "/Users/someone/.kube/config")
	})

	t.Run("should fail if home dir cannot be determined", func(t *testing.T) {
		homeDir := func() (string, error) { return "", fmt.Errorf("unable to find home dir") }

		f, err := kubeconfig.DefaultFile(homeDir)

		assert.Equal(t, f, "")
		assert.EqualError(t, err, "failed to get default kubeconfig file path: unable to find home dir")
	})
}

var userHomePath = "/Users/someone"
var mockHomePathResolver = func() (string, error) {
	return userHomePath, nil
}

func TestFindInConfigs(t *testing.T) {
	context := stevedore.Context{
		Name:              "context-a",
		Type:              "components",
		Environment:       "staging",
		KubernetesContext: "minikube",
		EnvironmentType:   "staging",
	}

	t.Run("it should return the config file without extension", func(t *testing.T) {
		t.Run("yaml", func(t *testing.T) {
			memFs := afero.NewMemMapFs()
			configPath := filepath.Join(userHomePath, ".kube", "configs", context.KubernetesContext)
			err := afero.WriteFile(memFs, configPath, []byte(""), 0644)
			assert.NoError(t, err)

			actual, err := kubeconfig.FindInConfigs(mockHomePathResolver, memFs, context)

			assert.NoError(t, err)
			assert.Equal(t, "/Users/someone/.kube/configs/minikube", actual)
		})
	})

	t.Run("it should return the config file with extension", func(t *testing.T) {

		t.Run("yaml", func(t *testing.T) {
			memFs := afero.NewMemMapFs()
			configPath := filepath.Join(userHomePath, ".kube", "configs", fmt.Sprintf("%s.yaml", context.KubernetesContext))
			err := afero.WriteFile(memFs, configPath, []byte(""), 0644)
			assert.NoError(t, err)

			actual, err := kubeconfig.FindInConfigs(mockHomePathResolver, memFs, context)

			assert.NoError(t, err)
			assert.Equal(t, "/Users/someone/.kube/configs/minikube.yaml", actual)
		})

		t.Run("yml", func(t *testing.T) {
			memFs := afero.NewMemMapFs()
			configPath := filepath.Join(userHomePath, ".kube", "configs", fmt.Sprintf("%s.yml", context.KubernetesContext))
			err := afero.WriteFile(memFs, configPath, []byte(""), 0644)
			assert.NoError(t, err)

			actual, err := kubeconfig.FindInConfigs(mockHomePathResolver, memFs, context)

			assert.NoError(t, err)
			assert.Equal(t, "/Users/someone/.kube/configs/minikube.yml", actual)
		})
	})
}

func TestResolveAndValidate(t *testing.T) {
	context := stevedore.Context{
		Name:              "context-a",
		Type:              "components",
		Environment:       "staging",
		KubernetesContext: "minikube",
		EnvironmentType:   "staging",
	}

	t.Run("should return the given kubeconfig if its not empty", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"
		_ = afero.WriteFile(memFs, filename, []byte("file content"), 0644)

		k, err := kubeconfig.ResolveAndValidate(mockHomePathResolver, "/mock/file", memFs, context)

		assert.NoError(t, err)
		assert.Equal(t, k, "/mock/file")
	})

	t.Run("should get value from context if given kubeconfig is empty and context has kubeconfig", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		contextWithKubeConfig := stevedore.Context{
			Name:              "context-a",
			Type:              "components",
			Environment:       "staging",
			KubernetesContext: "minikube",
			EnvironmentType:   "staging",
			KubeConfigFile:    "~/.kube/configs/minikube",
		}

		_ = afero.WriteFile(memFs, contextWithKubeConfig.KubeConfigFile, []byte("file content"), 0644)

		filename := "/mock/file"
		_ = afero.WriteFile(memFs, filename, []byte("file content"), 0644)

		_ = os.Setenv("KUBECONFIG", "/mock/file")
		defer func() { _ = os.Unsetenv("KUBECONFIG") }()

		actual, err := kubeconfig.ResolveAndValidate(mockHomePathResolver, "", memFs, contextWithKubeConfig)

		assert.NoError(t, err)
		assert.Equal(t, contextWithKubeConfig.KubeConfigFile, actual)
	})

	t.Run("should get value from env if given kubeconfig and kubeconfig in context is empty", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"
		_ = afero.WriteFile(memFs, filename, []byte("file content"), 0644)

		_ = os.Setenv("KUBECONFIG", "/mock/file")
		defer func() { _ = os.Unsetenv("KUBECONFIG") }()

		k, err := kubeconfig.ResolveAndValidate(mockHomePathResolver, "", memFs, context)

		assert.NoError(t, err)
		assert.Equal(t, k, "/mock/file")
	})

	t.Run("should get file with context name in $HOME/.kube/config/", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		userHomePath := "/Users/someone"
		mockHomePathResolver := func() (string, error) {
			return userHomePath, nil
		}

		configPath := filepath.Join(userHomePath, ".kube", "configs", context.KubernetesContext)
		err := afero.WriteFile(memFs, configPath, []byte(""), 0644)
		assert.NoError(t, err)

		actual, err := kubeconfig.ResolveAndValidate(mockHomePathResolver, "", memFs, context)
		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(userHomePath, ".kube", "configs", context.KubernetesContext), actual)
	})

	t.Run("should return default value if given kubeconfig and kubeconfig in context is empty and KUBECONFIG env is empty", func(t *testing.T) {
		homeDir := func() (string, error) { return "/Users/someone", nil }
		memFs := afero.NewMemMapFs()
		filename := "/Users/someone/.kube/config"
		_ = afero.WriteFile(memFs, filename, []byte("file content"), 0644)

		actual, err := kubeconfig.ResolveAndValidate(homeDir, "", memFs, context)

		assert.NoError(t, err)
		assert.Equal(t, "/Users/someone/.kube/config", actual)
	})

	t.Run("should return error if finding default value fails", func(t *testing.T) {
		homeDir := func() (string, error) { return "", fmt.Errorf("unable to find home dir") }
		memFs := afero.NewMemMapFs()

		k, err := kubeconfig.ResolveAndValidate(homeDir, "", memFs, context)

		assert.Empty(t, k)
		assert.EqualError(t, err, "failed to get default kubeconfig file path: unable to find home dir")
	})

	t.Run("should return error if kubeconfig file does not exists", func(t *testing.T) {
		memFs := afero.NewMemMapFs()

		k, err := kubeconfig.ResolveAndValidate(mockHomePathResolver, "/mock/file", memFs, context)

		assert.Empty(t, k)
		assert.EqualError(t, err, "unable to find kubeconfig file /mock/file due to open /mock/file: file does not exist")
	})
}

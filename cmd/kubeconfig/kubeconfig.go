package kubeconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
)

const kubeconfigEnv = "KUBECONFIG"

var homeDir = func() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("unable to find home directory: %v", err)
	}

	return home, nil
}

// HomeDirResolver resolves user home director on invoke
type HomeDirResolver = func() (string, error)

// OSHomeDirResolver resolves user home director based on the
// operating system
var OSHomeDirResolver = homeDir

// DefaultFile returns the path to default kubeconfig file
func DefaultFile(homeDirResolver HomeDirResolver) (string, error) {
	home, err := homeDirResolver()
	if err != nil {
		return "", fmt.Errorf("failed to get default kubeconfig file path: %v", err)
	}
	return filepath.Join(home, ".kube", "config"), nil
}

// FindInConfigs returns the path if find the default kube configs folder
func FindInConfigs(homeDirResolver HomeDirResolver, fs afero.Fs, context stevedore.Context) (string, error) {
	home, err := homeDirResolver()
	baseErrFmt := "failed to get default kubeconfig file path: %v"
	if err != nil {
		return "", fmt.Errorf(baseErrFmt, err)
	}

	contextPath := filepath.Join(home, ".kube", "configs", fmt.Sprintf("%s*", context.KubernetesContext))
	files, err := afero.Glob(fs, contextPath)

	if err != nil {
		return "", fmt.Errorf(baseErrFmt, err)
	}

	if len(files) == 0 {
		return "", nil
	}

	return files[0], nil
}

// ResolveAndValidate returns the final value of kubeconfig after
// considering the given kubeconfig value, env and
// default value respectively. Finally it checks if the
// resolved kubeconfig file exists in the filesystem
func ResolveAndValidate(homeDirResolver HomeDirResolver, kubeconfig string, fs afero.Fs, context stevedore.Context) (resolvedKubeconfig string, err error) {
	if strings.TrimSpace(kubeconfig) != "" {
		resolvedKubeconfig = kubeconfig
	} else if context.KubeConfigFile != "" {
		resolvedKubeconfig = context.KubeConfigFile
	} else if os.Getenv(kubeconfigEnv) != "" {
		resolvedKubeconfig = os.Getenv(kubeconfigEnv)
	} else {
		if resolvedKubeconfig, err = FindInConfigs(homeDirResolver, fs, context); err == nil && resolvedKubeconfig == "" {
			resolvedKubeconfig, err = DefaultFile(homeDirResolver)
		}
	}

	if _, err := fs.Stat(resolvedKubeconfig); err != nil {
		return "", fmt.Errorf("unable to find kubeconfig file %s due to %v", resolvedKubeconfig, err)
	}
	return resolvedKubeconfig, err
}

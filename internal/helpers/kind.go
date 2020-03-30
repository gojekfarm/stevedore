package helpers

import (
	"fmt"
	"github.com/gojek/stevedore/cmd/cli"
	"github.com/gojek/stevedore/internal/cli/kind"
	stringUtils "github.com/gojek/stevedore/pkg/utils/string"
	"github.com/spf13/afero"
	"path"
)

func kubeConfigFile(clusterName string) string {
	return path.Join(testDir(), fmt.Sprintf("%s.yaml", clusterName))
}

func isClusterAlreadyExists(clusterName string) (bool, error) {
	command, err := kind.GetClusters()
	if err != nil {
		return false, err
	}

	out, err := execute(command)
	if err != nil {
		return false, err
	}

	clusters := splitByNewLine(out)
	return stringUtils.Contains(clusters, clusterName), nil
}

func saveExistingKubeConfig(clusterName string) error {
	kubeConfigFile := kubeConfigFile(clusterName)
	command, err := kind.GetKubeConfig(clusterName)
	if err != nil {
		return err
	}

	kubeConfig, err := execute(command)
	if err != nil {
		return err
	}
	cli.Info(fmt.Sprintf("saving existing cluster info at %s", kubeConfigFile))
	return afero.WriteFile(afero.NewOsFs(), kubeConfigFile, []byte(kubeConfig), 0644)
}

// CreateCluster creates a K8s cluster with the clusterName and version
// optionally recreate can be specified
func CreateCluster(clusterName, version string, recreate bool) error {
	exists, err := isClusterAlreadyExists(clusterName)
	if err != nil {
		return err
	}

	if exists && !recreate {
		err = saveExistingKubeConfig(clusterName)
		if err != nil {
			return err
		}
	}
	if !exists || recreate {
		if exists {
			command, err := kind.Delete(clusterName)
			if err != nil {
				return err
			}

			_, err = execute(command)
			if err != nil {
				return err
			}
		}

		image := fmt.Sprintf("kindest/node:v%s", version)
		kubeConfig := kubeConfigFile(clusterName)

		command, err := kind.Create(image, kubeConfig, clusterName, 0)
		if err != nil {
			return err
		}

		_, err = execute(command)
		return err
	}
	return nil
}

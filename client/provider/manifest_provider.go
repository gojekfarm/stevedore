package provider

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/gojek/stevedore/client/yaml"
	"github.com/gojek/stevedore/pkg/plugin"
	stringUtils "github.com/gojek/stevedore/pkg/utils/string"

	"github.com/gojek/stevedore/pkg/file"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
)

// ManifestRecursiveFlag is a flag for recursive check
const ManifestRecursiveFlag = "recursive"

// ManifestPathFile is a flag for path
const ManifestPathFile = "path"

// EnvironmentTypeKey is a flag for environmentType
const EnvironmentTypeKey = "environmentType"

func init() {
	defaultPlugins["manifests"] = ClientPlugin{PluginImpl: DefaultManifestProvider{fs: afero.NewOsFs()}}
}

var _ plugin.ManifestInterface = DefaultManifestProvider{}

// DefaultManifestProvider represents the default manifest provider
// which reads the manifests from given dir
type DefaultManifestProvider struct {
	fs afero.Fs
}

// Version returns Version of the DefaultManifestProvider
func (provider DefaultManifestProvider) Version() (string, error) {
	return "v0.0.1", nil
}

// Flags returns Flags for the DefaultManifestProvider
func (provider DefaultManifestProvider) Flags() ([]plugin.Flag, error) {
	return []plugin.Flag{
		{Name: ManifestPathFile, Shorthand: "f", Required: true, Usage: "Stevedore manifest(s) path (can be yaml file or folder)"},
		{Name: ManifestRecursiveFlag, Default: "false", Usage: "Respect 'depends_on' in manifest file and load other manifest"},
	}, nil
}

// Type returns Type of the DefaultManifestProvider
func (provider DefaultManifestProvider) Type() (plugin.Type, error) {
	return plugin.TypeManifest, nil
}

// Help returns Help of the DefaultManifestProvider
func (provider DefaultManifestProvider) Help() (string, error) {
	return "This is the provider which reads the files/dir and create manifests for the stevedore", nil
}

// Close the DefaultManifestProvider
func (provider DefaultManifestProvider) Close() error {
	return nil
}

// DevEnvironments contains the environments which can be used to run locally
var DevEnvironments = []string{"dev", "local"}

type mappedReleases map[string]stevedore.ReleaseSpecification

// Manifests return all the manifest
func (provider DefaultManifestProvider) Manifests(context map[string]string) (stevedore.ManifestFiles, error) {
	manifestsPath, recursive, environmentType, err := getValues(provider.fs, context)
	if err != nil {
		return nil, err
	}
	if recursive && !stringUtils.Contains(DevEnvironments, environmentType) {
		return nil, fmt.Errorf("--recursive mode cannot be used in non dev environment")
	}
	isDir, err := afero.IsDir(provider.fs, manifestsPath)
	if err != nil {
		return nil, err
	}

	if isDir || !recursive {
		return getAllManifestFiles(provider.fs, manifestsPath)
	}

	manifestFile, err := getManifestFile(provider.fs, manifestsPath)
	if err != nil {
		return nil, err
	}

	releaseSpecs, err := getAllReleaseSpecifications(provider.fs, path.Dir(manifestsPath))
	if err != nil {
		return nil, err
	}

	result := getDependentReleaseSpecifications(manifestFile.Spec, releaseSpecs, mappedReleases{})
	dependentReleaseSpecifications := stevedore.ReleaseSpecifications{}
	for _, dependentReleaseSpecification := range result {
		dependentReleaseSpecifications = append(dependentReleaseSpecifications, dependentReleaseSpecification)
	}

	dependentManifest := stevedore.Manifest{DeployTo: manifestFile.DeployTo, Spec: dependentReleaseSpecifications}
	return stevedore.ManifestFiles{stevedore.ManifestFile{File: manifestFile.File, Manifest: dependentManifest}}, nil
}

func getValues(fs afero.Fs, context map[string]string) (string, bool, string, error) {
	manifestsPath := context[ManifestPathFile]
	recursive, err := strconv.ParseBool(context[ManifestRecursiveFlag])
	if err != nil {
		return "", false, "", err
	}
	if manifestsPath == "" {
		return "", false, "", fmt.Errorf("provide a valid path to stevedore manifests using --manifests-path")
	}
	if _, err := fs.Stat(manifestsPath); os.IsNotExist(err) {
		return "", false, "", fmt.Errorf("invalid file path. Provide a valid path to stevedore manifests using --manifests-path")
	}

	environmentType, ok := context[EnvironmentTypeKey]
	if !ok {
		return "", false, "", fmt.Errorf("environmentType is not set")
	}
	return manifestsPath, recursive, environmentType, nil
}

func getManifestFile(fs afero.Fs, filePath string) (stevedore.ManifestFile, error) {
	manifestFiles, err := getAllManifestFiles(fs, filePath)
	if err != nil {
		return stevedore.ManifestFile{}, err
	} else if len(manifestFiles) == 0 {
		return stevedore.ManifestFile{}, fmt.Errorf("unable to load manifest file %s", filePath)
	}
	return manifestFiles[0], nil
}

func getAllReleaseSpecifications(fs afero.Fs, dirPath string) (mappedReleases, error) {
	manifestFiles, err := getAllManifestFiles(fs, dirPath)
	releaseSpecs := map[string]stevedore.ReleaseSpecification{}
	if err != nil {
		return nil, err
	}
	for _, manifestFile := range manifestFiles {
		for _, spec := range manifestFile.Spec {
			releaseSpecs[spec.Release.Name] = spec
		}
	}
	return releaseSpecs, nil
}

func getDependentReleaseSpecifications(
	releaseSpecs stevedore.ReleaseSpecifications,
	knownReleases mappedReleases,
	visited mappedReleases,
) mappedReleases {

	for _, releaseSpec := range releaseSpecs {
		if _, ok := visited[releaseSpec.Release.Name]; !ok {
			visited[releaseSpec.Release.Name] = releaseSpec

			for _, dependentComponentName := range releaseSpec.DependsOn {
				if nextReleaseSpecification, ok := knownReleases[dependentComponentName]; ok {
					visited = getDependentReleaseSpecifications(stevedore.ReleaseSpecifications{nextReleaseSpecification}, knownReleases, visited)
				}
			}

		}
	}
	return visited
}

func getAllManifestFiles(fs afero.Fs, dirPath string) (stevedore.ManifestFiles, error) {
	yamlFiles, err := yaml.NewYamlFiles(fs, dirPath)

	if _, ok := err.(yaml.EmptyFolderError); ok {
		return stevedore.ManifestFiles{}, nil
	} else if err != nil {
		return nil, err
	}

	var fileErrors file.Errors
	manifests := stevedore.ManifestFiles{}
	for _, yamlFile := range yamlFiles {
		ok, err := yamlFile.Check(stevedore.KindStevedoreManifest, stevedore.ManifestCurrentVersion)
		if err != nil {
			fileErrors = append(fileErrors, file.Error{Filename: yamlFile.Name, Reason: err})
			continue
		}
		if ok {
			manifest, err := stevedore.NewManifest(yamlFile.Reader())
			if err != nil {
				fileErrors = append(fileErrors, file.Error{Filename: yamlFile.Name, Reason: err})
				continue
			}
			manifests = append(manifests, stevedore.ManifestFile{File: yamlFile.Name, Manifest: *manifest})
		}
	}

	if len(fileErrors) != 0 {
		return nil, fileErrors
	}

	return manifests, nil
}

// NewManifestProvider return new instance of ManifestProvider
func NewManifestProvider(fs afero.Fs) plugin.ManifestInterface {
	//TODO: remove this function
	return DefaultManifestProvider{fs: fs}
}

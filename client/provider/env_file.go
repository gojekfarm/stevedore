package provider

import (
	"github.com/gojek/stevedore/pkg/stevedore"
)

// EnvsFile wraps stevedore.Env with it corresponding file name
type EnvsFile struct {
	Name string
	stevedore.EnvSpecifications
}

// EnvsFiles is collection of EnvsFile
type EnvsFiles []EnvsFile

// Filter filters the manifests by removing non applicable ones
func (envsFiles EnvsFiles) Filter(context stevedore.Context) EnvsFiles {
	filteredEnvsFiles := EnvsFiles{}

	for _, envsFile := range envsFiles {
		applicableEnvs := stevedore.EnvSpecifications{}
		for _, env := range envsFile.EnvSpecifications {
			if env.IsApplicableFor(context) {
				applicableEnvs = append(applicableEnvs, env)
			}
		}

		if len(applicableEnvs) != 0 {
			filteredEnvsFiles = append(filteredEnvsFiles, EnvsFile{Name: envsFile.Name, EnvSpecifications: applicableEnvs})
		}
	}

	return filteredEnvsFiles
}

func (envsFiles EnvsFiles) extractEnvs() stevedore.EnvSpecifications {
	envs := make(stevedore.EnvSpecifications, 0, len(envsFiles))

	for _, envsFile := range envsFiles {
		envs = append(envs, envsFile.EnvSpecifications...)
	}

	return envs
}

// SortAndMerge will sort all envs based on a pre-determined order and merge all the env substitute into one
func (envsFiles EnvsFiles) SortAndMerge(envs stevedore.Substitute, labels stevedore.Labels) (stevedore.Substitute, error) {
	substitutes := stevedore.Substitute{}
	applicableEnvs := envsFiles.extractEnvs()
	applicableEnvs.Sort(labels)

	for _, env := range applicableEnvs {
		result, err := substitutes.Merge(env.Values)
		if err != nil {
			return nil, err
		}
		substitutes = result
	}

	return substitutes.Merge(envs)
}

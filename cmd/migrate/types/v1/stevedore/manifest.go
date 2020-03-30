package stevedore

import "github.com/gojek/stevedore/pkg/stevedore"

// Manifest is a Stevedore Manifest to wrap all configurations necessary for Stevedore to deploy
// Manifest is the struct representation of Stevedore yaml
type Manifest struct {
	Environments stevedore.Environments `json:"deployTo" yaml:"deployTo" validate:"required"`
	Applications Applications           `json:"applications" yaml:"applications" validate:"required,dive"`
}

func (manifest Manifest) matchers(contexts stevedore.Contexts, optimize bool) stevedore.Matchers {
	contextTypes := map[string]struct{}{}
	specificTypes := stevedore.Matchers{}

	for _, environment := range manifest.Environments {
		contextName := string(environment)
		if optimize {
			if index, ok := contexts.Find(contextName); ok {
				context := contexts[index]
				contextType := context.Type
				if contextType != "components" {
					contextTypes[contextType] = struct{}{}
					continue
				}
			}
		}
		specificTypes = append(specificTypes, stevedore.Conditions{stevedore.ConditionContextName: contextName})
	}

	result := stevedore.Matchers{}
	for contextType := range contextTypes {
		result = append(result, stevedore.Conditions{stevedore.ConditionContextType: contextType})
	}
	return append(result, specificTypes...)
}

// Convert convert Manifest as stevedore.Manifest
func (manifest Manifest) Convert(contexts stevedore.Contexts, optimize bool) stevedore.Manifest {
	matchers := manifest.matchers(contexts, optimize)
	return stevedore.Manifest{
		Kind:     stevedore.KindStevedoreManifest,
		Version:  stevedore.ManifestCurrentVersion,
		DeployTo: matchers,
		Spec:     manifest.Applications.ReleaseSpecifications(),
	}
}

package stevedore

import "github.com/gojek/stevedore/pkg/stevedore"

// Envs represents collection of Env
type Envs []Env

// Convert convert envs to newer stevedore envs format
func (envs Envs) Convert() stevedore.Env {
	specifications := stevedore.EnvSpecifications{}
	for _, ignore := range envs {
		specifications = append(specifications, ignore.Convert())
	}
	return stevedore.Env{
		Kind:    stevedore.KindStevedoreEnv,
		Version: stevedore.EnvCurrentVersion,
		Spec:    specifications,
	}
}

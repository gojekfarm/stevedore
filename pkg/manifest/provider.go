package manifest

import (
	"github.com/gojek/stevedore/pkg/stevedore"
)

// Provider is the interface which represents the contract of manifest plugins
type Provider interface {
	Manifests(map[string]string) (stevedore.ManifestFiles, error)
}

// ProviderImpl is a ProviderImpl
type ProviderImpl struct {
	Name     string
	Context  map[string]string
	Provider Provider
}

// MergeToContext will merge the given map to the context
func (p *ProviderImpl) MergeToContext(extraContext map[string]string) {
	for key, value := range extraContext {
		p.Context[key] = value
	}
}

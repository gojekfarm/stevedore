package config

import (
	"fmt"
)

// Provider is a provider for the configs and mounts section
type Provider interface {
	Fetch(
		context map[string]string,
		data interface{},
	) (map[string]interface{}, error)
}

// ProviderImpl is a ProviderImpl
type ProviderImpl struct {
	Name     string
	Context  map[string]string
	Provider Provider
}

// Providers represents the list of ProviderImpl
type Providers []ProviderImpl

// Fetch can be used to fetch configs from an external plugin
func (p Providers) Fetch(
	stevedoreCtx map[string]string,
	data map[string]interface{},
) (map[string]map[string]interface{}, error) {
	results := make(map[string]map[string]interface{}, len(data))

	pluginsMap := map[string]ProviderImpl{}

	for _, each := range p {
		pluginsMap[each.Name] = each
	}

	for rootKey, configs := range data {
		fetcher, ok := pluginsMap[rootKey]
		if !ok {
			return nil, fmt.Errorf("could not find a plugin binary for the config: %v", rootKey)
		}

		ctx := make(map[string]string)

		for k, v := range stevedoreCtx {
			ctx[k] = v
		}

		for k, v := range fetcher.Context {
			ctx[k] = v
		}

		result, err := fetcher.Provider.Fetch(ctx, configs)
		if err != nil {
			return nil, fmt.Errorf("error in fetching from provider: %v", err)
		}
		results[rootKey] = result
	}
	return results, nil
}

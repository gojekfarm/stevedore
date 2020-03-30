package stevedore

import (
	"github.com/gojek/stevedore/pkg/config"
	"github.com/gojek/stevedore/pkg/merger"
)

// Configs represents different store from where configurations can be fetched
type Configs map[string]interface{}

// Fetch returns substitute
func (configs Configs) Fetch(providers config.Providers, context Context) (Substitute, error) {
	contextMap, err := context.Map()
	if err != nil {
		return nil, err
	}

	//TODO: Fix: Non Deterministic Plugin merge for same config keys
	pluginConfigs, err := providers.Fetch(contextMap, configs)
	if err != nil {
		return nil, err
	}

	var pluginConfigList []map[string]interface{}
	for _, pluginConfig := range pluginConfigs {
		pluginConfigList = append(pluginConfigList, pluginConfig)
	}
	merge, err := merger.Merge(pluginConfigList...)
	return merge, err
}

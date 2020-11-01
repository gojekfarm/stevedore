package manifest

import (
	"github.com/gojek/stevedore/client/provider"
	"github.com/gojek/stevedore/pkg/config"
	"github.com/gojek/stevedore/pkg/manifest"
)

// NewManifests all stevedore configs in given folder
func NewManifests(
	environment config.Environment,
	contextProvider provider.ContextProvider,
	manifestProviderImpl manifest.ProviderImpl,
	overridesProvider provider.OverrideProvider,
	ignoreProvider provider.IgnoreProvider,
	envProvider provider.EnvProvider,
	reporter Reporter,
	providers config.Providers,
) (*Info, error) {
	ctx, err := contextProvider.Context()
	if err != nil {
		return nil, err
	}
	reporter.ReportContext(ctx)

	labels, err := contextProvider.Labels()
	if err != nil {
		return nil, err
	}

	envs, err := envProvider.Envs()
	if err != nil {
		return nil, err
	}

	filteredEnvs := envs.Filter(ctx)
	substitutes, err := filteredEnvs.SortAndMerge(environment.Fetch(), labels)
	if err != nil {
		return nil, err
	}

	reporter.ReportEnvs(filteredEnvs)

	ignores, err := ignoreProvider.Ignores()
	if err != nil {
		return nil, err
	}
	reporter.ReportIgnores(ignores)

	overrides, err := overridesProvider.Overrides()
	if err != nil {
		return nil, err
	}
	reporter.ReportOverrides(overrides)

	contextMap, err := ctx.Map()
	if err != nil {
		return nil, err
	}
	manifestProviderImpl.MergeToContext(contextMap)
	manifests, err := manifestProviderImpl.Provider.Manifests(manifestProviderImpl.Context)
	if err != nil {
		return nil, err
	}

	info, err := info(manifests, overrides, ctx, substitutes, ignores, providers, labels)
	if err != nil {
		return nil, err
	}
	reporter.ReportSkipped(info.Ignored)
	reporter.ReportManifest(info.ManifestFiles)

	return info, nil
}

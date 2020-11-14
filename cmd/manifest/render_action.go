package manifest

import (
	"io"

	"github.com/gojek/stevedore/cmd/cli"
)

// RenderAction to render manifest
type RenderAction struct {
	info Info
	out  io.Writer
}

// Do RenderAction to render manifest
func (action RenderAction) Do() (Info, error) {
	out := action.out
	if out == nil {
		out = cli.OutputStream()
	}

	for _, manifestFile := range action.info.ManifestFiles {
		cli.FPrintYaml(out, map[string]interface{}{"File": manifestFile.File})

		cli.FPrintYaml(out, map[string]interface{}{"DeployTo": manifestFile.Manifest.DeployTo})

		for _, releaseSpecification := range manifestFile.Manifest.Spec {
			overrides := releaseSpecification.Release.Overrides()
			substitute := releaseSpecification.SubstitutedVariables()

			cli.FPrintYaml(out, map[string]interface{}{"Manifest": releaseSpecification})

			if len(overrides.Spec) != 0 {
				cli.FPrintYaml(out, map[string]interface{}{"Used following Overrides": overrides.Spec})
			}

			if len(substitute) != 0 {
				cli.FPrintYaml(out, map[string]interface{}{"Used following variables": substitute})
			}
		}
	}
	return action.info, nil
}

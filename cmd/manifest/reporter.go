package manifest

import (
	"fmt"
	"github.com/gojek/stevedore/client/provider"

	"github.com/gojek/stevedore/cmd/cli"
	"github.com/gojek/stevedore/pkg/stevedore"
)

// Reporter reports the progress
type Reporter interface {
	ReportContext(stevedore.Context)
	ReportIgnores(ignores stevedore.Ignores)
	ReportOverrides(overrides stevedore.Overrides)
	ReportManifest(files stevedore.ManifestFiles)
	ReportSkipped(components stevedore.IgnoredReleases)
	ReportEnvs(files provider.EnvsFiles)
}

// DefaultReporter reports progress to cli
type DefaultReporter struct{}

// ReportContext prints the context information to cli
func (r DefaultReporter) ReportContext(ctx stevedore.Context) {
	cli.Info(fmt.Sprintf("Using context %s", ctx.Name))
}

// ReportIgnores prints the context information to cli
func (r DefaultReporter) ReportIgnores(ignores stevedore.Ignores) {
	count := len(ignores)
	if count == 0 {
		cli.Warn("No ignore rules found")
		return
	}
	cli.Info(fmt.Sprintf("Found %d ignore rule(s)", count))
}

// ReportEnvs prints the context information to cli
func (r DefaultReporter) ReportEnvs(files provider.EnvsFiles) {
	count := len(files)
	if count == 0 {
		cli.Warn("No envs found")
		return
	}
	cli.Info(fmt.Sprintf("Using %d env files(s):", count))

	for _, file := range files {
		cli.Info(fmt.Sprintf("  - %s", file.Name))
	}
}

// ReportOverrides prints the overrides information to cli
func (r DefaultReporter) ReportOverrides(overrides stevedore.Overrides) {
	count := len(overrides.Spec)
	if count == 0 {
		cli.Warn("No overrides found")
		return
	}
	cli.Info(fmt.Sprintf("Found %d override(s)", count))
}

// ReportManifest prints the manifests information to cli
func (r DefaultReporter) ReportManifest(files stevedore.ManifestFiles) {
	count := len(files)

	if count == 0 {
		cli.Warn("No manifest file found\n")
		return
	}

	cli.Info(fmt.Sprintf("Processing %d manifest file(s)", count))
	for _, file := range files {
		cli.Info(fmt.Sprintf("  - %s", file.File))
	}
}

// ReportSkipped prints the skipped components information to cli
func (r DefaultReporter) ReportSkipped(releases stevedore.IgnoredReleases) {
	count := len(releases)

	if count == 0 {
		return
	}

	cli.Info(fmt.Sprintf("%d releases(s) are ignored:", count))
	for _, release := range releases {
		cli.Info(fmt.Sprintf("  - '%s' is ignored. Reason: %s", release.Name, release.Reason))
	}
}

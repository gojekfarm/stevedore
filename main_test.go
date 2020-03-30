//+build integration

package main_test

import (
	"flag"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/cucumber/godog/gherkin"
	"github.com/gojek/stevedore/internal/cli/helm"
	"github.com/gojek/stevedore/internal/helpers"
)

func iHaveAKubernetesClusterWithNameAndVersion(clusterName, version string) error {
	return helpers.CreateCluster(clusterName, version, false)
}

func iHaveToInstallIntoMyCluster(_ string, _ string, _ string) error {
	return godog.ErrPending
}

func iHaveFollowingHelmRepos(helmRepos *gherkin.DataTable) error {
	repos := helm.NewRepos(helmRepos)
	return helpers.AddHelmRepos(repos)
}

func iRefreshHelmLocalCache() error {
	return helpers.UpdateHelmRepo()
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I have a kubernetes cluster with name "([^"]*)" and version "([^"]*)"$`, iHaveAKubernetesClusterWithNameAndVersion)
	s.Step(`^I have following helm repos:$`, iHaveFollowingHelmRepos)
	s.Step(`^I refresh helm local cache$`, iRefreshHelmLocalCache)
	s.Step(`^I have to install "([^"]*)" into cluster "([^"]*)" as "([^"]*)"$`, iHaveToInstallIntoMyCluster)
}

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	//Format: "progress", // can define default values
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

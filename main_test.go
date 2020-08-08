//+build integration

package main_test

import (
	"flag"
	"os"
	"testing"

	"github.com/cucumber/messages-go/v10"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/gojek/stevedore/internal/cli/helm"
	"github.com/gojek/stevedore/internal/helpers"
)

func iHaveAKubernetesClusterWithNameAndVersion(clusterName, version string) error {
	return helpers.CreateCluster(clusterName, version, false)
}

func iHaveToInstallIntoMyCluster(_ string, _ string, _ string) error {
	return godog.ErrPending
}

func iHaveFollowingHelmRepos(helmRepos *messages.PickleStepArgument_PickleTable) error {
	repos := helm.NewRepos(helmRepos)
	return helpers.AddHelmRepos(repos)
}

func iRefreshHelmLocalCache() error {
	return helpers.UpdateHelmRepo()
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {})
	ctx.Step(`^I have a kubernetes cluster with name "([^"]*)" and version "([^"]*)"$`, iHaveAKubernetesClusterWithNameAndVersion)
	ctx.Step(`^I have following helm repos:$`, iHaveFollowingHelmRepos)
	ctx.Step(`^I refresh helm local cache$`, iRefreshHelmLocalCache)
	ctx.Step(`^I have to install "([^"]*)" into cluster "([^"]*)" as "([^"]*)"$`, iHaveToInstallIntoMyCluster)
}

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opts.Paths = flag.Args()

	status := godog.TestSuite{
		Name:                 "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

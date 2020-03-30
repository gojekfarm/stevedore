package main

import (
	"github.com/gojek/stevedore/cmd"
	"github.com/gojek/stevedore/cmd/config"
	"github.com/gojek/stevedore/pkg/helm"
)

var version string
var build string
var helmVersion = "v2.16.9"
var configFilePath string

//go:generate ./scripts/mocks

func init() {
	helm.SetHelmVersion(helmVersion)
	helm.SetBuildMetadata("")
	config.SetBuildVersion(version, build)
	config.SetConfigFilePath(configFilePath)
}

func main() {
	cmd.Execute()
}

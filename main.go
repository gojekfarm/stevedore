package main

import (
	"github.com/gojek/stevedore/cmd"
	"github.com/gojek/stevedore/cmd/config"
)

var version string
var build string
var configFilePath string

//go:generate ./scripts/mocks

func init() {
	//helm.SetHelmVersion(helmVersion)
	//helm.SetBuildMetadata("")
	config.SetBuildVersion(version, build)
	config.SetConfigFilePath(configFilePath)
}

func main() {
	cmd.Execute()
}

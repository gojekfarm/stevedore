package cmd

import (
	"github.com/gojek/stevedore/cmd/dependency"
	"github.com/spf13/cobra"
)

var dependencyCmd = &cobra.Command{
	Use:   "dependency",
	Short: "Manage dependency for the manifest(s)",
}

func init() {
	buildCommand := dependency.NewBuildCommand()

	cobraCommand, err := buildCommand.CobraCommand(fs, &cfgFile, localStore)
	if err != nil {
		panic(err)
	}
	dependencyCmd.AddCommand(cobraCommand)
	rootCmd.AddCommand(dependencyCmd)
}

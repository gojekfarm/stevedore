package cmd

import (
	"github.com/gojek/stevedore/cmd/cli"
	"github.com/gojek/stevedore/cmd/manifest"
)

func init() {
	action := manifest.NewApplyCmd(fs, &cfgFile, true)
	command, err := action.CobraCommand()
	cli.DieIf(err, closePlugins)
	rootCmd.AddCommand(command)
}

package cmd

import (
	"github.com/gojek/stevedore/cmd/cli"
	"github.com/gojek/stevedore/cmd/manifest"
)

func init() {
	action := manifest.NewRenderCmd(fs, &cfgFile, false)
	command, err := action.CobraCommand()
	cli.DieIf(err, closePlugins)
	rootCmd.AddCommand(command)
}

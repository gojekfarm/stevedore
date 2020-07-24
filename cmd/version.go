package cmd

import (
	"fmt"

	"github.com/gojek/stevedore/pkg/helm"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:           "version",
	Short:         "Print the version",
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s version %s\n", rootCmd.Name(), rootCmd.Version)
		fmt.Printf("supported helm client version %s\n", helm.GetHelmVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

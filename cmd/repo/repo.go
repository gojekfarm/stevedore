package repo

import "github.com/spf13/cobra"

// defaultHelmRepoName is the default helm repo name to be used
// by stevedore. It can be set through the -ldflags at build time.
// If not set "chartmuseum" will be used.
var defaultHelmRepoName = "chartmuseum"

// AddRepoFlags add helm repo related flags to cobra command
func AddRepoFlags(cmd *cobra.Command, param *string) {
	cmd.PersistentFlags().StringVarP(param, "helm-repo-name", "r", defaultHelmRepoName, "Helm repo to which the charts need to be pushed")
}

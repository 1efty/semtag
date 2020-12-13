package cmd

import (
	"github.com/1efty/semtag/pkg/utils"
	"github.com/spf13/cobra"
)

var getCurrentCmd = &cobra.Command{
	Use:   "getcurrent",
	Short: "Returns the current version, based on the latest one.",
	Long: `Returns the current version, based on the latest one, if there are un-committed or
			un-staged changes, they will be reflected in the version, adding the number of
			pending commits, current branch and commit hash.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		initGit()
	},
	RunE: getCurrentAction,
}

func init() {
	rootCmd.AddCommand(getCurrentCmd)
}

func getCurrentAction(cmd *cobra.Command, args []string) error {
	utils.Info("Current tagged version: %s", repository.CurrentVersion.String())
	return nil
}

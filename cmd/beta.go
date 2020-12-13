package cmd

import (
	"github.com/1efty/semtag/pkg/utils"
	"github.com/spf13/cobra"
)

var betaCmd = &cobra.Command{
	Use:   "beta",
	Short: "Tags the current ref as a beta version, the tag will contain all the commits from the last final version.",
	PreRun: func(cmd *cobra.Command, args []string) {
		initGit()
	},
	RunE: betaAction,
}

func init() {
	rootCmd.AddCommand(betaCmd)
}

func betaAction(cmd *cobra.Command, args []string) error {
	v, err := bumpVersion(repository.LastVersion, Scope, "beta", Metadata)
	utils.CheckIfError(err)
	tagAction(repository, v.String(), Output)
	return nil
}

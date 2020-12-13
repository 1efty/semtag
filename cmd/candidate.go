package cmd

import (
	"github.com/1efty/semtag/pkg/utils"
	"github.com/spf13/cobra"
)

var candidateCmd = &cobra.Command{
	Use:   "candidate",
	Short: "Tags the current ref as a release candidate, the tag will contain all the commits from the last final version.",
	PreRun: func(cmd *cobra.Command, args []string) {
		initGit()
	},
	RunE: candidateAction,
}

func init() {
	rootCmd.AddCommand(candidateCmd)
}

func candidateAction(cmd *cobra.Command, args []string) error {
	v, err := bumpVersion(repository.LastVersion, Scope, "rc", Metadata)
	utils.CheckIfError(err)
	tagAction(repository, v.String(), Output)
	return nil
}

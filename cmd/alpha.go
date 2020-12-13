package cmd

import (
	"github.com/1efty/semtag/pkg/utils"
	"github.com/spf13/cobra"
)

var alphaCmd = &cobra.Command{
	Use:   "alpha",
	Short: "Tags the current ref as a alpha version, the tag will contain all the commits from the last final version.",
	PreRun: func(cmd *cobra.Command, args []string) {
		initGit()
	},
	RunE: alphaAction,
}

func init() {
	rootCmd.AddCommand(alphaCmd)
}

func alphaAction(cmd *cobra.Command, args []string) error {
	v, err := bumpVersion(repository.LastVersion, Scope, "alpha", Metadata)
	utils.CheckIfError(err)
	tagAction(repository, v.String(), Output)
	return nil
}

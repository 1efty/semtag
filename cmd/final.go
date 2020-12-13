package cmd

import (
	"github.com/1efty/semtag/pkg/utils"
	"github.com/spf13/cobra"
)

var finalCmd = &cobra.Command{
	Use:   "final",
	Short: "Tags the current ref as a final version, this only be done on the master branch.",
	PreRun: func(cmd *cobra.Command, args []string) {
		initGit()
	},
	RunE: finalAction,
}

func init() {
	rootCmd.AddCommand(finalCmd)
}

func finalAction(cmd *cobra.Command, args []string) error {
	v, err := bumpVersion(repository.LastVersion, Scope, "", "")
	utils.CheckIfError(err)
	tagAction(repository, v.String(), Output)
	return nil
}

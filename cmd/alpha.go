package cmd

import (
	"github.com/1efty/semtag/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(alphaCmd)
}

var alphaCmd = &cobra.Command{
	Use:   "alpha",
	Short: "Tags the current ref as a alpha version, the tag will contain all the commits from the last final version.",
	PreRun: func(cmd *cobra.Command, args []string) {
		initGit()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := alphaAction(); err != nil {
			return err
		}
		return nil
	},
}

func alphaAction() error {
	v, err := lib.BumpVersion(lastVersion, Scope, "alpha", Metadata)
	lib.CheckIfError(err)
	tagAction(repository, v.String(), Output)
	return nil
}

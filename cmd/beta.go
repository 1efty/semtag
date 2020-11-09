package cmd

import (
	"github.com/1efty/semtag/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(betaCmd)
}

var betaCmd = &cobra.Command{
	Use:   "beta",
	Short: "Tags the current ref as a beta version, the tag will contain all the commits from the last final version.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := betaAction(); err != nil {
			return err
		}
		return nil
	},
}

func betaAction() error {
	v, err := lib.BumpVersion(lastVersion, Scope, "beta", Metadata)
	lib.CheckIfError(err)
	lib.CreateTag(repository, v.String())
	return nil
}

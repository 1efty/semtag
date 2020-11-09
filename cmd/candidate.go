package cmd

import (
	"github.com/1efty/semtag/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(candidateCmd)
}

var candidateCmd = &cobra.Command{
	Use:   "candidate",
	Short: "Tags the current ref as a release candidate, the tag will contain all the commits from the last final version.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := candidateAction(); err != nil {
			return err
		}
		return nil
	},
}

func candidateAction() error {
	v, err := lib.BumpVersion(lastVersion, Scope, "rc", Metadata)
	lib.CheckIfError(err)
	lib.CreateTag(repository, v.String())
	return nil
}

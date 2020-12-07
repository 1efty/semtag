package cmd

import (
	"github.com/1efty/semtag/lib"
	"github.com/spf13/cobra"
)

var candidateCmd *cobra.Command

var _ = RegisterCommandVar(func() {
	candidateCmd = &cobra.Command{
		Use:   "candidate",
		Short: "Tags the current ref as a release candidate, the tag will contain all the commits from the last final version.",
		PreRun: func(cmd *cobra.Command, args []string) {
			initGit()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := candidateAction(); err != nil {
				return err
			}
			return nil
		},
	}
})

var _ = RegisterCommandInit(func() {
	rootCmd.AddCommand(candidateCmd)
})

func candidateAction() error {
	v, err := bumpVersion(repository.LastVersion, Scope, "rc", Metadata)
	lib.CheckIfError(err)
	tagAction(repository, v.String(), Output)
	return nil
}

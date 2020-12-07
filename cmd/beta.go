package cmd

import (
	"github.com/1efty/semtag/pkg/utils"
	"github.com/spf13/cobra"
)

var betaCmd *cobra.Command

var _ = RegisterCommandVar(func() {
	betaCmd = &cobra.Command{
		Use:   "beta",
		Short: "Tags the current ref as a beta version, the tag will contain all the commits from the last final version.",
		PreRun: func(cmd *cobra.Command, args []string) {
			initGit()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := betaAction(); err != nil {
				return err
			}
			return nil
		},
	}
})

var _ = RegisterCommandInit(func() {
	rootCmd.AddCommand(betaCmd)
})

func betaAction() error {
	v, err := bumpVersion(repository.LastVersion, Scope, "beta", Metadata)
	utils.CheckIfError(err)
	tagAction(repository, v.String(), Output)
	return nil
}

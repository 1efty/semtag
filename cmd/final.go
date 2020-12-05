package cmd

import (
	"github.com/1efty/semtag/lib"
	"github.com/spf13/cobra"
)

var finalCmd *cobra.Command

var _ = RegisterCommandVar(func() {
	finalCmd = &cobra.Command{
		Use:   "final",
		Short: "Tags the current ref as a final version, this only be done on the master branch.",
		PreRun: func(cmd *cobra.Command, args []string) {
			initGit()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := finalAction(); err != nil {
				return err
			}
			return nil
		},
	}
})

var _ = RegisterCommandInit(func() {
	rootCmd.AddCommand(finalCmd)
})

func finalAction() error {
	v, err := lib.BumpVersion(lastVersion, Scope, "", "")
	lib.CheckIfError(err)
	tagAction(repository, v.String(), Output)
	return nil
}

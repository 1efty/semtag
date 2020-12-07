package cmd

import (
	"fmt"

	"github.com/1efty/semtag/pkg/utils"
	"github.com/spf13/cobra"
)

var getCurrentCmd *cobra.Command

var _ = RegisterCommandVar(func() {
	getCurrentCmd = &cobra.Command{
		Use:   "getcurrent",
		Short: "Returns the current version, based on the latest one.",
		Long: `Returns the current version, based on the latest one, if there are un-committed or
			un-staged changes, they will be reflected in the version, adding the number of
			pending commits, current branch and commit hash.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			initGit()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := getCurrentAction(); err != nil {
				return err
			}
			return nil
		},
	}
})

var _ = RegisterCommandInit(func() {
	rootCmd.AddCommand(getCurrentCmd)
})

func getCurrentAction() error {
	utils.Info(fmt.Sprintf("Current tagged version: %s", repository.CurrentVersion.String()))
	return nil
}

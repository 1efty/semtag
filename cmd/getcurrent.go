package cmd

import (
	"fmt"

	"github.com/1efty/semtag/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCurrentCmd)
}

var getCurrentCmd = &cobra.Command{
	Use:   "getcurrent",
	Short: "Returns the current version, based on the latest one.",
	Long: `Returns the current version, based on the latest one, if there are un-committed or
			un-staged changes, they will be reflected in the version, adding the number of
			pending commits, current branch and commit hash.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := getCurrentAction(); err != nil {
			return err
		}
		return nil
	},
}

func getCurrentAction() error {
	lib.Info(fmt.Sprintf("Current tagged version: %s", currentVersion.String()))
	return nil
}

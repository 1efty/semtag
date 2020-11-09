package cmd

import (
	"github.com/1efty/semtag/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(finalCmd)
}

var finalCmd = &cobra.Command{
	Use:   "final",
	Short: "Tags the current ref as a final version, this only be done on the master branch.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := finalAction(); err != nil {
			return err
		}
		return nil
	},
}

func finalAction() error {
	v, err := lib.BumpVersion(lastVersion, Scope, "", "")
	lib.CheckIfError(err)

	lib.CreateTag(repository, v.String())

	return nil
}

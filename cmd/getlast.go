package cmd

import (
	"fmt"

	"github.com/1efty/semtag/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getLastCmd)
}

var getLastCmd = &cobra.Command{
	Use:   "getlast",
	Short: "Returns the latest tagged version.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := getLastAction(); err != nil {
			return err
		}
		return nil
	},
}

func getLastAction() error {
	lib.Info(fmt.Sprintf("Last tagged version: %v", lastVersion))
	return nil
}

package cmd

import (
	"fmt"

	"github.com/1efty/semtag/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns both current, final, and last tagged versions.",
	PreRun: func(cmd *cobra.Command, args []string) {
		initGit()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := getAction(); err != nil {
			return err
		}
		return nil
	},
}

func getAction() error {
	lib.Info(fmt.Sprintf("Current final version: %v", finalVersion))
	lib.Info(fmt.Sprintf("Last tagged version: %v", lastVersion))
	return nil
}

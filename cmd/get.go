package cmd

import (
	"github.com/1efty/semtag/pkg/utils"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns both current, final, and last tagged versions.",
	PreRun: func(cmd *cobra.Command, args []string) {
		initGit()
	},
	RunE: getAction,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func getAction(cmd *cobra.Command, args []string) error {
	utils.Info("Current final version: %v", repository.FinalVersion.String())
	utils.Info("Last tagged version: %v", repository.LastVersion.String())
	return nil
}

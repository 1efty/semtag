package cmd

import (
	"github.com/1efty/semtag/pkg/utils"
	"github.com/spf13/cobra"
)

var getFinalCmd = &cobra.Command{
	Use:   "getfinal",
	Short: "Returns latest tagged final version.",
	PreRun: func(cmd *cobra.Command, args []string) {
		initGit()
	},
	RunE: getFinalAction,
}

func init() {
	rootCmd.AddCommand(getFinalCmd)
}

func getFinalAction(cmd *cobra.Command, args []string) error {
	utils.Info("Final tagged version: %v", repository.FinalVersion.String())
	return nil
}

package cmd

import (
	"fmt"

	"github.com/1efty/semtag/pkg/utils"
	"github.com/spf13/cobra"
)

var getLastCmd = &cobra.Command{
	Use:   "getlast",
	Short: "Returns the latest tagged version.",
	PreRun: func(cmd *cobra.Command, args []string) {
		initGit()
	},
	RunE: getLastAction,
}

func init() {
	rootCmd.AddCommand(getLastCmd)
}

func getLastAction(cmd *cobra.Command, args []string) error {
	utils.Info(fmt.Sprintf("Last tagged version: %v", repository.LastVersion.String()))
	return nil
}

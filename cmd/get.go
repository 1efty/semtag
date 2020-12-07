package cmd

import (
	"fmt"

	"github.com/1efty/semtag/pkg/utils"
	"github.com/spf13/cobra"
)

var getCmd *cobra.Command

var _ = RegisterCommandVar(func() {
	getCmd = &cobra.Command{
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
})

var _ = RegisterCommandInit(func() {
	rootCmd.AddCommand(getCmd)
})

func getAction() error {
	utils.Info(fmt.Sprintf("Current final version: %v", repository.FinalVersion.String()))
	utils.Info(fmt.Sprintf("Last tagged version: %v", repository.LastVersion.String()))
	return nil
}

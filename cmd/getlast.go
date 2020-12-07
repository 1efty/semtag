package cmd

import (
	"fmt"

	"github.com/1efty/semtag/lib"
	"github.com/spf13/cobra"
)

var getLastCmd *cobra.Command

var _ = RegisterCommandVar(func() {
	getLastCmd = &cobra.Command{
		Use:   "getlast",
		Short: "Returns the latest tagged version.",
		PreRun: func(cmd *cobra.Command, args []string) {
			initGit()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := getLastAction(); err != nil {
				return err
			}
			return nil
		},
	}
})

var _ = RegisterCommandInit(func() {
	rootCmd.AddCommand(getLastCmd)
})

func getLastAction() error {
	lib.Info(fmt.Sprintf("Last tagged version: %v", repository.LastVersion.String()))
	return nil
}

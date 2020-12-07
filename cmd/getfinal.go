package cmd

import (
	"fmt"

	"github.com/1efty/semtag/lib"
	"github.com/spf13/cobra"
)

var getFinalCmd *cobra.Command

var _ = RegisterCommandVar(func() {
	getFinalCmd = &cobra.Command{
		Use:   "getfinal",
		Short: "Returns latest tagged final version.",
		PreRun: func(cmd *cobra.Command, args []string) {
			initGit()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := getFinalAction(); err != nil {
				return err
			}
			return nil
		},
	}
})

var _ = RegisterCommandInit(func() {
	rootCmd.AddCommand(getFinalCmd)
})

func getFinalAction() error {
	lib.Info(fmt.Sprintf("Final tagged version: %v", repository.FinalVersion.String()))
	return nil
}

package cmd

import (
	"github.com/1efty/semtag/lib"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getFinalCmd)
}

var getFinalCmd = &cobra.Command{
	Use:   "getfinal",
	Short: "Returns latest tagged final version.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := getFinalAction(); err != nil {
			return err
		}
		return nil
	},
}

func getFinalAction() error {
	lib.Info(finalVersion.String())
	return nil
}

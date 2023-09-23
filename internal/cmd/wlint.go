package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "wlint",
		Short: "Linter for writing projects",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func AddCommand(command *cobra.Command) error {
	rootCmd.AddCommand(command)
	return nil
}

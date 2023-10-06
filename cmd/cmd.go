package cmd

import (
	"github.com/spf13/cobra"

	icmd "github.com/snewell/wlint-go/internal/cmd"
)

var (
	Purifier string
)

func AddCommonFlags(command *cobra.Command) {
	command.PersistentFlags().StringVarP(&Purifier, "purifier", "p", "", "Text purifier to use")
}

func AddCommand(command *cobra.Command) error {
	return icmd.AddCommand(command)
}

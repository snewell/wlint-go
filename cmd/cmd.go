package cmd

import (
	"github.com/spf13/cobra"

	icmd "github.com/snewell/wlint-go/internal/cmd"
)

func AddCommand(command *cobra.Command) error {
	return icmd.AddCommand(command)
}

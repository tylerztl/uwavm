package cmd

import (
	"github.com/spf13/cobra"
)

var contractQueryCmd *cobra.Command

const queryCmdName = "query"

func QueryCmd() *cobra.Command {
	contractQueryCmd = &cobra.Command{
		Use:       queryCmdName,
		Short:     "Query the specified wasm contract.",
		Long:      "Query the specified wasm contract.",
		ValidArgs: []string{"1"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return contractInvokeOrQuery(cmd, args, contractMethod)
		},
	}
	flagList := []string{
		"name",
		"language",
		"method",
		"args",
		"caller",
	}
	attachFlags(contractQueryCmd, flagList)

	return contractQueryCmd
}

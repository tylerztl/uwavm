package cmd

import (
	"errors"
	"fmt"

	"github.com/BeDreamCoder/uwavm/bridge"
	"github.com/spf13/cobra"
)

var contractInvokeCmd *cobra.Command

const invokeCmdName = "invoke"

func InvokeCmd() *cobra.Command {
	contractInvokeCmd = &cobra.Command{
		Use:       invokeCmdName,
		Short:     "Invoke the specified wasm contract.",
		Long:      "Invoke the specified wasm contract.",
		ValidArgs: []string{"1"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return contractInvokeOrQuery(cmd, args, invokeCmdName)
		},
	}
	flagList := []string{
		"name",
		"language",
		"args",
		"caller",
	}
	attachFlags(contractInvokeCmd, flagList)

	return contractInvokeCmd
}

func contractInvokeOrQuery(cmd *cobra.Command, args []string, method string) error {
	if err := checkContractCmdParams(cmd); err != nil {
		return err
	}
	vm, ok := bridge.GetBridge(nil).GetVirtualMachine("wasm")
	if !ok {
		return errors.New("not found VirtualMachine name wasm")
	}

	if resp, err := vm.InvokeContract(method, makeInvokeOrQueryArgs()); err != nil {
		return err
	} else {
		fmt.Println("Status:", resp.GetStatus())
		fmt.Println("Message:", resp.GetMessage())
		fmt.Println("Bdoy:", string(resp.GetBody()))
		return nil
	}
}

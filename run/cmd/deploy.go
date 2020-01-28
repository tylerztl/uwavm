package cmd

import (
	"errors"
	"fmt"

	"github.com/BeDreamCoder/uwavm/bridge"
	"github.com/spf13/cobra"
)

var contractDeployCmd *cobra.Command

const deployCmdName = "deploy"

func DeployCmd() *cobra.Command {
	contractDeployCmd = &cobra.Command{
		Use:       deployCmdName,
		Short:     "Deploy the specified wasm contract.",
		Long:      "Deploy the specified wasm contract.",
		ValidArgs: []string{"1"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return contractDeploy(cmd, args)
		},
	}
	flagList := []string{
		"name",
		"language",
		"args",
		"path",
		"caller",
	}
	attachFlags(contractDeployCmd, flagList)

	return contractDeployCmd
}

func contractDeploy(cmd *cobra.Command, args []string) error {
	if err := checkContractCmdParams(cmd); err != nil {
		return err
	}
	vm, ok := bridge.GetBridge(nil).GetVirtualMachine("wasm")
	if !ok {
		return errors.New("not found VirtualMachine name wasm")
	}

	if resp, resourceUsed, err := vm.DeployContract(makeDeployArgs()); err != nil {
		return err
	} else {
		fmt.Println("Status:", resp.GetStatus())
		fmt.Println("Message:", resp.GetMessage())
		fmt.Println("Bdoy:", string(resp.GetBody()))
		fmt.Println("Gas:", resourceUsed.TotalGas())
		return nil
	}
}

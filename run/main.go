package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	cmdpkg "github.com/BeDreamCoder/uwavm/run/cmd"
	_ "github.com/BeDreamCoder/uwavm/vm/interpreter"
	"github.com/spf13/cobra"
)

// The main command describes the service and
// defaults to printing the help message.
var mainCmd = &cobra.Command{
	Use:   "uwavm",
	Short: "Decode wasm binary files.",
	Long:  "Decode wasm binary files that compiled by golang, javascript, c/c++, rust",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmdpkg.InitCmd(cmd, args)
	},
}

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Operate a contract: deploy|invoke|query|list.",
	Long:  "Operate a contract: deploy|invoke|query|list.",
}

// Cmd returns the cobra command for Chaincode
func ContractCmd() *cobra.Command {
	contractCmd.AddCommand(cmdpkg.DeployCmd())
	contractCmd.AddCommand(cmdpkg.InvokeCmd())
	contractCmd.AddCommand(cmdpkg.QueryCmd())

	return contractCmd
}

func makeDeployArgs(modulePath string) map[string][]byte {
	codebuf, err := ioutil.ReadFile(modulePath)
	if err != nil {
		panic(err)
	}

	args := map[string][]byte{
		"initSupply": []byte("1000000"),
	}
	argsbuf, _ := json.Marshal(args)
	return map[string][]byte{
		"contract_name": []byte("erc20"),
		"contract_code": codebuf,
		"language":      []byte("go"),
		"args":          argsbuf,
		"caller":        []byte("alice"),
	}
}

func makeInvokeArgs() map[string][]byte {
	args := map[string][]byte{
		"action":  []byte("balanceOf"),
		"address": []byte("alice"),
	}
	argsbuf, _ := json.Marshal(args)
	return map[string][]byte{
		"contract_name": []byte("erc20"),
		"language":      []byte("go"),
		"args":          argsbuf,
		"caller":        []byte("alice"),
	}
}

func main() {
	// Define command-line flags that are valid for all commands and
	// subcommands.
	mainCmd.AddCommand(ContractCmd())

	// On failure Cobra prints the usage message and error string, so we only
	// need to exit with a non-0 status
	if mainCmd.Execute() != nil {
		os.Exit(1)
	}
}

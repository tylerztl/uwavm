package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/BeDreamCoder/uwavm/bridge"
	"github.com/BeDreamCoder/uwavm/common/db/leveldb"
	"github.com/BeDreamCoder/uwavm/common/log"
	"github.com/BeDreamCoder/uwavm/vm"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	contractName   string
	contractLang   string
	contractMethod string
	contractArgs   string
	contractPath   string
	contractCaller string
)

var flags *pflag.FlagSet

func InitCmd(cmd *cobra.Command, args []string) {
	db := leveldb.NewProvider().GetDBHandle("uwavm")
	bridge := bridge.GetBridge(db)
	vm := vm.NewVMManager(db, bridge)
	bridge.RegisterExecutor("wasm", vm)
}

func init() {
	resetFlags()
}

// Explicitly define a method to facilitate tests
func resetFlags() {
	flags = &pflag.FlagSet{}

	flags.StringVarP(&contractName, "name", "n", "",
		fmt.Sprint("Name of the contract"))
	flags.StringVarP(&contractLang, "language", "l", "go",
		fmt.Sprintf("Language the contract is written in"))
	flags.StringVarP(&contractMethod, "method", "m", "invoke",
		fmt.Sprintf("Invoke contract method name"))
	flags.StringVarP(&contractArgs, "args", "a", "{}",
		fmt.Sprintf("Constructor message for the contract initialize args in JSON format"))
	flags.StringVarP(&contractPath, "path", "p", "",
		fmt.Sprintf("Path to wasm binary files"))
	flags.StringVarP(&contractCaller, "caller", "c", "",
		fmt.Sprint("Contract caller name"))
}

func attachFlags(cmd *cobra.Command, names []string) {
	cmdFlags := cmd.Flags()
	for _, name := range names {
		if flag := flags.Lookup(name); flag != nil {
			cmdFlags.AddFlag(flag)
		} else {
			log.GetLogger().Error("Could not find flag  to attach to command", "flag", name, "cmd", cmd.Name())
		}
	}
}

func checkContractCmdParams(cmd *cobra.Command) error {
	if contractName == "" {
		return errors.Errorf("must provide contract name")
	}

	if contractCaller == "" {
		return errors.Errorf("must provide contract caller")
	}

	if cmd.Name() == deployCmdName {
		if contractPath == "" {
			return errors.Errorf("must provide contract wasm file path")
		}
	} else {
		if contractMethod == "" {
			return errors.Errorf("must provide contract method name")
		}
	}

	if contractArgs != "{}" {
		var f map[string]string
		err := json.Unmarshal([]byte(contractArgs), &f)
		if err != nil {
			return errors.Wrap(err, "contract argument error")
		}
		m := make(map[string][]byte)
		for k := range f {
			m[k] = []byte(f[k])
		}
		fmtArgs, err := json.Marshal(m)
		if err != nil {
			return errors.Wrap(err, "contract argument error")
		}
		contractArgs = string(fmtArgs)
	}

	return nil
}

func makeDeployArgs() map[string][]byte {
	codebuf, err := ioutil.ReadFile(contractPath)
	if err != nil {
		panic(err)
	}

	return map[string][]byte{
		"contract_name": []byte(contractName),
		"contract_code": codebuf,
		"language":      []byte(contractLang),
		"args":          []byte(contractArgs),
		"caller":        []byte(contractCaller),
	}
}

func makeInvokeOrQueryArgs() map[string][]byte {
	return map[string][]byte{
		"contract_name": []byte(contractName),
		"language":      []byte(contractLang),
		"args":          []byte(contractArgs),
		"caller":        []byte(contractCaller),
	}
}

package wasm

import (
	"encoding/json"
	"errors"

	"github.com/BeDreamCoder/uwavm/bridge"
	"github.com/BeDreamCoder/uwavm/common/db"
	"github.com/BeDreamCoder/uwavm/common/util"
	"github.com/BeDreamCoder/uwavm/contract/go/pb"
	"github.com/BeDreamCoder/uwavm/wasm/vm"
	_ "github.com/BeDreamCoder/uwavm/wasm/vm"
	log "github.com/inconshreveable/log15"
)

// VMManager manages wasm contracts, include deploy contracts, instance wasm virtual machine, etc...
type VMManager struct {
	db     db.Database
	vmimpl vm.InstanceCreator
	bridge *bridge.Bridge
}

// New instances a new VMManager
func NewVMManager(db db.Database, bridge *bridge.Bridge) *VMManager {
	return &VMManager{
		db:     db,
		bridge: bridge,
	}
}

// RegisterSyscallService implements bridge.Executor
func (v *VMManager) RegisterSyscallService(syscall *bridge.SyscallService) {
	vmimpl, err := vm.Open("uwavm", syscall, v.db)
	if err != nil {
		panic(err)
	}
	v.vmimpl = vmimpl
}

// NewInstance implements bridge.Executor
func (v *VMManager) NewCreatorInstance(ctx *bridge.ContractState) (bridge.Instance, error) {
	ins, err := v.vmimpl.CreateInstance(ctx)
	if err != nil {
		return nil, err
	}
	return &bridgeInstance{
		ctx:        ctx,
		vmInstance: ins,
	}, nil
}

// TODO:校验名字
func (v *VMManager) verifyContractName(name string) error {
	return nil
}

// DeployContract deploy contract and initialize contract
func (v *VMManager) DeployContract(args map[string][]byte) (*pb.Response, error) {
	name := args["contract_name"]
	if name == nil {
		return nil, errors.New("bad contract name")
	}
	contractName := string(name)
	err := v.verifyContractName(contractName)
	if err != nil {
		return nil, err
	}

	code := args["contract_code"]
	if code == nil {
		return nil, errors.New("missing contract code")
	}

	language := args["language"]
	if language == nil {
		return nil, errors.New("missing contract language")
	}

	initArgsBuf := args["args"]
	if initArgsBuf == nil {
		return nil, errors.New("missing args field in args")
	}
	var initArgs map[string][]byte
	if err = json.Unmarshal(initArgsBuf, &initArgs); err != nil {
		return nil, err
	}

	caller := args["caller"]
	if caller == nil {
		return nil, errors.New("missing contract caller")
	}

	if err = v.db.Put(util.ContractCodeKey(contractName), code); err != nil {
		return nil, err
	}
	if err = v.db.Put(util.ContractCodeDescKey(contractName), language); err != nil {
		return nil, err
	}

	state := &bridge.ContractState{
		ContractName: contractName,
		Language:     string(language),
		Caller:       string(caller),
	}

	out, err := v.invokeContract(state, bridge.InitMethod, initArgs)
	if err != nil {
		if _, ok := err.(*bridge.ContractError); !ok {
			v.vmimpl.RemoveCache(contractName)
		}
		log.Error("call contract initialize method error", "error", err, "contract", contractName)
		return nil, err
	}
	return out, nil
}

func (v *VMManager) InvokeContract(method string, args map[string][]byte) (*pb.Response, error) {
	name := args["contract_name"]
	if name == nil {
		return nil, errors.New("bad contract name")
	}
	contractName := string(name)
	err := v.verifyContractName(contractName)
	if err != nil {
		return nil, err
	}

	caller := args["caller"]
	if caller == nil {
		return nil, errors.New("missing contract caller")
	}

	language := args["language"]
	if language == nil {
		return nil, errors.New("missing contract language")
	}

	argsBuf := args["args"]
	if argsBuf == nil {
		return nil, errors.New("missing args field in args")
	}
	var invokeArgs map[string][]byte
	if err = json.Unmarshal(argsBuf, &invokeArgs); err != nil {
		return nil, err
	}

	state := &bridge.ContractState{
		ContractName: contractName,
		Language:     string(language),
		Caller:       string(caller),
	}

	out, err := v.invokeContract(state, method, invokeArgs)
	if err != nil {
		if _, ok := err.(*bridge.ContractError); !ok {
			v.vmimpl.RemoveCache(contractName)
		}
		log.Error("call contract initialize method error", "error", err, "contract", contractName)
		return nil, err
	}
	return out, nil
}

func (v *VMManager) invokeContract(state *bridge.ContractState, method string, args map[string][]byte) (*pb.Response, error) {
	vm, ok := v.bridge.GetVirtualMachine("wasm")
	if !ok {
		return nil, errors.New("wasm vm not registered")
	}

	ctx, err := vm.NewContext(state)
	if err != nil {
		return nil, err
	}
	out, err := ctx.Invoke(method, args)
	if err != nil {
		return nil, err
	}
	return out, nil
}

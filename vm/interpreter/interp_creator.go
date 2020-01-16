package interpreter

import (
	"errors"
	"fmt"

	"github.com/BeDreamCoder/uwavm/bridge"
	"github.com/BeDreamCoder/uwavm/common/db"
	"github.com/BeDreamCoder/uwavm/common/util"
	"github.com/BeDreamCoder/uwavm/vm"
	"github.com/BeDreamCoder/uwavm/wasm/exec"
	gowasm "github.com/BeDreamCoder/uwavm/wasm/runtime/go"
)

type interpCreator struct {
	chd            vm.CodeHandle
	db             db.Database
	syscallService *bridge.SyscallService
}

func newInterpCreator(syscallService *bridge.SyscallService, db db.Database) (vm.InstanceCreator, error) {
	creator := &interpCreator{
		syscallService: syscallService,
		db:             db,
	}
	creator.chd = vm.NewCodeManager(creator.makeExecCode)
	return creator, nil
}

func (x *interpCreator) makeExecCode(contractName string) (exec.WasmExec, error) {
	codebuf, err := x.GetContractCode(contractName)
	if err != nil {
		return nil, err
	}
	resolver := exec.NewMultiResolver(
		gowasm.NewResolver(),
		newSyscallResolver(x.syscallService))
	return exec.NewInterpCode(codebuf, resolver)
}

func (x *interpCreator) CreateInstance(ctx *bridge.ContractState) (bridge.Instance, error) {
	code, err := x.chd.GetExecCode(ctx.ContractName)
	if err != nil {
		return nil, err
	}
	return createInstance(ctx, code)
}

func (x *interpCreator) RemoveCache(contractName string) {
	x.chd.RemoveCode(contractName)
}

func (x *interpCreator) GetContractCode(name string) ([]byte, error) {
	codebuf, err := x.db.Get(util.ContractCodeKey(name))
	if err != nil {
		return nil, fmt.Errorf("get contract code for '%s' error:%s", name, err)
	}
	if len(codebuf) == 0 {
		return nil, errors.New("empty wasm code")
	}
	return codebuf, nil
}

func init() {
	vm.Register("uwavm", newInterpCreator)
}

package vm

import (
	"errors"
	"fmt"
	"github.com/BeDreamCoder/uwavm/bridge"
	"github.com/BeDreamCoder/uwavm/common/db"
	"github.com/BeDreamCoder/uwavm/common/util"
	"github.com/BeDreamCoder/uwavm/exec"
	gowasm "github.com/BeDreamCoder/uwavm/runtime/go"
)

type interpCreator struct {
	cm             *codeManager
	db             db.Database
	syscallService *bridge.SyscallService
}

func newInterpCreator(syscallService *bridge.SyscallService, db db.Database) (InstanceCreator, error) {
	creator := &interpCreator{
		syscallService: syscallService,
		db:             db,
	}
	creator.cm = newCodeManager(creator.makeExecCode)
	return creator, nil
}

func (x *interpCreator) makeExecCode(contractName string) (exec.Code, error) {
	codebuf, err := x.GetContractCode(contractName)
	if err != nil {
		return nil, err
	}
	resolver := exec.NewMultiResolver(
		gowasm.NewResolver(),
		newSyscallResolver(x.syscallService))
	return exec.NewInterpCode(codebuf, resolver)
}

func (x *interpCreator) CreateInstance(ctx *bridge.ContractState) (Instance, error) {
	code, err := x.cm.GetExecCode(ctx.ContractName)
	if err != nil {
		return nil, err
	}
	return createInstance(ctx, code)
}

func (x *interpCreator) RemoveCache(contractName string) {
	x.cm.RemoveCode(contractName)
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
	Register("uwavm", newInterpCreator)
}

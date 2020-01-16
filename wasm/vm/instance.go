package vm

import (
	"errors"

	"github.com/BeDreamCoder/uwavm/bridge"
	"github.com/BeDreamCoder/uwavm/common/log"
	"github.com/BeDreamCoder/uwavm/exec"
	gowasm "github.com/BeDreamCoder/uwavm/runtime/go"
)

type vmInstance struct {
	bridgeCtx *bridge.ContractState
	execCtx   exec.Context
}

func createInstance(ctx *bridge.ContractState, code *contractCode) (Instance, error) {
	execCtx, err := code.ExecCode.NewContext(exec.DefaultContextConfig())
	if err != nil {
		log.GetLogger().Error("create contract context error", "error", err, "contract", ctx.ContractName)
		return nil, err
	}
	switch ctx.Language {
	case "go":
		gowasm.RegisterRuntime(execCtx)
	}
	execCtx.SetUserData(contextIDKey, ctx.ID)
	instance := &vmInstance{
		bridgeCtx: ctx,
		execCtx:   execCtx,
	}
	return instance, nil
}

func (x *vmInstance) Exec(function string) error {
	mem := x.execCtx.Memory()
	if mem == nil {
		return errors.New("bad contract, no memory")
	}
	var args []int64
	// go's entry function expects argc and argv these two arguments
	if x.bridgeCtx.Language == "go" {
		args = []int64{0, 0}
	}
	_, err := x.execCtx.Exec(function, args)
	if err != nil {
		log.GetLogger().Error("exec contract error", "error", err, "contract", x.bridgeCtx.ContractName)
	}
	return err
}

func (x *vmInstance) Release() {
	x.execCtx.Release()
}

func (x *vmInstance) Abort(msg string) {
	exec.Throw(exec.NewTrap(msg))
}

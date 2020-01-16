package wasm

import (
	"errors"

	"github.com/BeDreamCoder/uwavm/bridge"
	"github.com/BeDreamCoder/uwavm/wasm/vm"
)

type bridgeInstance struct {
	ctx        *bridge.ContractState
	vmInstance vm.Instance
}

func (v *bridgeInstance) guessEntry() (string, error) {
	switch v.ctx.Language {
	case "go":
		return "run", nil
	case "c":
		return "_" + v.ctx.Method, nil
	default:
		return "", errors.New("bad runtime")
	}
}

func (v *bridgeInstance) getEntry() (string, error) {
	return v.guessEntry()
}

func (v *bridgeInstance) Exec() error {
	entry, err := v.getEntry()
	if err != nil {
		return err
	}
	return v.vmInstance.Exec(entry)
}

func (v *bridgeInstance) Release() {
	v.vmInstance.Release()
}

func (v *bridgeInstance) Abort(msg string) {
	v.vmInstance.Abort(msg)
}

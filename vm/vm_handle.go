package vm

import (
	"errors"

	"github.com/BeDreamCoder/uwavm/bridge"
	"github.com/BeDreamCoder/uwavm/vm/gas"
)

type vmHandle struct {
	ctx        *bridge.ContractState
	vmInstance bridge.Instance
}

func (v *vmHandle) guessEntry() (string, error) {
	switch v.ctx.Language {
	case "go":
		return "run", nil
	case "c":
		return "_" + v.ctx.Method, nil
	default:
		return "", errors.New("bad runtime")
	}
}

func (v *vmHandle) getEntry() (string, error) {
	return v.guessEntry()
}

func (v *vmHandle) Exec(function string) error {
	entry, err := v.getEntry()
	if err != nil {
		return err
	}
	return v.vmInstance.Exec(entry)
}

func (v *vmHandle) ResourceUsed() gas.Limits {
	return v.vmInstance.ResourceUsed()
}

func (v *vmHandle) Release() {
	v.vmInstance.Release()
}

func (v *vmHandle) Abort(msg string) {
	v.vmInstance.Abort(msg)
}

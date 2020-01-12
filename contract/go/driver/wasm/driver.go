// +build wasm

package wasm

import (
	"github.com/BeDreamCoder/uwavm/contract/go/code"
	"github.com/BeDreamCoder/uwavm/contract/go/exec"
)

type driver struct {
}

// New returns a wasm driver
func New() code.Driver {
	return new(driver)
}

func (d *driver) Serve(contract code.Contract) {
	initDebugLog()
	exec.RunContract(0, contract, syscall)
}

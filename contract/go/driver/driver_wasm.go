// +build wasm

package driver

import (
	"github.com/BeDreamCoder/uwavm/contract/go/code"
	"github.com/BeDreamCoder/uwavm/contract/go/driver/wasm"
)

// Serve run contract in wasm environment
func Serve(contract code.Contract) {
	driver := wasm.New()
	driver.Serve(contract)
}

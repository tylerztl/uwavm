package vm

import (
	"errors"
	"sync"

	"github.com/BeDreamCoder/uwavm/wasm/exec"
)

type makeExecCodeFunc func(contractName string) (exec.WasmExec, error)

type ContractCode struct {
	ContractName string
	ExecCode     exec.WasmExec
}

type CodeManager struct {
	makeExecCode makeExecCodeFunc
	codes        map[string]*ContractCode
	mutex        sync.Mutex
}

func NewCodeManager(makeExec makeExecCodeFunc) *CodeManager {
	return &CodeManager{
		makeExecCode: makeExec,
		codes:        make(map[string]*ContractCode),
	}
}

func (c *CodeManager) GetExecCode(name string) (*ContractCode, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.codes[name]; ok {
		return nil, errors.New("old contract code not purged")
	}

	execCode, err := c.makeExecCode(name)
	if err != nil {
		return nil, err
	}
	code := &ContractCode{
		ContractName: name,
		ExecCode:     execCode,
	}
	c.codes[name] = code

	return code, nil
}

func (c *CodeManager) RemoveCode(name string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	code, ok := c.codes[name]
	if ok {
		code.ExecCode.Release()
	}
	delete(c.codes, name)
}

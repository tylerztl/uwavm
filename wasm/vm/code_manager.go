package vm

import (
	"errors"
	"sync"

	"github.com/BeDreamCoder/uwavm/exec"
)

type makeExecCodeFunc func(contractName string) (exec.Code, error)

type contractCode struct {
	ContractName string
	ExecCode     exec.Code
}

type codeManager struct {
	makeExecCode makeExecCodeFunc
	codes        map[string]*contractCode
	mutex        sync.Mutex
}

func newCodeManager(makeExec makeExecCodeFunc) *codeManager {
	return &codeManager{
		makeExecCode: makeExec,
		codes:        make(map[string]*contractCode),
	}
}

func (c *codeManager) GetExecCode(name string) (*contractCode, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.codes[name]; ok {
		return nil, errors.New("old contract code not purged")
	}

	execCode, err := c.makeExecCode(name)
	if err != nil {
		return nil, err
	}
	code := &contractCode{
		ContractName: name,
		ExecCode:     execCode,
	}
	c.codes[name] = code

	return code, nil
}

func (c *codeManager) RemoveCode(name string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	code, ok := c.codes[name]
	if ok {
		code.ExecCode.Release()
	}
	delete(c.codes, name)
}

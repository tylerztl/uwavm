package bridge

import (
	"fmt"

	"github.com/BeDreamCoder/uwavm/contract/go/pb"
)

// ContractError indicates the error of the contract running result
type ContractError struct {
	Status  int
	Message string
}

// Error implements error interface
func (c *ContractError) Error() string {
	return fmt.Sprintf("contract error status:%d message:%s", c.Status, c.Message)
}

// contractHandle 为vm.Context的实现，
// 它组合了合约内核态数据(cts)以及用户态的虚拟机数据(instance)
type contractHandle struct {
	cts      *ContractState
	instance Instance
	release  func()
}

func (c *contractHandle) Invoke(method string, args map[string][]byte) (*pb.Response, error) {
	c.cts.Method = method
	c.cts.Args = args
	err := c.instance.Exec("")
	if err != nil {
		return nil, err
	}
	if c.cts.Output == nil {
		return nil, &ContractError{
			Status:  500,
			Message: "internal error",
		}
	}

	return c.cts.Output, nil
}

func (c *contractHandle) ReleaseCache() error {
	// release the context of instance
	c.instance.Release()
	c.release()
	return nil
}

// vmImpl 为vm.VirtualMachine的实现
// 它是vmContextImpl的工厂类，根据不同的虚拟机类型(Executor)生成对应的vmContextImpl
type vmImpl struct {
	name  string
	state *StateManager
	exec  Executor
}

func (v *vmImpl) GetName() string {
	return v.name
}

func (v *vmImpl) NewVM(state *ContractState) (Contract, error) {
	ctx := v.state.CreateContractState()
	ctx.ContractName = state.ContractName
	ctx.Language = state.Language
	ctx.Caller = state.Caller

	release := func() {
		v.state.DestroyContractState(ctx)
	}

	instance, err := v.exec.NewCreatorInstance(ctx)
	if err != nil {
		v.state.DestroyContractState(ctx)
		return nil, err
	}
	return &contractHandle{
		cts:      ctx,
		instance: instance,
		release:  release,
	}, nil
}

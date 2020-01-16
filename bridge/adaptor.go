package bridge

import (
	"fmt"

	"github.com/BeDreamCoder/uwavm/contract/go/pb"
)

const InitMethod = "initialize"

// ContractError indicates the error of the contract running result
type ContractError struct {
	Status  int
	Message string
}

// Error implements error interface
func (c *ContractError) Error() string {
	return fmt.Sprintf("contract error status:%d message:%s", c.Status, c.Message)
}

// vmContextImpl 为vm.Context的实现，
// 它组合了合约内核态数据(ctx)以及用户态的虚拟机数据(instance)
type vmContextImpl struct {
	ctx      *ContractState
	instance Instance
	release  func()
}

func (v *vmContextImpl) Invoke(method string, args map[string][]byte) (*pb.Response, error) {
	v.ctx.Method = method
	v.ctx.Args = args
	err := v.instance.Exec()
	if err != nil {
		return nil, err
	}
	if v.ctx.Output == nil {
		return nil, &ContractError{
			Status:  500,
			Message: "internal error",
		}
	}

	return v.ctx.Output, nil
}

func (v *vmContextImpl) Release() error {
	// release the context of instance
	v.instance.Release()
	v.release()
	return nil
}

// vmImpl 为vm.VirtualMachine的实现
// 它是vmContextImpl的工厂类，根据不同的虚拟机类型(Executor)生成对应的vmContextImpl
type vmImpl struct {
	ctxmgr *StateManager
	name   string
	exec   Executor
}

func (v *vmImpl) GetName() string {
	return v.name
}

func (v *vmImpl) NewContext(state *ContractState) (Contract, error) {
	ctx := v.ctxmgr.CreateContractState()
	ctx.ContractName = state.ContractName
	ctx.Language = state.Language
	ctx.Caller = state.Caller

	release := func() {
		v.ctxmgr.DestroyContractState(ctx)
	}

	instance, err := v.exec.NewCreatorInstance(ctx)
	if err != nil {
		v.ctxmgr.DestroyContractState(ctx)
		return nil, err
	}
	return &vmContextImpl{
		ctx:      ctx,
		instance: instance,
		release:  release,
	}, nil
}

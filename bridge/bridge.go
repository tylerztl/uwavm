package bridge

import (
	"sync"

	"github.com/BeDreamCoder/uwavm/common/db"
	"github.com/BeDreamCoder/uwavm/contract/go/pb"
	"github.com/BeDreamCoder/uwavm/vm/gas"
)

// Executor 为用户态虚拟机工厂类
type Executor interface {
	// RegisterSyscallService 用于虚拟机把系统调用链接到合约代码上，类似vdso
	// 注册到Registry的时候被调用一次
	RegisterSyscallService(*SyscallService)
	// NewCreatorInstance 根据合约Context返回合约虚拟机的一个实例
	NewCreatorInstance(ctx *ContractState) (Instance, error)
	CallContract
}

// Instance is an instance of a contract run
type Instance interface {
	// Exec根据ctx里面的参数执行合约代码
	Exec(function string) error
	// ResourceUsed returns the resource used by contract
	ResourceUsed() gas.Limits
	// ReleaseCache releases contract instance
	Release()
	// Abort terminates running contract with error message
	Abort(msg string)
}

type Contract interface {
	Invoke(method string, args map[string][]byte) (*pb.Response, error)
	ResourceUsed() gas.Limits
	ReleaseCache() error
}

type CallContract interface {
	DeployContract(args map[string][]byte) (*pb.Response, gas.Limits, error)
	InvokeContract(method string, args map[string][]byte) (*pb.Response, gas.Limits, error)
}

// VirtualMachine define virtual machine interface
type VirtualMachine interface {
	GetName() string
	NewVM(state *ContractState) (Contract, error)
	CallContract
}

// Bridge 用于注册用户虚拟机以及向Xchain Core注册可被识别的vm.VirtualMachine
type Bridge struct {
	state   *StateManager
	syscall *SyscallService
	vms     map[string]VirtualMachine
}

var bridgeInstance *Bridge
var bridgeOnce sync.Once

// New instances a new Bridge
func GetBridge(db db.Database) *Bridge {
	bridgeOnce.Do(func() {
		state := NewStateManager()
		bridgeInstance = &Bridge{
			state:   state,
			syscall: NewSyscallService(state, db),
			vms:     make(map[string]VirtualMachine),
		}
	})
	return bridgeInstance
}

// RegisterExecutor register a Executor to Bridge
func (v *Bridge) RegisterExecutor(name string, exec Executor) VirtualMachine {
	wraper := &vmImpl{
		state: v.state,
		name:  name,
		exec:  exec,
	}
	exec.RegisterSyscallService(v.syscall)
	v.vms[name] = wraper
	return wraper
}

// GetVirtualMachine returns a contract.VirtualMachine from the given name
func (v *Bridge) GetVirtualMachine(name string) (VirtualMachine, bool) {
	vm, ok := v.vms[name]
	return vm, ok
}

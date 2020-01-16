package bridge

import (
	"sync"

	"github.com/BeDreamCoder/uwavm/contract/go/pb"
)

// ContractState 保存了合约执行的内核状态，
// 所有的系统调用产生的状态保存在这里
type ContractState struct {
	ID int64
	// 合约名字
	ContractName string

	Method string

	Args map[string][]byte

	Language string

	Initiator string

	Output *pb.Response
}

// StateManager 用于管理产生和销毁ContractState
type StateManager struct {
	// 保护如下两个变量
	// 合约进行系统调用以及合约执行会并发访问ctxs
	ctxlock sync.Mutex
	ctxid   int64
	ctxs    map[int64]*ContractState
}

// NewStateManager instances a new StateManager
func NewStateManager() *StateManager {
	return &StateManager{
		ctxs: make(map[int64]*ContractState),
	}
}

// ContractState 根据ContractState的id返回当前运行当前合约的上下文
func (n *StateManager) GetContractState(id int64) (*ContractState, bool) {
	n.ctxlock.Lock()
	defer n.ctxlock.Unlock()
	ctx, ok := n.ctxs[id]
	return ctx, ok
}

// CreateContractState allocates a ContractState with unique context id
func (n *StateManager) CreateContractState() *ContractState {
	n.ctxlock.Lock()
	defer n.ctxlock.Unlock()
	n.ctxid++
	ctx := new(ContractState)
	ctx.ID = n.ctxid
	n.ctxs[ctx.ID] = ctx
	return ctx
}

// DestroyContractState 一定要在合约执行完毕（成功或失败）进行销毁
func (n *StateManager) DestroyContractState(ctx *ContractState) {
	n.ctxlock.Lock()
	defer n.ctxlock.Unlock()
	delete(n.ctxs, ctx.ID)
}

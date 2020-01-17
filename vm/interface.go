package vm

import (
	"github.com/BeDreamCoder/uwavm/bridge"
	"github.com/BeDreamCoder/uwavm/common/db"
)

// NewInstanceCreatorFunc instances a new InstanceCreator from InstanceCreatorConfig
type NewInstanceCreatorFunc func(syscallService *bridge.SyscallService, db db.Database) (InstanceCreator, error)

// InstanceCreator is the creator of wasm virtual machine instance
type InstanceCreator interface {
	// CreateInstance instances a wasm virtual machine instance which can run a single contract call
	CreateInstance(ctx *bridge.ContractState) (bridge.Instance, error)
	RemoveCache(name string)
}

type CodeHandle interface {
	GetExecCode(name string) (*ContractCode, error)
	RemoveCode(name string)
}

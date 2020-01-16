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
	CreateInstance(ctx *bridge.ContractState) (Instance, error)
	RemoveCache(name string)
}

// Instance is a wasm virtual machine instance which can run a single contract call
type Instance interface {
	Exec(function string) error
	Release()
	Abort(msg string)
}

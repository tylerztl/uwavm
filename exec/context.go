package exec

import (
	"fmt"
)

const (
	// MaxGasLimit is the maximum gas limit
	MaxGasLimit = 0xFFFFFFFF
)

type ErrFuncNotFound struct {
	Name string
}

func (e *ErrFuncNotFound) Error() string {
	return fmt.Sprintf("%s not found", e.Name)
}

// ContextConfig configures an execution context
type ContextConfig struct {
	GasLimit int64
}

// DefaultContextConfig returns the default configuration of ContextConfig
func DefaultContextConfig() *ContextConfig {
	return &ContextConfig{
		GasLimit: MaxGasLimit,
	}
}

// GetContractState hold the context data when running a wasm instance
type Context interface {
	Exec(name string, param []int64) (ret int64, err error)
	GasUsed() int64
	ResetGasUsed()
	Memory() []byte
	StaticTop() uint32
	SetUserData(key string, value interface{})
	GetUserData(key string) interface{}
	Release()
}

type Code interface {
	NewContext(cfg *ContextConfig) (ictx Context, err error)
	Release()
}

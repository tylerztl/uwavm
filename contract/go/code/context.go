package code

import (
	"math/big"
)

// Context is the context in which the contract runs
type Context interface {
	Args() map[string][]byte
	Caller() string
	Initiator() string
	AuthRequire() []string

	PutObject(key []byte, value []byte) error
	GetObject(key []byte) ([]byte, error)
	DeleteObject(key []byte) error

	Transfer(to string, amount *big.Int) error
	Call(module, contract, method string, args map[string][]byte) (*Response, error)
}

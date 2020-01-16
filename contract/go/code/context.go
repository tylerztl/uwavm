package code

// GetContractState is the context in which the contract runs
type Context interface {
	Caller() string
	Args() map[string][]byte
	Method() string
	PutObject(key []byte, value []byte) error
	GetObject(key []byte) ([]byte, error)
	DeleteObject(key []byte) error
	Call(module, contract, method string, args map[string][]byte) (*Response, error)
}

package vm

import (
	"fmt"
	"github.com/BeDreamCoder/uwavm/bridge"
	"github.com/BeDreamCoder/uwavm/common/db"
)

var defaultRegistry = newRegistry()

type registry struct {
	drivers map[string]NewInstanceCreatorFunc
}

func newRegistry() *registry {
	return &registry{
		drivers: make(map[string]NewInstanceCreatorFunc),
	}
}

func (r *registry) Register(name string, driver NewInstanceCreatorFunc) {
	r.drivers[name] = driver
}

func (r *registry) Open(name string, syscallService *bridge.SyscallService, db db.Database) (InstanceCreator, error) {
	driverFunc, ok := r.drivers[name]
	if !ok {
		return nil, fmt.Errorf("driver %s not found", name)
	}
	return driverFunc(syscallService, db)
}

// Register makes a wasm driver available by the provided name
func Register(name string, driver NewInstanceCreatorFunc) {
	defaultRegistry.Register(name, driver)
}

// Open opens a wasm virtual machine specified by its driver name
func Open(name string, syscallService *bridge.SyscallService, db db.Database) (InstanceCreator, error) {
	return defaultRegistry.Open(name, syscallService, db)
}

package exec

import (
	"fmt"
)

var (
	// TrapOOB is raised when memory access out of bound
	TrapOOB = NewTrap("memory access out of bound")
	// TrapIntOverflow is raised when math overflow
	TrapIntOverflow = NewTrap("integer overflow on divide or truncation")
	// TrapDivByZero is raised when divide by zero
	TrapDivByZero = NewTrap("integer divide by zero")
	// TrapInvalidConvert is raised when convert from NaN to integer
	TrapInvalidConvert = NewTrap("conversion from NaN to integer")
	// TrapUnreachable is raised when unreachable instruction executed
	TrapUnreachable = NewTrap("unreachable instruction executed")
	// TrapInvalidIndirectCall is raised when run invalid call_indirect instruction
	TrapInvalidIndirectCall = NewTrap("invalid call_indirect")
	// TrapCallStackExhaustion is raised when call stack exhausted
	TrapCallStackExhaustion = NewTrap("call stack exhausted")
	// TrapGasExhaustion is raised when runnning out of gas limit
	TrapGasExhaustion = NewTrap("run out of gas limit")
	// TrapInvalidArgument is raised when running function with invalid argument
	TrapInvalidArgument = NewTrap("invalid function argument")
)

// Trap 用于表示虚拟机运行过程中的错误，中断虚拟机的运行
type Trap interface {
	Reason() string
}

// TrapError 用于包装一个Trap到Error
type TrapError struct {
	Trap Trap
}

func (t *TrapError) Error() string {
	return fmt.Sprintf("trap error:%s", t.Trap.Reason())
}

// Throw 用于抛出一个Trap
func Throw(trap Trap) {
	panic(trap)
}

// CaptureTrap 用于捕获潜在的Trap，如果是其他panic则不会捕获
func CaptureTrap(err *error) {
	ret := recover()
	if ret == nil {
		return
	}
	trap, ok := ret.(Trap)
	if ok {
		*err = &TrapError{
			Trap: trap,
		}
		return
	}
	panic(ret)
}

type stringTrap struct {
	reason string
}

func (s *stringTrap) Reason() string {
	return s.reason
}

// NewTrap returns a trap with the given reason
func NewTrap(reason string) Trap {
	return &stringTrap{
		reason,
	}
}

// TrapSymbolNotFound is raised when resolving symbol failed
type TrapSymbolNotFound struct {
	Module string
	Name   string
}

// Reason implements Trap interface
func (s *TrapSymbolNotFound) Reason() string {
	return fmt.Sprintf("%s.%s can't be resolved", s.Module, s.Name)
}

// TrapFuncSignatureNotMatch is raised when calling function signature is not matched
type TrapFuncSignatureNotMatch struct {
	Module string
	Name   string
}

// Reason implements Trap interface
func (s *TrapFuncSignatureNotMatch) Reason() string {
	return fmt.Sprintf("%s.%s not match with host signature", s.Module, s.Name)
}
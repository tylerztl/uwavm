package exec

import (
	"io"
)

const (
	debugWriterKey = "debugWriter"
)

// SetWriter set debug writer to GetContractState
func SetWriter(ctx Context, w io.Writer) {
	ctx.SetUserData(debugWriterKey, w)
}

// GetDebugWriter get debug writer
func GetWriter(ctx Context) io.Writer {
	value := ctx.GetUserData(debugWriterKey)
	if value == nil {
		return nil
	}
	w, ok := value.(io.Writer)
	if !ok {
		return nil
	}
	return w
}

// Write write debug message
// if SetWriter is not set, message will be ignored
func Write(ctx Context, p []byte) {
	w := GetWriter(ctx)
	if w == nil {
		return
	}
	w.Write(p)
}

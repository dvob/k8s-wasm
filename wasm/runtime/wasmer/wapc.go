package wasmer

import (
	"wasm/runner"
	"wasm/wapc"
)

func NewWAPCRuntime(source []byte) (runner.RawRunner, error) {
	return wapc.NewWasmerRuntime(source)
}

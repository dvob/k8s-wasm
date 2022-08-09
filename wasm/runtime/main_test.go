package wasm

import (
	"os"
	"testing"
	"wasm/runner"
	"wasm/wasmer"
	"wasm/wasmtime"
	"wasm/wazero"
)

const fnName = "run"

type InitRuntimeFunc func([]byte) (runner.RawRunner, error)

var InitRawRuntimes = map[string]InitRuntimeFunc{
	"wazero":   wazero.NewRawRuntime,
	"wasmer":   wasmer.NewRawRuntime,
	"wasmtime": wasmtime.NewRawRuntime,
}

var InitWASIRuntimes = map[string]InitRuntimeFunc{
	"wazero": wazero.NewWASIRuntime,
	// does not work with wasmer and wazero
	//"wasmer":   wasmer.NewWASIRuntime,
	//"wasmtime": wasmtime.NewWASIRuntime,
}

var InitWAPCRuntimes = map[string]InitRuntimeFunc{
	"wazero":   wazero.NewWAPCRuntime,
	"wasmer":   wasmer.NewWAPCRuntime,
	"wasmtime": wasmtime.NewWAPCRuntime,
}

func getRawRuntimes(t *testing.T, initFuncs map[string]InitRuntimeFunc, modulePath string) map[string]runner.RawRunner {
	module, err := os.ReadFile(modulePath)
	if err != nil {
		t.Fatal(err)
	}
	rawRunners := map[string]runner.RawRunner{}
	for name, initFn := range initFuncs {
		rawRunner, err := initFn(module)
		if err != nil {
			t.Fatalf("failed to initialize %s: %s", name, err)
		}
		rawRunners[name] = rawRunner
	}
	return rawRunners
}

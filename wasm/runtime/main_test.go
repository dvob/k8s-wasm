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

type ProtocolType string
type RuntimeType string

const (
	ProtocolTypeRaw  ProtocolType = "raw"
	ProtocolTypeWASI              = "wasi"
	ProtocolTypeWAPC              = "wapc"
)

const (
	RuntimeTypeWazero   RuntimeType = "wazero"
	RuntimeTypeWasmtime             = "wasmtime"
	RuntimeTypeWasmer               = "wasmer"
)

type runtime struct {
	protocol ProtocolType
	runtime  RuntimeType
	runner.RawRunner
}

type InitRuntimeFunc func([]byte) (runner.RawRunner, error)

var InitRawRuntimes = map[RuntimeType]InitRuntimeFunc{
	"wazero":   wazero.NewRawRuntime,
	"wasmer":   wasmer.NewRawRuntime,
	"wasmtime": wasmtime.NewRawRuntime,
}

var InitWASIRuntimes = map[RuntimeType]InitRuntimeFunc{
	"wazero": wazero.NewWASIRuntime,
	// does not work with wasmer and wazero
	//"wasmer":   wasmer.NewWASIRuntime,
	//"wasmtime": wasmtime.NewWASIRuntime,
}

var InitWAPCRuntimes = map[RuntimeType]InitRuntimeFunc{
	"wazero":   wazero.NewWAPCRuntime,
	"wasmer":   wasmer.NewWAPCRuntime,
	"wasmtime": wasmtime.NewWAPCRuntime,
}

func getRawRuntimes(t *testing.T, setup map[ProtocolType]string) []runtime {
	runtimes := []runtime{}
	for proto, modulePath := range setup {
		module, err := os.ReadFile(modulePath)
		if err != nil {
			t.Fatal(err)
		}
		var initFuncs map[RuntimeType]InitRuntimeFunc
		switch proto {
		case ProtocolTypeRaw:
			initFuncs = InitRawRuntimes
		case ProtocolTypeWASI:
			initFuncs = InitWASIRuntimes
		case ProtocolTypeWAPC:
			initFuncs = InitWAPCRuntimes
		}
		for runtimeType, initFn := range initFuncs {
			rawRunner, err := initFn(module)
			if err != nil {
				t.Fatalf("failed to initialize runtime=%s proto=%s: %s", runtimeType, proto, err)
			}
			runtimes = append(runtimes, runtime{
				protocol:  proto,
				runtime:   runtimeType,
				RawRunner: rawRunner,
			})
		}
	}
	return runtimes
}

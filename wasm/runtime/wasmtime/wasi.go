package wasmtime

import (
	"context"
	"fmt"
	"os"
	"wasm/runner"

	"github.com/bytecodealliance/wasmtime-go"
)

type WASIRuntime struct {
	engine *wasmtime.Engine
	module *wasmtime.Module
	linker *wasmtime.Linker
}

func NewWASIRuntime(moduleSource []byte) (runner.RawRunner, error) {

	e := &WASIRuntime{
		engine: wasmtime.NewEngine(),
	}
	e.linker = wasmtime.NewLinker(e.engine)
	err := e.linker.DefineWasi()
	if err != nil {
		return nil, err
	}

	e.module, err = wasmtime.NewModule(e.engine, moduleSource)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *WASIRuntime) MemorySize() uint64 {
	return 0
}

func (e *WASIRuntime) Run(_ context.Context, fnName string, input []byte) ([]byte, error) {
	// To implement this is really a pain because we can not use the Reader and Writer interface. See: https://github.com/bytecodealliance/wasmtime-go/issues/34
	// There seems to be another problem that if we use tempfiles Wasmtime does not close the files after use and we run into a "too many open files" error.
	// At least i could not see where I leak the open files. To say it is a bug in Wasmtime for sure we have to investigate further.

	panic("not implemented")

	store := wasmtime.NewStore(e.engine)
	wasiConfig := wasmtime.NewWasiConfig()

	// stdin
	stdin, err := os.CreateTemp("", "wasmtime-in-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(stdin.Name())
	_, err = stdin.Write(input)
	if err != nil {
		return nil, err
	}
	err = stdin.Close()
	if err != nil {
		return nil, err
	}

	// stdout
	stdout, err := os.CreateTemp("", "wasmtime-out-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(stdout.Name())
	err = stdout.Close()
	if err != nil {
		return nil, err
	}

	// stderr
	stderr, err := os.CreateTemp("", "wasmtime-err-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(stderr.Name())
	err = stderr.Close()
	if err != nil {
		return nil, err
	}

	wasiConfig.SetStdinFile(stdin.Name())
	wasiConfig.SetStdoutFile(stdout.Name())
	wasiConfig.SetStderrFile(stderr.Name())

	store.SetWasi(wasiConfig)

	instance, err := e.linker.Instantiate(store, e.module)
	if err != nil {
		return nil, err
	}

	fn := instance.GetFunc(store, fnName)
	if fn == nil {
		return nil, fmt.Errorf("missing function %s", fnName)
	}
	_, err = fn.Call(store)
	if err != nil {
		errout, _ := os.ReadFile(stderr.Name())
		return nil, fmt.Errorf("call failed. stderr='%s': %w", errout, err)
	}

	output, err := os.ReadFile(stdout.Name())
	if err != nil {
		return nil, err
	}

	return output, nil
}

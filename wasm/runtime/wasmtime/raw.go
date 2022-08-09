package wasmtime

import (
	"context"
	"fmt"
	"wasm/runner"

	"github.com/bytecodealliance/wasmtime-go"
)

type RawRuntime struct {
	store    *wasmtime.Store
	instance *wasmtime.Instance
	alloc    *wasmtime.Func
	dealloc  *wasmtime.Func
}

func NewRawRuntime(moduleSource []byte) (runner.RawRunner, error) {
	store := wasmtime.NewStore(wasmtime.NewEngine())
	module, err := wasmtime.NewModule(store.Engine, moduleSource)
	if err != nil {
		return nil, err
	}

	instance, err := wasmtime.NewInstance(store, module, []wasmtime.AsExtern{})
	if err != nil {
		return nil, err
	}

	engine := &RawRuntime{
		store:    store,
		instance: instance,
	}

	engine.alloc = engine.instance.GetFunc(store, "alloc")
	if engine.alloc == nil {
		return nil, fmt.Errorf("function alloc missing")
	}

	engine.dealloc = engine.instance.GetFunc(store, "dealloc")
	if engine.dealloc == nil {
		return nil, fmt.Errorf("function dealloc missing")
	}

	return engine, nil
}

func (e *RawRuntime) Run(_ context.Context, fnName string, input []byte) ([]byte, error) {

	runFn := e.instance.GetFunc(e.store, fnName)
	if runFn == nil {
		return nil, fmt.Errorf("function %s missing", fnName)
	}

	res, err := e.alloc.Call(e.store, len(input))
	if err != nil {
		return nil, err
	}

	inputPtr, ok := res.(int32)
	if !ok {
		return nil, fmt.Errorf("allocation did not return int32")
	}

	data := e.instance.GetExport(e.store, "memory").Memory().UnsafeData(e.store)

	n := copy(data[inputPtr:], input)
	if n != len(input) {
		return nil, fmt.Errorf("not all input bytes copied")
	}

	res, err = runFn.Call(e.store, inputPtr, len(input))
	if err != nil {
		return nil, err
	}

	outputPtrSize, ok := res.(int64)
	if !ok {
		return nil, fmt.Errorf("output did not return int64")
	}

	outputPtr := uint32(outputPtrSize >> 32)
	outputSize := uint32(outputPtrSize)

	data = e.instance.GetExport(e.store, "memory").Memory().UnsafeData(e.store)

	output := make([]byte, outputSize)

	copy(output, data[outputPtr:])

	_, err = e.dealloc.Call(e.store, outputPtr, outputSize)

	return output, nil
}

package wasmer

import (
	"context"
	"fmt"
	"wasm/runner"

	"github.com/wasmerio/wasmer-go/wasmer"
)

type RawRuntime struct {
	instance *wasmer.Instance
	alloc    wasmer.NativeFunction
	dealloc  wasmer.NativeFunction
}

func NewRawRuntime(moduleSource []byte) (runner.RawRunner, error) {
	store := wasmer.NewStore(wasmer.NewEngine())

	module, err := wasmer.NewModule(store, moduleSource)

	if err != nil {
		return nil, err
	}

	instance, err := wasmer.NewInstance(module, wasmer.NewImportObject())

	if err != nil {
		return nil, fmt.Errorf("failed to instantiate the module: %w", err)
	}

	engine := &RawRuntime{
		instance: instance,
	}

	engine.alloc, err = instance.Exports.GetFunction("alloc")
	if err != nil {
		return nil, fmt.Errorf("function alloc missing: %w", err)
	}

	engine.dealloc, err = instance.Exports.GetFunction("dealloc")
	if err != nil {
		return nil, fmt.Errorf("function dealloc missing: %w", err)
	}

	return engine, nil
}

func (e *RawRuntime) Run(_ context.Context, fnName string, input []byte) ([]byte, error) {
	runFn, err := e.instance.Exports.GetFunction(fnName)
	if err != nil {
		return nil, fmt.Errorf("could not get function '%s': %w", fnName, err)
	}
	if runFn == nil {
		return nil, fmt.Errorf("function %s missing", fnName)
	}

	res, err := e.alloc(len(input))
	if err != nil {
		return nil, err
	}

	inputPtr, ok := res.(int32)
	if !ok {
		return nil, fmt.Errorf("allocation did not return int32")
	}

	data, err := e.instance.Exports.GetMemory("memory")
	if err != nil {
		return nil, fmt.Errorf("failed to get memory")
	}

	n := copy(data.Data()[inputPtr:], input)
	if n != len(input) {
		return nil, fmt.Errorf("not all input bytes copied")
	}

	res, err = runFn(inputPtr, len(input))
	if err != nil {
		return nil, err
	}

	outputPtrSize, ok := res.(int64)
	if !ok {
		return nil, fmt.Errorf("output did not return int64")
	}

	outputPtr := uint32(outputPtrSize >> 32)
	outputSize := uint32(outputPtrSize)

	output := make([]byte, outputSize)

	copy(output, data.Data()[outputPtr:])

	_, err = e.dealloc(outputPtr, outputSize)

	return output, nil
}

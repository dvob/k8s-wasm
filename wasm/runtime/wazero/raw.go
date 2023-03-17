package wazero

import (
	"context"
	"fmt"
	"wasm/runner"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type RawRuntime struct {
	module  api.Module
	alloc   api.Function
	dealloc api.Function
}

func NewRawRuntime(moduleSource []byte) (runner.RawRunner, error) {
	ctx := context.Background()

	engine := &RawRuntime{}

	module, err := wazero.NewRuntime(ctx).Instantiate(ctx, moduleSource)
	if err != nil {
		return nil, err
	}
	engine.module = module

	engine.alloc = engine.module.ExportedFunction("alloc")
	if engine.alloc == nil {
		return nil, fmt.Errorf("function alloc missing")
	}

	engine.dealloc = engine.module.ExportedFunction("dealloc")
	if engine.dealloc == nil {
		return nil, fmt.Errorf("function dealloc missing")
	}

	return engine, nil
}

func (e *RawRuntime) Run(ctx context.Context, fnName string, input []byte) ([]byte, error) {
	runFn := e.module.ExportedFunction(fnName)
	if runFn == nil {
		return nil, fmt.Errorf("function '%s' not found", fnName)
	}
	inputPtr, err := e.alloc.Call(ctx, uint64(len(input)))
	if err != nil {
		return nil, fmt.Errorf("failed to call alloc function: %w", err)
	}

	mem := e.module.Memory()
	if !mem.Write(uint32(inputPtr[0]), input) {
		return nil, fmt.Errorf("Memory.Write(%d, %d) out of range of memory size %d", inputPtr[0], len(input), mem.Size())
	}

	resPtrSize, err := runFn.Call(ctx, inputPtr[0], uint64(len(input)))
	if err != nil {
		return nil, fmt.Errorf("failed to call run function: %w", err)
	}

	resultPtr := uint32(resPtrSize[0] >> 32)
	resultSize := uint32(resPtrSize[0])

	bytes, ok := mem.Read(resultPtr, resultSize)

	if !ok {
		return nil, fmt.Errorf("Memory.Read(%d, %d) out of range of memory size %d", resultPtr, resultSize, mem.Size())
	}

	_, err = e.dealloc.Call(ctx, uint64(resultPtr), uint64(resultSize))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

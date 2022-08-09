package wasmer

import (
	"context"
	"fmt"
	"wasm/runner"

	"github.com/wasmerio/wasmer-go/wasmer"
)

type WASIRuntime struct {
	store  *wasmer.Store
	module *wasmer.Module
}

func NewWASIRuntime(moduleSource []byte) (runner.RawRunner, error) {
	var err error
	e := &WASIRuntime{}
	e.store = wasmer.NewStore(wasmer.NewEngine())
	e.module, err = wasmer.NewModule(e.store, moduleSource)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *WASIRuntime) Run(_ context.Context, fnName string, input []byte) ([]byte, error) {
	// It seems that there is no way to pass a standard input
	// https://github.com/wasmerio/wasmer-go/issues/338
	panic("not implemented")
	wasiEnv, err := wasmer.NewWasiStateBuilder(e.module.Name()).
		CaptureStderr().
		CaptureStdout().
		Finalize()
	if err != nil {
		return nil, err
	}

	importObject, err := wasiEnv.GenerateImportObject(e.store, e.module)
	if err != nil {
		return nil, err
	}

	instance, err := wasmer.NewInstance(e.module, importObject)
	if err != nil {
		return nil, err
	}
	run, err := instance.Exports.GetFunction(fnName)
	if err != nil {
		return nil, err
	}

	if run == nil {
		return nil, fmt.Errorf("function %s not found", fnName)
	}

	_, err = run()
	if err != nil {
		stderr := wasiEnv.ReadStderr()
		return nil, fmt.Errorf("failed to run. stderr='%s': %w", stderr, err)
	}
	return wasiEnv.ReadStdout(), nil
}

package wapc

import (
	"context"
	"os"

	"github.com/wapc/wapc-go"
	"github.com/wapc/wapc-go/engines/wasmer"
	"github.com/wapc/wapc-go/engines/wasmtime"
	"github.com/wapc/wapc-go/engines/wazero"
)

type Runtime struct {
	module   wapc.Module
	instance wapc.Instance
}

func NewWasmerRuntime(code []byte) (*Runtime, error) {
	return NewRuntime(code, wasmer.Engine())
}

func NewWasmtimeRuntime(code []byte) (*Runtime, error) {
	return NewRuntime(code, wasmtime.Engine())
}

func NewWazeroRuntime(code []byte) (*Runtime, error) {
	return NewRuntime(code, wazero.Engine())
}

func NewRuntime(code []byte, engine wapc.Engine) (*Runtime, error) {
	ctx := context.Background()

	module, err := engine.New(ctx, wapc.NoOpHostCallHandler, code, &wapc.ModuleConfig{
		Logger: wapc.PrintlnLogger,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		return nil, err
	}
	// TODO: defer module.Close()

	instance, err := module.Instantiate(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: defer instance.Close(ctx)

	return &Runtime{
		module:   module,
		instance: instance,
	}, nil
}

func (e *Runtime) Run(ctx context.Context, fnName string, input []byte) ([]byte, error) {
	return e.instance.Invoke(ctx, fnName, input)
}

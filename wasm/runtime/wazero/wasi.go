package wazero

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"wasm/runner"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type WASIRuntime struct {
	mu       *sync.Mutex
	runtime  wazero.Runtime
	code     wazero.CompiledModule
	instance api.Module
	stdin    bytes.Buffer
	stdout   bytes.Buffer
	stderr   bytes.Buffer
}

func NewWASIRuntime(moduleSource []byte) (runner.RawRunner, error) {
	ctx := context.Background()

	runtime := wazero.NewRuntime(ctx)

	if _, err := wasi.Instantiate(ctx, runtime); err != nil {
		return nil, err
	}

	code, err := runtime.CompileModule(ctx, moduleSource)
	if err != nil {
		return nil, err
	}

	return &WASIRuntime{
		mu:      &sync.Mutex{},
		runtime: runtime,
		code:    code,
	}, nil
}

func (r *WASIRuntime) Run(ctx context.Context, fnName string, input []byte) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stdin.Reset()
	r.stdout.Reset()
	r.stderr.Reset()

	r.stdin.Write(input)

	config := wazero.NewModuleConfig().WithStdin(&r.stdin).WithStdout(&r.stdout).WithStderr(&r.stderr).WithStartFunctions()

	instance, err := r.runtime.InstantiateModule(ctx, r.code, config)
	if err != nil {
		return nil, fmt.Errorf("failed with stderr '%s': %w)", r.stderr.String(), err)
	}
	defer instance.Close(ctx)

	fn := instance.ExportedFunction(fnName)
	if fn == nil {
		return nil, fmt.Errorf("function '%s' missing", fnName)
	}

	_, err = fn.Call(ctx)
	if err != nil {
		errOut := r.stderr.String()
		if errOut != "" {
			return nil, fmt.Errorf("call to %s failed. stderr: '%s', err: %w", fnName, errOut, err)
		}
		return nil, err
	}

	output := make([]byte, r.stdout.Len())
	copy(output, r.stdout.Bytes())

	return output, nil
}

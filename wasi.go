package wasm

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	wasi "github.com/tetratelabs/wazero/wasi_snapshot_preview1"
)

type wasiExecutor struct {
	mu       *sync.Mutex
	runtime  wazero.Runtime
	code     wazero.CompiledModule
	instance api.Module
	stdin    bytes.Buffer
	stdout   bytes.Buffer
	stderr   bytes.Buffer
}

func newWasiExecutor(moduleSource []byte) (*wasiExecutor, error) {
	ctx := context.Background()

	runtime := wazero.NewRuntime()

	// Instantiate WASI, which implements system I/O such as console output.
	if _, err := wasi.Instantiate(ctx, runtime); err != nil {
		return nil, err
	}

	// Compile the WebAssembly module using the default configuration.
	code, err := runtime.CompileModule(ctx, moduleSource, wazero.NewCompileConfig())
	if err != nil {
		return nil, err
	}

	return &wasiExecutor{
		mu:      &sync.Mutex{},
		runtime: runtime,
		code:    code,
	}, nil
}

func (e *wasiExecutor) hasFunction(fnName string) bool {
	exportedFunctions := e.code.ExportedFunctions()
	_, ok := exportedFunctions[fnName]
	return ok
}

func (e *wasiExecutor) Close(ctx context.Context) error {
	return e.runtime.Close(ctx)
}

func (e *wasiExecutor) Run(ctx context.Context, fnName string, input []byte) ([]byte, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.stdin.Reset()
	e.stdout.Reset()
	e.stderr.Reset()

	e.stdin.Write(input)

	config := wazero.NewModuleConfig().WithStdin(&e.stdin).WithStdout(&e.stdout).WithStderr(&e.stderr).WithStartFunctions()

	instance, err := e.runtime.InstantiateModule(ctx, e.code, config)
	if err != nil {
		return nil, fmt.Errorf("failed with stderr '%s': %w)", e.stderr.String(), err)
	}
	defer instance.Close(ctx)

	fn := instance.ExportedFunction(fnName)
	if fn == nil {
		return nil, fmt.Errorf("function '%s' missing", fnName)
	}

	_, err = fn.Call(ctx)
	if err != nil {
		errOut := e.stderr.String()
		if errOut != "" {
			return nil, fmt.Errorf("call to %s failed. stderr: '%s', err: %w", fnName, errOut, err)
		}
		return nil, err
	}

	output := make([]byte, e.stdout.Len())
	copy(output, e.stdout.Bytes())

	return output, nil
}

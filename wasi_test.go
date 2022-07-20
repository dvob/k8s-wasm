package wasm

import (
	"context"
	"os"
	"strings"
	"testing"
)

const (
	wasiTestModuleFile = "testdata/test-wasi/target/wasm32-wasi/debug/test_wasi.wasm"
	// the function my_function writes the following outputs
	expectedStdout = "stdout output\n"
	expectedStderr = "stderr output\n"

	// output of my_panic function contains the following string
	panicOutput = "panic output"
)

func newTestWasiExecutor(t *testing.T) *wasiExecutor {
	source, err := os.ReadFile(wasiTestModuleFile)
	if err != nil {
		t.Fatal(err)
	}
	wasiExec, err := newWasiExecutor(source)
	if err != nil {
		t.Fatal(err)
	}

	return wasiExec
}

func TestWasiExecutorRun(t *testing.T) {
	wasiExec := newTestWasiExecutor(t)
	ctx := context.Background()

	outputRaw, err := wasiExec.Run(ctx, "my_function", []byte("stdin input"))
	output := string(outputRaw)
	if err != nil {
		t.Fatal(err)
	}
	if output != expectedStdout {
		t.Fatalf("want=%s, got=%s", expectedStdout, output)
	}
}

func TestWasiExecutorRunMultiple(t *testing.T) {
	wasiExec := newTestWasiExecutor(t)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		outputRaw, err := wasiExec.Run(ctx, "my_function", []byte("stdin input"))
		output := string(outputRaw)
		if err != nil {
			t.Fatal(err)
		}
		if output != expectedStdout {
			t.Fatalf("want=%s, got=%s", expectedStdout, output)
		}
	}
}

func TestWasiExecutorRunPanic(t *testing.T) {
	wasiExec := newTestWasiExecutor(t)
	ctx := context.Background()
	_, err := wasiExec.Run(ctx, "my_panic", []byte("stdin input"))
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), panicOutput) {
		t.Fatalf("panic output string '%s' not found in error output: '%s'", panicOutput, err.Error())
	}
}

func TestWasiExecutorRunError(t *testing.T) {
	wasiExec := newTestWasiExecutor(t)
	ctx := context.Background()
	_, err := wasiExec.Run(ctx, "my_error", []byte("stdin input"))
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), expectedStderr) {
		t.Fatalf("error output string '%s' not found in error output: '%s'", expectedStderr, err.Error())
	}
}

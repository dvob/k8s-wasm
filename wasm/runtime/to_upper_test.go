package wasm

import (
	"context"
	"testing"
)

func TestToUpperRaw(t *testing.T) {
	runtimes := getRawRuntimes(t, InitRawRuntimes, "../modules/rs/target/wasm32-unknown-unknown/debug/to_upper_raw.wasm")
	for name, runtime := range runtimes {
		t.Run(name, func(t *testing.T) {
			input := "foo"
			wantOutput := "FOO"

			ctx := context.Background()
			output, err := runtime.Run(ctx, fnName, []byte(input))

			if err != nil {
				t.Fatal(err)
			}

			if string(output) != wantOutput {
				t.Fatalf("want=%s, got=%s", wantOutput, string(output))
			}
		})
	}
}

func TestToUpperWASI(t *testing.T) {
	runtimes := getRawRuntimes(t, InitWASIRuntimes, "../modules/rs/target/wasm32-wasi/debug/to_upper_wasi.wasm")
	for name, runtime := range runtimes {
		t.Run(name, func(t *testing.T) {
			input := "foo"
			wantOutput := "FOO"

			ctx := context.Background()
			output, err := runtime.Run(ctx, fnName, []byte(input))

			if err != nil {
				t.Fatal(err)
			}

			if string(output) != wantOutput {
				t.Fatalf("want=%s, got=%s", wantOutput, string(output))
			}
		})
	}
}

func TestToUpperWAPC(t *testing.T) {
	runtimes := getRawRuntimes(t, InitWAPCRuntimes, "../modules/rs/target/wasm32-unknown-unknown/debug/to_upper_wapc.wasm")
	for name, runtime := range runtimes {
		t.Run(name, func(t *testing.T) {
			input := "foo"
			wantOutput := "FOO"

			ctx := context.Background()
			output, err := runtime.Run(ctx, fnName, []byte(input))

			if err != nil {
				t.Fatal(err)
			}

			if string(output) != wantOutput {
				t.Fatalf("want=%s, got=%s", wantOutput, string(output))
			}
		})
	}
}

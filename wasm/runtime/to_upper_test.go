package wasm

import (
	"context"
	"testing"
)

func TestToUpper(t *testing.T) {
	runtimes := getRawRuntimes(t, map[ProtocolType]string{
		ProtocolTypeRaw:  "../modules/rs/target/wasm32-unknown-unknown/debug/to_upper_raw.wasm",
		ProtocolTypeWASI: "../modules/rs/target/wasm32-wasi/debug/to_upper_wasi.wasm",
		ProtocolTypeWAPC: "../modules/rs/target/wasm32-unknown-unknown/debug/to_upper_wapc.wasm",
	})
	for _, runtime := range runtimes {
		name := string(runtime.runtime) + "_" + string(runtime.protocol)
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

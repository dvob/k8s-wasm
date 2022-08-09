package wasm

import (
	"context"
	"encoding/json"
	"testing"

	v1 "k8s.io/api/authentication/v1"
)

func TestTokenReview(t *testing.T) {
	runtimes := getRawRuntimes(t, map[ProtocolType]string{
		ProtocolTypeRaw:  "../modules/rs/target/wasm32-unknown-unknown/debug/k8s_raw.wasm",
		ProtocolTypeWASI: "../modules/rs/target/wasm32-wasi/debug/k8s_wasi.wasm",
		ProtocolTypeWAPC: "../modules/rs/target/wasm32-unknown-unknown/debug/k8s_wapc.wasm",
	})
	for _, runtime := range runtimes {
		name := string(runtime.runtime) + "_" + string(runtime.protocol)
		t.Run(name, func(t *testing.T) {
			t.Run("correct", func(t *testing.T) {
				tr := &v1.TokenReview{
					Spec: v1.TokenReviewSpec{
						Token: "correct-token",
					},
				}

				input, err := json.Marshal(tr)
				if err != nil {
					t.Fatal(err)
				}

				ctx := context.Background()
				output, err := runtime.Run(ctx, fnName, input)
				if err != nil {
					t.Fatal(err)
				}

				response := &v1.TokenReview{}
				err = json.Unmarshal(output, response)
				if err != nil {
					t.Fatal(err)
				}

				if !response.Status.Authenticated {
					t.Fatalf("resonse should be authenticated: '%s'", output)
				}
			})
			t.Run("wrong", func(t *testing.T) {
				tr := &v1.TokenReview{
					Spec: v1.TokenReviewSpec{
						Token: "wrong-token",
					},
				}

				input, err := json.Marshal(tr)
				if err != nil {
					t.Fatal(err)
				}

				ctx := context.Background()
				output, err := runtime.Run(ctx, fnName, input)
				if err != nil {
					t.Fatal(err)
				}

				response := &v1.TokenReview{}
				err = json.Unmarshal(output, response)
				if err != nil {
					t.Fatal(err)
				}

				if response.Status.Authenticated {
					t.Fatalf("resonse should not be authenticated: '%s'", output)
				}
			})
		})
	}
}

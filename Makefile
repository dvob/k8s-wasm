
build:
	cd testdata/test-wasi/ && cargo build --target wasm32-wasi
	cd testdata/test-authn/ && cargo build --target wasm32-wasi

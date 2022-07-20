
testdata/test-wasi/target/wasm32-wasi/debug/test_wasi.wasm: testdata/test-wasi/src/*
	cd testdata/test-wasi/ && cargo build --target wasm32-wasi

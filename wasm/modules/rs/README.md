# Rust Example Modules
This directory contains example WASM modules implemented with Rust.
These examples were used to experiment how we can pass data between runtime and module.

## Build
We can't use a Cargo workspace here because we want to build the individual crates for different targets.
Some for `wasm32-unknown-unknown` and some for ` wasm32-wasi`.

Run `./build.sh` to build all modules with the correct target.
```
./build.sh
```

The target directory for all modules is set to `./target` (see `.cargo/config`), so you'll find the compiled wasm modules in this directory under `./target/`.

## Modules
The `-raw` modules use the method to pass data between host and module described in a [blog](https://radu-matei.com/blog/practical-guide-to-wasm-memory/) post of Radu Matei.
The `-wasi` modules are compiled with the target `wasm32-wasi` and use the standard input and standard output to pass data between host and module.
The `k8s-` modules show how to process a Kubernetes [TokenReview](https://pkg.go.dev/k8s.io/api/authentication/v1#TokenReview) object.
For this we use the [k8s-openapi](https://crates.io/crates/k8s-openapi) and [serde_json](https://crates.io/crates/serde_json) crates.

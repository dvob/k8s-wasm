# WASM

* [Simple example (use C code from JavaScript)](./simple-example/)
* [Runtimes](./runtime/)
* [Modules](./modules/rs/)

This directory contains code to explore how we can use WebAssembly modules from Go and how we can pass data between the host and the module.

To extend the Kubernetes Authentication, Authorization and Admission process we have to pass the relevant data to the WASM module and reveive the result from the module.
For this we use data in the same format as it is used for the webhooks (TokenReview, SubjectAccessReview, AdmissionReview).

First we have to verify that we can pass raw bytes (`[]byte`) back and forth.
For this we implement a module which takes a string and transforms it to upper case.
In a second step we send a JSON encoded `TokenReview` back and forth.

We implement these basic steps which we later use in our Kubernetes extension with different data passing mechanisms:
* `raw`: Run alloc, dealloc manually and pass pointers.
* `wasi`: Read and write from and to standard input and standard output
* `wapc`: Use [waPC specification](https://wapc.io/docs/spec/)

# Links
* Core Spec: https://webassembly.github.io/spec/core/
* Additional Features: https://webassembly.org/roadmap/
* Memory Interactions
  * https://radu-matei.com/blog/practical-guide-to-wasm-memory/
  * https://blog.nishtahir.com/interacting-with-wasm-memory/
  * https://rob-blackbourn.github.io/blog/webassembly/wasm/array/arrays/javascript/c/2020/06/07/wasm-arrays.html
* Memory Instances: https://webassembly.github.io/spec/core/exec/runtime.html#memory-instances
* Terms: https://webassembly.github.io/spec/core/intro/overview.html
* Component Model / Interface Types (types):  https://github.com/WebAssembly/component-model/blob/main/design/mvp/Explainer.md#type-definitions
* Runtimes:
  * Wasmer: https://wasmer.io/
  * Wasmtime (Bytecode allicance): https://wasmtime.dev/
  * Wazero (pure Go): https://github.com/tetratelabs/wazero
* Data passing:
  * WaPC: https://wapc.io/docs/spec/
  * WAGI: https://github.com/deislabs/wagi
  * WIT: https://radu-matei.com/blog/wasm-components-host-implementations/

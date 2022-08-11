# Extend Kubernetes with WASM

This repository explains an extension of the Kubernetes API-Server which allows to use WebAssembly modules to perform the following actions:
* Authentication
* Authorization
* Admission (validating and mutating)

For this I forked the Kubernetes project and extended the API-Server accordingly: https://github.com/dvob/kubernetes/tree/wasm.
The extension is based on the [release-1.24](https://github.com/kubernetes/kubernetes/tree/release-1.24) branch.

To code is **not intended for production use**. Its a proof of concept to show how WebAssembly could be used to extend Kubernetes.

In the fork I created a new package `pkg/wasm` under which I implemented a new Authenticator, Authorizer and AdmissionController.
Most of the implementation lives in this new package.
There are only a few changes in the `pkg/kube-apiserver` to add command line options which allow to enable the WASM Authenticator, Authorizer and AdmissionController.

To run the WebAssembly modules we use the [Wazero](https://github.com/tetratelabs/wazero) runtime.
Wazero has zero dependencies and does not rely on CGO. Hence it can be easy integrated in a Go project without adding a ton of dependencies.

To pass data between our extension (host) and the WASM module we expect that the modules target [WASI](https://wasi.dev/).
The host then writes the request JSON encoded to the standard input (stdin) of the module.
The module then can write the response to the standard output.
See [Module Specification](./spec/) for more details.

For Admission the extension also supports to use [Kubewarden policies](https://hub.kubewarden.io/) which are not context aware.

See [User Documentation](./docs/main/) for more details on how to setup and configure the extended API-Server.

## Links Overview
* [User Documentation](./docs/main): How to setup and configure the WASM extension
* [Kubernetes Cluster Setup](./docs/cluseter_setup/): How to setup a Kubernetes cluster
* [Module Specification](./spec/): 
  * [Rust module library](https://github.com/dvob/k8s-wasi-rs): Rust library which simplifies the creation of WASM modules according to the Module Specification.
* Experiments (PoCs)
  * [WASM](./wasm/)
    * [Runtimes](./wasm/runtime): Comparison of WebAssembly Runtimes
    * Data passing ([Modules](./wasm/modules/rs), [Runtimes](./wasm/runtimes/)
  * [Kubernetes Integration](./k8s/)
    * [Webhook](./k8s/webhook/)
    * [Direct in API-Server](./k8s/api-server/)

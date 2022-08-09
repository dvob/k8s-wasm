# Extend Kubernetes with WASM

This repository explains an extension of the Kubernetes API-Server which allows to configure WebAssembly modules to perform the following actions:
* Authentication
* Authorization
* Admission (validating and mutating)

For this I [forked](https://github.com/dvob/kubernetes/tree/wasm) the Kubernetes project and extended the API-Server accordingly.
The changes are done based on the [realease-1.24](https://github.com/kubernetes/kubernetes/tree/release-1.24) branch.
In the Kubernetes project I created a new package `pkg/wasm` under which I implemented a new Authenticator, Authorizer and AdmissionController.
Most of the implementation lives in this new package.
There are only a few changes in the `pkg/kube-apiserver` to add options which allow to enable the WASM Authenticator, Authorizer and AdmissionController.
This should make it easy to merge/rebase the additions onto new/other Kubernetes versions.

For more information check the documentation in this repository:

* [Setup documentation](./docs/setup.md): How to setup a Kubernetes cluster with the extended API-Server.
* [User documentatin](./docs/user.md): How to configure the WASM extensions in the API-Server (Authenticator, Authroizer, AdmissionController).
* [Module specification](./spec/README.md): How to implement your own modules.
  * In addition to this specification the extension supports to run [Kubewarden policies](https://hub.kubewarden.io/) which are not context aware.
* [Rust module library](https://github.com/dvob/k8s-wasi-rs): A rust library which simplifies the creation of WASM modules.

## PoCs and Expreiments
* [WASM](./wasm/README.md)
  * Runtime
  * Data passing
* [Kubernetes](./k8s/README.md)
  * Integration

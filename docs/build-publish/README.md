# Build and test Kubernetes API server

This page explains how to build and publish the API server which contains the WASM extension.

## Dependencies
To run the tests and build the API server you need the following tools:

* System tools
  * Git
  * Curl
  * Make
  * C compiler (`cc`)
  * Rsync

```
apt-get install git curl gcc make rsync
```

* Docker (https://docs.docker.com/engine/install/)
* Go (https://go.dev/doc/install)
* Rustup (Rust, Cargo, etc.)o (https://www.rust-lang.org/tools/install)

Install wasm32-wasi target:
```
rustup target add wasm32-wasi
```

## Run tests
Clone the extended API server code:
```bash
git clone -b wasm git@github.com:dvob/kubernetes.git
cd kubernetes/
```

Before you can run the tests for the `pkg/wasm` package you have to run the prepare script which builds the test modules and downloads some Kubewarden policies (modules).
```
cd pkg/wasm/testmodules/
./prepare_integration_tests.sh
```

Then you can run all WASM related tests in the `./pkg/wasm` directory:
```
cd $( git rev-parse --show-toplevel )/pkg/wasm
go test -v ./...
```

## Build & Publish
From the root of the repository you can run the build of the whole Kubernetes project:
```
cd $( git rev-parse --show-toplevel )
make quick-release-images KUBE_BUILD_PLATFORMS=linux/amd64
```

This publishes the build artifacts to the `_output` directory. For the API server you can find the following artifacts:
* Binary: `:/_output/release-stage/server/linux-amd64/kubernetes/server/bin/kube-apiserver`
* Docker image (TAR): `./_output/release-images/amd64/kube-apiserver.tar`

Depending on your cluster setup you either need the `kube-apiserver` binary or the Docker image.
In the case of a `kubeadm` Kubernetes installation you need the Docker image.

To use the Docker image you have to publish it to a Docker registry. In the following example we publish the image to `dvob/kube-apiserver:dev`.
```
docker login
docker load -i ./_output/release-images/amd64/kube-apiserver.tar
docker tag <name of loaded image> dvob/kube-apiserver:dev
docker push dvob/kube-apiserver:dev
```

To publish the TAR directly you can also use the tool [crane](https://github.com/google/go-containerregistry/tree/main/cmd/crane):
```
crane push _output/release-images/amd64/kube-apiserver.tar dvob/kube-apiserver:magic-example
```

### Run
See [Setup Kubernetes cluster](../cluster-setup/) for instructions on how to setup a Kubernetes cluster and configure a custom API server image.

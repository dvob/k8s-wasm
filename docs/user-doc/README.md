# Setup and Configure WASM Extension

The following page explains how to setup and configure the WASM extension for the Kubernetes API server.

## Setup
First you have to build the API server which contains the WASM extension. See [Build and test Kubernetes API server](../build-publish/).

Then you have to setup a Kubernetes cluster with the build of the Kubernetes API server which contains the WASM extension.
See [cluster setup documentation](../cluster-setup/) on how to setup and run a Kubernetes cluster with a custom build of an API server.

If you don't want to build the API server by your own you can use the following image:
* `dvob/kube-apiserver:wasm` (`dvob/kube-apiserver@sha256:69f9bc68e50bffb0db5ed105ee10b8098adc5a029449ad91543eb97e37440f15`)

## Mount configuration
To configure the WASM extension you have to prepare the configuration files and the actual WASM modules.
Copy the files to the server which runs your API-Server.
In a typical kubeadm setup you also have to update the kube-apiserver mainfest to mount the files from your server into the API server Pod.

The easiest way to mount all required files into the API server Pod is to place all files in one directory and mount that directory into the API server.
For this you have to extend `/etc/kubernetes/manifests/kube-apiserver.yaml` with the following parts:
```yaml
# spec.containers[0].volumeMounts
    - mountPath: /etc/kubernetes/wasm
      name: wasm
      readOnly: true

# spec.volumes
  - hostPath:
      path: /etc/kubernetes/wasm
      type: DirectoryOrCreate
    name: wasm
```

## Module Basic Configuration
For all use cases (Authentication, Authorization and Admission) you configure a list of WASM modules which should be consulted on a request:
```yaml
modules:
- module: /etc/kubernetes/wasm/my_module1.wasm
- module: /etc/kubernetes/wasm/my_module2.wasm
```

All use cases use the same basic module configuration.
```yaml
# modules specifies the modules which should be consulted on a request
modules:
- # name (optional) specifies the name to identify the wasm module (e.g. in log messages).
  # if not specified it defaults to the basename of the module path (in this example 'file.wasm')
  name: my-module-1

  # module is the path to the wasm file
  module: /path/to/the/file.wasm

  # settings (optional) for the module. these are passed with each invocation
  # (see module specification).
  settings: {}

  # debug (optional) enables debug output. this prints all inputs and outputs
  # which are passed between the host and the module.
  debug: false
```

The admission modules support additional configurations (see below).

If you update the configuration or the modules you have to restart the API server that the changes become active.

# Authentication
To enable the WASM authentication you have to configure the following option on the API server:
```
--authentication-wasm-config-file=/etc/kubernetes/wasm/authn.conf
```

The authentication extension consults each module in the module list until one successfully authenticates the token.

## Authentication Example
`/etc/kubernetes/wasm/authn.conf`:
```yaml
modules:
- name: magic-authentication
  module: /etc/kubernetes/wasm/magic_authenticator.wasm
```

Copy the module file from https://github.com/dvob/k8s-wasi-rs/releases/download/v0.1.1/magic_authenticator.wasm to `/etc/kubernetes/wasm/magic_authenticator.wasm`.
The magic authenticator authenticates requests which use the token `magic-token`.

# Authorization
To enable the WASM authorization you have to add `WASM` to the authorization modes and specify a modules configuration:
```
--authorization-mode=Node,RBAC,WASM
--authorization-wasm-config-file=/etc/kubernetes/wasm/authz.conf
```

The authorization extension consults each module in the module list until one successfully authorizes the request.

## Authorization Example
`/etc/kubernetes/wasm/authz.conf`:
```yaml
modules:
- name: magic-authorization
  module: /etc/kubernetes/wasm/magic_authorizer.wasm
```

Copy the module file from https://github.com/dvob/k8s-wasi-rs/releases/download/v0.1.1/magic_authorizer.wasm to `/etc/kubernetes/wasm/magic_authenticator.wasm`.
The magic authorizer authorizes your request if you are a member of the `magic-group` and you try to access a config map.

# Admission
To enable the WASM admission you have to add the `WASM` admission controller to the list of enabled admission plugins `--enable-admission-plugins`.
To configure the WASM admission controller you have to provide the configuration with the admission control config file `--admission-control-config-file`.
```
--enable-admission-plugins=WASM
--admission-control-config-file=/etc/kubernetes/wasm/admission.conf
```

The admission configuration supports the same basic module configuration as described above. The following additional settings are supported for admission modules:
```yaml
modules:
- module: /path/to/module.wasm

  # mutating (optional) specifies if a module is mutating or not. If mutating is
  # set to true the module gets called during the mutating phase of the admission.
  # By default mutating is false.
  mutating: false

  # type (optional) specifies the module type. Valid types are 'wasi' and 'kubewarden'.
  type: wasi

  # rules is a list of rules which define for which resources the admission
  # module should be called. See https://pkg.go.dev/k8s.io/api/admissionregistration/v1#RuleWithOperations
  # for documentation or look at the examples. if no rules are specified your module never gets called.
  rules:
  - operations: ["CREATE", "UPDATE"]
    apiGroups: [""]
    apiVersions: ["v1"]
    resources: ["configmaps"]
```

If you specify the type `wasi` the module has to conform to the [module specification](../../spec/).
If `kubewarden` is used as type the modules are called as described in the [Kubewarden policy specification](https://docs.kubewarden.io/writing-policies/spec/intro-spec).
You can find Kubewarden modules here: https://hub.kubewarden.io/

The WASM admission configuration is part of the full admission configuration and is either included as separate file or directly in the admission configuration.

File:
```yaml
apiVersion: apiserver.config.k8s.io/v1
kind: AdmissionConfiguration
plugins:
- name: WASM
  path: /path/to/admission-module-configuration.conf
```

Direct:
```yaml
apiVersion: apiserver.config.k8s.io/v1
kind: AdmissionConfiguration
plugins:
- name: WASM
  configuration:
    modules:
    - name: magic-validation
      module: /etc/kubernetes/wasm/magic_validator.wasm
      rules:
      - operations: ["CREATE", "UPDATE"]
        apiGroups: [""]
        apiVersions: ["v1"]
        resources: ["configmaps"]
    - name: magic-mutator
      module: /etc/kubernetes/wasm/magic_mutator.wasm
      mutating: true
      rules:
      - operations: ["CREATE", "UPDATE"]
        apiGroups: [""]
        apiVersions: ["v1"]
        resources: ["configmaps"]
```

## Admission Examples

### Example with Magic-Modules
`/etc/kubernetes/wasm/admission.conf`:
```yaml
apiVersion: apiserver.config.k8s.io/v1
kind: AdmissionConfiguration
plugins:
- name: WASM
  configuration:
    modules:
    - name: magic-validation
      module: /etc/kubernetes/wasm/magic_validator.wasm
      rules:
      - operations: ["CREATE", "UPDATE"]
        apiGroups: [""]
        apiVersions: ["v1"]
        resources: ["configmaps"]
    - name: magic-mutator
      module: /etc/kubernetes/wasm/magic_mutator.wasm
      mutating: true
      rules:
      - operations: ["CREATE", "UPDATE"]
        apiGroups: [""]
        apiVersions: ["v1"]
        resources: ["configmaps"]
```

Copy the following module files to `/etc/kubernetes/wasm/`
* https://github.com/dvob/k8s-wasi-rs/releases/download/v0.1.1/magic_validator.wasm -> `/etc/kubernetes/wasm/magic_validator.wasm`
* https://github.com/dvob/k8s-wasi-rs/releases/download/v0.1.1/magic_mutator.wasm -> `/etc/kubernetes/wasm/magic_mutator.wasm`

The magic mutator adds the value `magic-value=foobar` to every config map.
The magic validator rejects the creation or update of config maps which contain the value `not-allowed-value`.

### Example with Kubewarden Policy
In the following example we ensure that `configmaps` and `namespaces` have an annotation `puzzle.ch/owner`.
For this we use the Kubewarden policy [safe-annotations](https://github.com/kubewarden/safe-annotations-policy) and set the appropriate settings (`mandatory_annotations`).

`/etc/kubernetes/wasm/admission.conf`:
```yaml
apiVersion: apiserver.config.k8s.io/v1
kind: AdmissionConfiguration
plugins:
- name: WASM
  configuration:
    modules:
    - name: safe-annotations
      type: kubewarden
      module: /etc/kubernetes/wasm/safe-annotations.wasm
      settings:
        mandatory_annotations:
        - puzzle.ch/owner
      rules:
      - operations: ["CREATE", "UPDATE"]
        apiGroups: [""]
        apiVersions: ["v1"]
        resources:
        - configmaps
        - namespaces
```

Copy https://github.com/kubewarden/safe-annotations-policy/releases/download/v0.2.0/policy.wasm to `/etc/kubernetes/wasm/safe-annotations.wasm`

# Kubernetes API server Extension

To explore the extension of the API server with an Authentication, an Authorizer and an AdmissionController we implemented the simple example logic described [here](../).
You can find the code here:
* Branch: https://github.com/dvob/kubernetes/tree/magic-examples
* Diff (commit): https://github.com/dvob/kubernetes/commit/0919013fca6558437beb4afc7c2aeaeba66d2683

## Implementation
To implement theses components we have to implement the appropriate interfaces and then integrate the implementation in the actual API server.

### Authentication
* [Interface](https://github.com/kubernetes/kubernetes/blob/0425c85cfc612cecdc4a333f5025163afec06615/staging/src/k8s.io/apiserver/pkg/authentication/authenticator/interfaces.go#L28)
* [Implementation](https://github.com/dvob/kubernetes/blob/0919013fca6558437beb4afc7c2aeaeba66d2683/staging/src/k8s.io/apiserver/plugin/pkg/authenticator/token/magic/authenticator.go)
* [Integration](https://github.com/dvob/kubernetes/blob/0919013fca6558437beb4afc7c2aeaeba66d2683/pkg/kubeapiserver/authenticator/config.go#L192)

### Authorization
* [Interface](https://github.com/kubernetes/kubernetes/blob/0425c85cfc612cecdc4a333f5025163afec06615/staging/src/k8s.io/apiserver/pkg/authorization/authorizer/interfaces.go#L70)
* [Implementation](https://github.com/dvob/kubernetes/blob/0919013fca6558437beb4afc7c2aeaeba66d2683/staging/src/k8s.io/apiserver/plugin/pkg/authorizer/magic/authorizer.go)
* [Integration](https://github.com/dvob/kubernetes/blob/0919013fca6558437beb4afc7c2aeaeba66d2683/pkg/kubeapiserver/authorizer/config.go#L144)

### Admission
* Interface
  * [Mutation](https://github.com/kubernetes/kubernetes/blob/0425c85cfc612cecdc4a333f5025163afec06615/staging/src/k8s.io/apiserver/pkg/admission/interfaces.go#L129)
  * [Validation](https://github.com/kubernetes/kubernetes/blob/0425c85cfc612cecdc4a333f5025163afec06615/staging/src/k8s.io/apiserver/pkg/admission/interfaces.go#L138)
* [Implementation](https://github.com/dvob/kubernetes/blob/0919013fca6558437beb4afc7c2aeaeba66d2683/plugin/pkg/admission/magic/admission.go)
* [Integration](https://github.com/dvob/kubernetes/blob/0919013fca6558437beb4afc7c2aeaeba66d2683/pkg/kubeapiserver/options/plugins.go#L144)

## Build & Run

Clone the fork:
```bash
git clone -b magic-examples git@github.com:dvob/kubernetes.git
cd kubernetes/
```

Build Kubernetes (see the official Kubernetes [build documentation](https://github.com/kubernetes/community/blob/master/contributors/devel/development.md#building-kubernetes) for more information):
```bash
make quick-release-images KUBE_BUILD_PLATFORMS=linux/amd64
```

This publishes the build artifacts to the `_output` directory. For the API server you can find the following artifacts:
* Binary: `:/_output/release-stage/server/linux-amd64/kubernetes/server/bin/kube-apiserver`
* Docker image (TAR): `./_output/release-images/amd64/kube-apiserver.tar`

Depending on your setup you either have to copy the `kube-apiserver` binary to the server or publish the `kube-apiserver` Docker image somewhere to use it as your new API server.

With the tool [crane](https://github.com/google/go-containerregistry/tree/main/cmd/crane) you can publish a Docker image in the TAR format to a registry easily:
```
crane push _output/release-images/amd64/kube-apiserver.tar dvob/kube-apiserver:magic-example
```

Now you have to run the API server with the following options to enable our own MagicAuthenticator, MagicAuthorizer, and MagicAdmissionController:
```
--magic-auth
--authorization-mode=Node,RBAC,Magic
--enable-admission-plugins=NodeRestriction,MagicAdmission
```

# Links
* Kubernetes Development
  * Logging: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-instrumentation/logging.md
  * Use dependencies: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/vendor.md

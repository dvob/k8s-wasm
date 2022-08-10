# Kubernetes API-Server Webhook Integration

The webhook server in this directory implements the following webhooks:
* Authentication `/authn`
* Authorization `/authz`
* Validating Admission `/admit`
* Mutating Admission `/admit-mut`

For each of the action a minimal example got implemented.

The webhooks perform the following actions:
* `/authn` handels TokenReview requests for authentication. If the token equals to `magic-token` the user is authenticated to user `magic-user` with the group `magic-group`.
* `/authz` handles SubjectAccessReview requests for authorization. It allows members of the group `magic-group` to manage `configmaps`.
* `/admit` handles AdmissionReview requests and performs validation on configmaps. If the configmap contains a value `not-allowed-value` it rejects the request.
* `/admit-mut` handles AdmissionReview requests and performs mutation. It adds the value `magic-value=foobar` to every configmap.

## Setup Test Cluster
The following section General explains the test setup for the webhooks in general. For a quick setup with kind see section "Setup with kind" below.

### General
You can test the webhook server with any Kubernetes cluster.
To configure the authentication and authorization webhooks you have to set the following options on the API-Server:
```
--authorization-mode=Node,RBAC,Webhook
--authorization-webhook-config-file=/etc/kubernetes/authz-webhook.conf
--authentication-token-webhook-config-file=/etc/kubernetes/authn-webhook.conf
```

In a typical Kubernetes installation where your API-Server is started by the kubelet you have to add the appropriate options to the manifest of the kube-apiserver (`/etc/kubernetes/manifests/kube-apiserver.yaml`):
```
# api server options
    - --authorization-mode=Node,RBAC,Webhook
    - --authorization-webhook-config-file=/etc/kubernetes/authz-webhook.conf
    - --authentication-token-webhook-config-file=/etc/kubernetes/authn-webhook.conf
# mounts
    - mountPath: /etc/kubernetes/authn-webhook.conf
      name: authn-webhook-conf
      readOnly: true
    - mountPath: /etc/kubernetes/authz-webhook.conf
      name: authz-webhook-conf
      readOnly: true
# volumes
  - hostPath:
      path: /etc/kubernetes/authn-webhook.conf
      type: File
    name: authn-webhook-conf
  - hostPath:
      path: /etc/kubernetes/authz-webhook.conf
      type: File
    name: authz-webhook-conf
```

To configure the admission webhooks you have to apply the MutatingWebhookConfiguration and the ValidatingWebhookConfiguration:
```
kubectl apply -f mutatingwebhookconfiguration.yaml -f validatingwebhookconfiguration.yaml
```

The example configurations `authz-webhook.conf`, `authn-webhook.conf`, `mutatingwebhookconfiguration.yaml` and `validatingwebhookconfiguration.yaml` assume that the Kubernetes API-Server can reach your webhook server on the IP `172.18.0.1`.
If your setup is different you have to change the IP configurations accordingly (you can also use a hostname) and also create a certificate which contains your IP/hostname.
If you generate a new certificate you also have to update the `caBundle` in the `mutatingwebhookconfiguration.yaml` and `validatingwebhookconfiguration.yaml`.

To generate a new certificate you can use [pcert](https://github.com/dvob/pcert), OpenSSL or any other tool which allows to create certificates and set the SubjectAlternativeName accordingly.
```bash
# create server.crt and server.key
# use --dns instead if you use a hostname
pcert create server --ip <YOUR IP>
```

To generate the value for the `caBundle` in the `mutatingwebhookconfiguration.yaml` and `validatingwebhookconfiguration.yaml` you have to base64-encode the `server.crt`:
```
base64 -w 0 server.crt
```

### Setup with Kind
You can quickly start a cluster with the correct configuration using [kind](https://kind.sigs.k8s.io/).
To make the setup work with the example configurations make sure your kind Docker network uses the IP range `172.18.0.0/16`.
If this is not the case see above on how to configure a different IP.

Then you can simply start a cluster as follows:
```
kind create cluster --config kind-config.yaml
kubectl config use-context kind-kind
kubectl apply -f mutatingwebhookconfiguration.yaml -f validatingwebhookconfiguration.yaml
```

## Build & Run
* Build and run the webhook server
```
go build
./webhook
```

## Test
After everything is setup can test if your setup works.
The following commands assume that your cluster is configured in a context named `kind-kind` (default for the kind setup described above).
```bash
kubectl config use-context kind-kind

# set a user which uses the magic-token
kubectl config set-credentials kind-user --token magic-token
kubectl config set-context kind-kind --user kind-user

# User "magic-user" cannot list resource "pods"
kubectl get pod

# The user is allowed to list the configmaps because of our authorizer
kubectl get configmap

# We are not allowed to create a configmap which contains the value 'not-allowed-value' because of our validating admission
# error: failed to create configmap: admission webhook "configmap-example.puzzle.ch" denied the request: value 'not-allowed-value' not allowed in configmap
kubectl create configmap mytest --from-literal=not-allowed-value=abc

# create a empty configmap
kubectl create configmap mytest

# if we verify the contents we see that the mutating admission added a value 'magic-value=foobar' to our configmap mytest
kubectl get configmaps mytest -o yaml
```

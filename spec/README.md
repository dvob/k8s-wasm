# Module Specification

To pass data between the Kubernetes extension (host) and the WASM module we use the capabilities of [WASI](https://wasi.dev/) to read from the standard input and write to the standard output.
Specifically the module reads the input data from the standard input (stdin) and writes the result to the standard output (stdout).

In general the data passed over stdin and stdout is JSON encoded and has the following format.
The actual inner `requestData` and `responseData` depends on the use case (Authentication, Authorization, Admission) and is described in the following sections.

**Input** ([`Request`](https://github.com/dvob/kubernetes/blob/704f41c20a83b76e1542fca89b046cf854106df2/pkg/wasm/runner.go#L66)):
```json
{
	"request": requestData,
	"settings": {}
}
```
With the `settings` object you can pass parameters to your module. The `settings`, can be omitted if your module does not use them.
The settings can have any format (e.g. contain a list, or nested objects, etc.).

**Output** ([`Response`](https://github.com/dvob/kubernetes/blob/704f41c20a83b76e1542fca89b046cf854106df2/pkg/wasm/runner.go#L71)):
```json
{
	"response": responseData,
	"error": "error message"
}
```
If the `error` is set in the output, the action is considered as failed and the `response` is ignored.
Usually you send a result which contains either a `response` or an `error`.

If the module returns with a exit code not equal to `0` the action is considered as failed and the output is ignored.

# Use cases
To implement one of the uses cases the module has to export the appropriate function name.
If the function is called it has to read the input as described above and provide the appropriate output.

## Authentication
* Function Name: `authn`
* RequestData: [`v1.TokenReview`](https://pkg.go.dev/k8s.io/api/authentication/v1#TokenReview)
* ResponseData: [`v1.TokenReview`](https://pkg.go.dev/k8s.io/api/authentication/v1#TokenReview)

### Example
**Input**: User uses token `magic-token`:
```json
{
  "request": {
    "apiVersion": "authentication.k8s.io/v1",
    "kind": "TokenReview",
    "spec": {
      "token": "magic-token",
      "audiences": [
        "https://kubernetes.default.svc.cluster.local"
      ]
    }
  }
}
```

**Output**: Token gets authenticated:
```json
{
  "response": {
    "apiVersion": "authentication.k8s.io/v1",
    "kind": "TokenReview",
    "status": {
      "authenticated": true,
      "user": {
        "groups": [
          "magic-group"
        ],
        "uid": "0",
        "username": "magic-user"
      }
    }
  }
}
```

## Authorization
* Function Name: `authz`
* RequestData: [`v1.SubjectAccessReview`](https://pkg.go.dev/k8s.io/api/authorization/v1#SubjectAccessReview)
* ResponseData: [`v1.SubjectAccessReview`](https://pkg.go.dev/k8s.io/api/authorization/v1#SubjectAccessReview)

### Example
**Input**: User would like to list pods in the namespace default:
```json
{
  "request": {
    "apiVersion": "authorization.k8s.io/v1",
    "kind": "SubjectAccessReview",
    "spec": {
      "resourceAttributes": {
        "namespace": "default",
        "verb": "list",
        "version": "v1",
        "resource": "pods"
      },
      "uid": "0",
      "user": "magic-user",
      "groups": [
        "magic-group",
        "system:authenticated"
      ]
    }
  }
}
```

**Output**: Request is not authorized:
```json
{
  "response": {
    "apiVersion": "authorization.k8s.io/v1",
    "kind": "SubjectAccessReview",
    "status": {
      "allowed": false
    }
  }
}
```

## Admission
* Function Name: `validate`
* RequestData: [`v1.AdmissionReview`](https://pkg.go.dev/k8s.io/api/admission/v1#AdmissionReview)
* ResponseData: [`v1.AdmissionReview`](https://pkg.go.dev/k8s.io/api/admission/v1#AdmissionReview)

Both mutating and validating admissions use the same function name.

If you provide a patch in your response to implement a mutating admission the following difference to the official `AdmissionReview` applies.
The host does not understand the official patch type `JSONPatch` but it implements its on patch type `Full`.
Instead of a patch you provide the full edited object as patch.
This simplifies the module implementation because you don't have to generate a `JSONPatch` (RFC 6902).

### Example
#### Validating Admission
**Input**: User wants to create a config map:
```json
{
  "request": {
    "kind": "AdmissionReview",
    "apiVersion": "admission.k8s.io/v1",
    "request": {
      "uid": "678b2f02-0837-4262-95ea-5781b2864ac0",
      "kind": {
        "group": "",
        "version": "v1",
        "kind": "ConfigMap"
      },
      "resource": {
        "group": "",
        "version": "v1",
        "resource": "configmaps"
      },
      "requestKind": {
        "group": "",
        "version": "v1",
        "kind": "ConfigMap"
      },
      "requestResource": {
        "group": "",
        "version": "v1",
        "resource": "configmaps"
      },
      "name": "my-config",
      "namespace": "default",
      "operation": "CREATE",
      "userInfo": {
        "username": "magic-user",
        "uid": "0",
        "groups": [
          "magic-group",
          "system:authenticated"
        ]
      },
      "object": {
        "kind": "ConfigMap",
        "apiVersion": "v1",
        "metadata": {
          "name": "my-config",
          "namespace": "default",
          "uid": "869d6f79-dbe3-4ef0-8255-ef817e6e673e",
          "creationTimestamp": "2022-08-11T05:05:28Z"
        },
        "data": {
          "magic-value": "foobar",
          "not-allowed-value": "bar"
        }
      }
    }
  }
}
```

**Output**: Creation of configmap is rejected because it contains a value which is not allowed:
```json
{
  "response": {
    "kind": "AdmissionReview",
    "apiVersion": "admission.k8s.io/v1",
    "response": {
      "uid": "678b2f02-0837-4262-95ea-5781b2864ac0",
      "allowed": false,
      "status": {
        "apiVersion": "v1",
        "kind": "Status",
        "message": "value not-allowed-value not allowed in configmap"
      }
    }
  }
}
```

#### Mutating Admission
**Input**: User creates configmap:
```json
{
  "request": {
    "kind": "AdmissionReview",
    "apiVersion": "admission.k8s.io/v1",
    "request": {
      "uid": "695570da-9d1d-476a-a58a-15e051768042",
      "kind": {
        "group": "",
        "version": "v1",
        "kind": "ConfigMap"
      },
      "resource": {
        "group": "",
        "version": "v1",
        "resource": "configmaps"
      },
      "requestKind": {
        "group": "",
        "version": "v1",
        "kind": "ConfigMap"
      },
      "requestResource": {
        "group": "",
        "version": "v1",
        "resource": "configmaps"
      },
      "name": "my-config",
      "namespace": "default",
      "operation": "CREATE",
      "userInfo": {
        "username": "magic-user",
        "uid": "0",
        "groups": [
          "magic-group",
          "system:authenticated"
        ]
      },
      "object": {
        "kind": "ConfigMap",
        "apiVersion": "v1",
        "metadata": {
          "name": "my-config",
          "namespace": "default"
        },
        "data": {
          "not-allowed-value": "bar"
        }
      }
    }
  }
}
```

**Output**: Add `magic-value: foobar` to the configmap:
```json
{
  "response": {
    "kind": "AdmissionReview",
    "apiVersion": "admission.k8s.io/v1",
    "response": {
      "uid": "695570da-9d1d-476a-a58a-15e051768042",
      "allowed": true,
      "patchType": "Full",
      "patch": "eyJhcGlWZXJzaW9uIjoidjEiLCJraW5kIjoiQ29uZmlnTWFwIiwiZGF0YSI6eyJtYWdpYy12YWx1ZSI6ImZvb2JhciIsIm5vdC1hbGxvd2VkLXZhbHVlIjoiYmFyIn0sIm1ldGFkYXRhIjp7Im5hbWUiOiJteS1jb25maWciLCJuYW1lc3BhY2UiOiJkZWZhdWx0In19Cg=="
    }
  }
}
```

# Develop Modules

## Rust

See: https://github.com/dvob/k8s-wasi-rs

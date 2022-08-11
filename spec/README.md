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
The `settings`, can be omitted if your module does not use them.

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
Input:
```json
```

Output:
```json
```

## Authorization
* Function Name: `authz`
* RequestData: [`v1.SubjectAccessReview`](https://pkg.go.dev/k8s.io/api/authorization/v1#SubjectAccessReview)
* ResponseData: [`v1.SubjectAccessReview`](https://pkg.go.dev/k8s.io/api/authorization/v1#SubjectAccessReview)

### Example
Input:
```json
```

Output:
```json
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
Input:
```json
```

Output:
```json
```

#### Mutating Admission
Input:
```json
```

Output:
```json
```

# Develop Modules

## Rust

See: https://github.com/dvob/k8s-wasi-rs

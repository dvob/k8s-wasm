# Kubernetes Integration

There are two main ways how we can implement custom logic for the API-Server which gets called during Authentication, Authorization and Admission:

* Webhooks: Configure webhooks in the API-Server
* Direct: Include the custom logic in the API-Server code

To explore these to variants we implement the following logic in both ways:

* **Authentication**: If the token `magic-token` is provided the request is authenticated as user `magic-user` which is a member of the group `magic-group`.
* **Authorization**: Allow users which are member of the group `magic-group` to manage configmaps.
* **Validating Admission**: Reject configmaps which contain the value `not-allowed-value`.
* **Mutating Admission**: Add the value `magic-value: foobar` to all configmaps.

See the subdirectories for a description of the two implementations:
* [webhook](./webhook/)
* [api-server](./api-server/)

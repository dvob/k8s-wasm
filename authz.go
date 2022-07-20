package wasm

import (
	"context"

	k8s "k8s.io/apiserver/pkg/authorization/authorizer"
)

var _ k8s.Authorizer = (*Authorizer)(nil)

type Authorizer struct {
}

func NewAuthorizer() *Authorizer {
	return &Authorizer{}
}

func (a *Authorizer) Authorize(ctx context.Context, attrs k8s.Attributes) (authorized k8s.Decision, reason string, err error) {
	return k8s.DecisionDeny, "", nil
}

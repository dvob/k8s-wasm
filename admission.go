package wasm

import (
	"context"

	k8s "k8s.io/apiserver/pkg/admission"
)

var _ k8s.MutationInterface = (*Admiter)(nil)
var _ k8s.ValidationInterface = (*Admiter)(nil)

type Admiter struct {
}

func NewAdmiter() *Admiter {
	return &Admiter{}
}

func (a *Admiter) Handles(operation k8s.Operation) bool {
	return true
}

func (a *Admiter) Admit(ctx context.Context, attr k8s.Attributes, o k8s.ObjectInterfaces) (err error) {
	return nil
}

func (a *Admiter) Validate(ctx context.Context, attr k8s.Attributes, o k8s.ObjectInterfaces) (err error) {
	return nil
}

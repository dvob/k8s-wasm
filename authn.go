package wasm

import (
	"context"

	k8s "k8s.io/apiserver/pkg/authentication/authenticator"
)

var _ k8s.Token = (*Authenticator)(nil)

type AuthenticationConfig struct {
	Modules []AuthenticationModuleConfig `json:"modules"`
}

type AuthenticationModuleConfig struct {
	File     string `json:"file"`
	Settings any
}

type Authenticator struct {
}

func NewTokenAuthenticator() *Authenticator {
	return &Authenticator{}
}

func (t *Authenticator) AuthenticateToken(ctx context.Context, token string) (*k8s.Response, bool, error) {
	return nil, false, nil
}

package wasm

import (
	"context"
	"reflect"
	"testing"
)

const (
	authnTestModuleFile = "testdata/test-authn/target/wasm32-wasi/debug/test_authn.wasm"
	testToken           = "my-test-token"
	testUser            = "my-user"
	testUID             = "1337"
)

var (
	testGroups = []string{"system:masters"}
)

func newTestAuthenticator(t *testing.T) *Authenticator {
	config := &AuthenticationModuleConfig{
		File: authnTestModuleFile,
	}
	authenticator, err := NewAuthenticatorWithConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	return authenticator
}

func TestAuthenticatorSuccess(t *testing.T) {
	authenticator := newTestAuthenticator(t)
	ctx := context.Background()

	resp, ok, err := authenticator.AuthenticateToken(ctx, testToken)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Errorf("token '%s' should be authenticated", testToken)
	}

	if resp.User.GetName() != testUser {
		t.Errorf("wrong username: want=%s, got=%s", testUser, resp.User.GetName())
	}

	if resp.User.GetUID() != testUID {
		t.Errorf("wrong UID: want=%s, got=%s", testUID, resp.User.GetUID())
	}

	if !reflect.DeepEqual(resp.User.GetGroups(), testGroups) {
		t.Errorf("wrong groups: want=%s, got=%s", testGroups, resp.User.GetGroups())
	}
}

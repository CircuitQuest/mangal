package auth

import (
	"context"

	"golang.org/x/oauth2"
)

// LoginOption is a per Provider login option that handles
// authorization and providing the authorization token (generalized to an oauth2 token).
type LoginOption interface {
	// String the name of the login option, for logging purposes.
	String() string

	// Authorize will perform the authorization steps to obtain the oauth2 Token.
	Authorize(ctx context.Context) error

	// Token returns the authorization token.
	//
	// Calls to Token should not perform any HTTP requests.
	Token() *oauth2.Token
}

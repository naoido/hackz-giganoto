package authapi

import (
	"context"

	"goa.design/clue/log"
	auth "object-t.com/hackz-giganoto/microservices/auth/gen/auth"
)

// auth service example implementation.
// The example methods log the requests and return zero values.
type authsrvc struct{}

// NewAuth returns the auth service implementation.
func NewAuth() auth.Service {
	return &authsrvc{}
}

// Introspect opaque token and return internal JWT token for Kong Gateway
func (s *authsrvc) Introspect(ctx context.Context, p *auth.IntrospectPayload) (res *auth.IntrospectResult, err error) {
	res = &auth.IntrospectResult{}
	log.Printf(ctx, "auth.introspect")
	return
}

// Get GitHub OAuth authorization URL with state parameter
func (s *authsrvc) AuthURL(ctx context.Context) (res *auth.AuthURLResult, err error) {
	res = &auth.AuthURLResult{}
	log.Printf(ctx, "auth.auth_url")
	return
}

// Handle GitHub OAuth callback and return opaque token
func (s *authsrvc) OauthCallback(ctx context.Context, p *auth.OauthCallbackPayload) (res *auth.OauthCallbackResult, err error) {
	res = &auth.OauthCallbackResult{}
	log.Printf(ctx, "auth.oauth_callback")
	return
}

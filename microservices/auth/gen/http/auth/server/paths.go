// Code generated by goa v3.21.1, DO NOT EDIT.
//
// HTTP request path constructors for the auth service.
//
// Command:
// $ goa gen object-t.com/hackz-giganoto/microservices/auth/design

package server

// IntrospectAuthPath returns the URL path to the auth service introspect HTTP endpoint.
func IntrospectAuthPath() string {
	return "/introspect"
}

// AuthURLAuthPath returns the URL path to the auth service auth_url HTTP endpoint.
func AuthURLAuthPath() string {
	return "/auth/github"
}

// OauthCallbackAuthPath returns the URL path to the auth service oauth_callback HTTP endpoint.
func OauthCallbackAuthPath() string {
	return "/auth/github/callback"
}

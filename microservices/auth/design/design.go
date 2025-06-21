package design

import . "goa.design/goa/v3/dsl"

var _ = Service("auth", func() {
	Description("Authentication service that converts opaque tokens to internal JWT tokens for Kong Gateway")

	HTTP(func() {
		Path("/")
	})

	Method("introspect", func() {
		Description("Introspect opaque token and return internal JWT token for Kong Gateway")

		Payload(func() {
			Field(1, "token", String, "Opaque token to introspect")
			Required("token")
		})

		Result(func() {
			Field(1, "jwt", String, "Internal JWT token for downstream services")
			Required("jwt")
		})

		Error("invalid_token", String, "Token is invalid or expired")
		Error("internal_error", String, "Internal server error")

		HTTP(func() {
			POST("/introspect")
			Response(StatusOK)
			Response("invalid_token", StatusUnauthorized)
			Response("internal_error", StatusInternalServerError)
		})
	})

	Method("auth_url", func() {
		Description("Get GitHub OAuth authorization URL with state parameter")

		Result(func() {
			Field(1, "auth_url", String, "GitHub OAuth authorization URL")
			Field(2, "state", String, "OAuth state parameter for CSRF protection")
			Required("auth_url", "state")
		})

		Error("internal_error", String, "Internal server error")

		HTTP(func() {
			GET("/auth/github")
			Response(StatusOK)
			Response("internal_error", StatusInternalServerError)
		})
	})

	Method("oauth_callback", func() {
		Description("Handle GitHub OAuth callback and return opaque token")

		Payload(func() {
			Field(1, "code", String, "Authorization code from GitHub")
			Field(2, "state", String, "OAuth state parameter for validation")
			Required("code", "state")
		})

		Result(func() {
			Field(1, "access_token", String, "Opaque access token")
			Field(2, "token_type", String, "Token type (Bearer)")
			Field(3, "expires_in", Int64, "Token expiration in seconds")
			Field(4, "user_id", String, "GitHub user ID")
			Field(5, "userName", String, "GitHub user name")
			Required("access_token", "token_type", "expires_in", "user_id")
		})

		Error("github_error", String, "GitHub API error")
		Error("invalid_state", String, "Invalid or expired state parameter")
		Error("invalid_code", String, "Invalid authorization code")
		Error("github_error", String, "GitHub API error")
		Error("internal_error", String, "Internal server error")

		HTTP(func() {
			GET("/auth/github/callback")
			Param("code")
			Param("state")
			Response(StatusOK)
			Response("invalid_state", StatusBadRequest)
			Response("invalid_code", StatusBadRequest)
			Response("github_error", StatusBadGateway)
			Response("internal_error", StatusInternalServerError)
		})
	})
})

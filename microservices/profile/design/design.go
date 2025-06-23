package design

import . "goa.design/goa/v3/dsl"

var JWTAuth = JWTSecurity("jwt", func() {
	Description("JWT")
	Scope("api:read", "Read access to API resources")
	Scope("api:write", "Write access to API resources")
	Scope("api:register", "Register new user")
})

var _ = API("profile", func() {
	Title("Profile Service")
	Description("User profile management service")
	Version("1.0")

	Server("profile", func() {
		Host("localhost", func() {
			URI("grpc://localhost:50052")
		})
	})
})

var _ = Service("profile", func() {
	Description("User profile management service")

	GRPC(func() {
		Package("profile.v1")
	})

	Security(JWTAuth, func() {
		Scope("api:read")
		Scope("api:write")
	})

	Method("get_profile", func() {
		Description("Get user profile")

		Security(JWTAuth, func() {
			Scope("api:read")
		})

		Payload(func() {
			Token("token", String, "JWT token")
			Field(1, "user_id", String, "The user_id of the user")

			Required("token", "user_id")
		})

		Result(func() {
			Field(1, "user_id", String, "User ID")
			Field(2, "name", String, "User name")
			Required("user_id", "name")
		})

		Error("unauthorized", String, "Unauthorized access")
		Error("not_found", String, "Profile not found")
		Error("internal_error", String, "Internal server error")

		GRPC(func() {
			Response(CodeOK)
			Response("unauthorized", CodeUnauthenticated)
			Response("not_found", CodeNotFound)
			Response("internal_error", CodeInternal)
		})
	})

	Method("update_profile", func() {
		Description("Update user profile")

		Security(JWTAuth, func() {
			Scope("api:write")
		})

		Payload(func() {
			Token("token", String, "JWT token")
			Field(1, "name", String, "User name")
			Required("name")
		})

		Result(func() {
			Field(1, "user_id", String, "User ID")
			Field(2, "name", String, "User name")
			Required("user_id", "name")
		})

		Error("unauthorized", String, "Unauthorized access")
		Error("bad_request", String, "Invalid request")
		Error("internal_error", String, "Internal server error")

		GRPC(func() {
			Response(CodeOK)
			Response("unauthorized", CodeUnauthenticated)
			Response("bad_request", CodeInvalidArgument)
			Response("internal_error", CodeInternal)
		})
	})
})

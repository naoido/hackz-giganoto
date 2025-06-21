package design

import . "goa.design/goa/v3/dsl"

var JWTAuth = JWTSecurity("jwt", func() {
	Description("JWT")
	Scope("api:read", "Read access to API resources")
	Scope("api:write", "Write access to API resources")
})

var ChatMessage = Type("ChatMessage", func() {
	Description("Chat message")
	Field(1, "message_id", String, "Message ID")
	Field(2, "user_id", String, "Sender user ID")
	Field(3, "user_name", String, "Sender user name")
	Field(4, "message", String, "Message content")
	Field(5, "message_type", String, "Message type")
	Field(6, "timestamp", String, "Message timestamp")
	Required("message_id", "user_id", "user_name", "message", "timestamp")
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
		})

		Result(func() {
			Field(1, "user_id", String, "User ID")
			Field(2, "name", String, "User name")
			Field(3, "created_at", String, "Creation timestamp")
			Field(4, "updated_at", String, "Last update timestamp")
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
			Field(3, "created_at", String, "Creation timestamp")
			Field(4, "updated_at", String, "Last update timestamp")
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

	Method("get_profile_by_id", func() {
		Description("Get profile by user ID (for internal service communication)")

		Security(JWTAuth, func() {
			Scope("api:read")
		})

		Payload(func() {
			Token("token", String, "JWT token")
			Field(1, "user_id", String, "User ID to lookup")
			Required("user_id")
		})

		Result(func() {
			Field(1, "user_id", String, "User ID")
			Field(2, "name", String, "User name")
			Field(3, "creation_timestamp", String, "Creation timestamp")
			Field(4, "updated_at", String, "Last update timestamp")
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
})

var _ = Service("chat", func() {
	Description("Bidirectional chat service")

	GRPC(func() {
		Package("chat.v1")
	})

	Security(JWTAuth, func() {
		Scope("api:read")
		Scope("api:write")
	})

	Method("send_message", func() {
		Description("Send a message to a chat room")

		Security(JWTAuth, func() {
			Scope("api:write")
		})

		Payload(func() {
			Token("token", String, "JWT token")
			Field(1, "room_id", String, "Chat room ID")
			Field(2, "message", String, "Message content")
			Field(3, "message_type", String, "Message type (text, image, file)")
			Required("room_id", "message")
		})

		Result(func() {
			Field(1, "message_id", String, "Message ID")
			Field(2, "timestamp", String, "Message timestamp")
			Required("message_id", "timestamp")
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

	Method("join_chat", func() {
		Description("Join a chat room and receive messages")

		Security(JWTAuth, func() {
			Scope("api:read")
		})

		Payload(func() {
			Token("token", String, "JWT token")
			Field(1, "room_id", String, "Chat room ID")
			Required("room_id")
		})

		StreamingResult(func() {
			Field(1, "message_id", String, "Message ID")
			Field(2, "user_id", String, "Sender user ID")
			Field(3, "user_name", String, "Sender user name")
			Field(4, "message", String, "Message content")
			Field(5, "message_type", String, "Message type")
			Field(6, "timestamp", String, "Message timestamp")
			Required("message_id", "user_id", "user_name", "message", "timestamp")
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

	Method("get_chat_history", func() {
		Description("Get chat history for a room")

		Security(JWTAuth, func() {
			Scope("api:read")
		})

		Payload(func() {
			Token("token", String, "JWT token")
			Field(1, "room_id", String, "Chat room ID")
			Field(2, "limit", Int, "Number of messages to retrieve")
			Field(3, "offset", Int, "Offset for pagination")
			Required("room_id")
		})

		Result(func() {
			Field(1, "messages", ArrayOf("ChatMessage"), "Chat messages")
			Field(2, "total_count", Int, "Total number of messages")
			Required("messages", "total_count")
		})

		Error("unauthorized", String, "Unauthorized access")
		Error("bad_request", String, "Invalid request")
		Error("not_found", String, "Room not found")
		Error("internal_error", String, "Internal server error")

		GRPC(func() {
			Response(CodeOK)
			Response("unauthorized", CodeUnauthenticated)
			Response("bad_request", CodeInvalidArgument)
			Response("not_found", CodeNotFound)
			Response("internal_error", CodeInternal)
		})
	})
})

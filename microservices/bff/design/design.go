package design

import . "goa.design/goa/v3/dsl"

var JWTAuth = JWTSecurity("jwt", func() {
	Description("JWT")
	Scope("api:read", "Read access to API resources")
	Scope("api:write", "Write access to API resources")
})

var EnrichedMessage = Type("EnrichedMessage", func() {
	Description("Chat message enriched with user profile information")

	Field(1, "message_id", String, "Message ID")
	Field(2, "room_id", String, "Room ID")
	Field(3, "user_id", String, "Sender user ID")
	Field(4, "message", String, "Message content")
	Field(5, "created_at", Int64, "Sent timestamp")
	Field(6, "updated_at", Int64, "Created timestamp")
	Required("room_id", "user_id", "message")
})

var RoomInfo = Type("RoomInfo", func() {
	Description("Chat room information enriched with creator profile")

	Field(1, "room_id", String, "Room ID")
	Field(2, "name", String, "Room name")
	Field(3, "created_by", String, "User ID who created the room")
	Field(4, "creator_name", String, "Creator user name from profile")
	Field(5, "created_at", String, "Creation timestamp")
	Required("room_id", "name", "created_by", "creator_name", "created_at")
})

var _ = Service("bff", func() {
	Description("Backend for Frontend service for chat application")

	GRPC(func() {
		Package("bff.v1")
	})

	Security(JWTAuth, func() {
		Scope("api:read")
		Scope("api:write")
	})

	Method("create_room", func() {
		Description("Create a new chat room")

		Security(JWTAuth, func() {
			Scope("api:write")
		})

		Payload(func() {
			Token("token", String, "JWT token")
			Required("token")
		})

		Result(String)

		Error("unauthorized", String, "Unauthorized access")
		Error("internal_error", String, "Internal server error")

		GRPC(func() {
			Response(CodeOK)
			Response("unauthorized", CodeUnauthenticated)
			Response("internal_error", CodeInternal)
		})
	})

	Method("history", func() {
		Description("Get chat room history with enriched user names")

		Security(JWTAuth, func() {
			Scope("api:read")
		})

		Payload(func() {
			Token("token", String, "JWT token")
			Field(1, "room_id", String, "Room ID")
			Required("token", "room_id")
		})

		Result(ArrayOf(EnrichedMessage))

		Error("unauthorized", String, "Unauthorized access")
		Error("permission-denied", String, "Permission denied")
		Error("internal_error", String, "Internal server error")

		GRPC(func() {
			Response(CodeOK)
			Response("unauthorized", CodeUnauthenticated)
			Response("permission-denied", CodePermissionDenied)
			Response("internal_error", CodeInternal)
		})
	})

	Method("room-list", func() {
		Description("Get all chat rooms history")
		Security(JWTAuth, func() {
			Scope("api:read")
		})

		Payload(func() {
			Token("token", String, "The access token")
			Required("token")
		})

		Result(ArrayOf(String))

		Error("unauthorized", String, "Unauthorized access")

		GRPC(func() {
			Response(CodeOK)
			Response("unauthorized", CodeInvalidArgument)
		})
	})

	Method("join-room", func() {
		Description("Creates a new chat room")

		Security(JWTAuth, func() {
			Scope("api:read")
		})

		Payload(func() {
			Token("token", String, "The access token")
			Field(1, "invite_key", String, "Invite key")
			Required("token", "invite_key")
		})

		Result(String)

		Error("notfound", String)

		GRPC(func() {
			Response(CodeOK)
			Response("notfound", CodeNotFound)
		})
	})

	Method("invite-room", func() {
		Description("Creates a new chat room")

		Security(JWTAuth, func() {
			Scope("api:read")
		})

		Payload(func() {
			Token("token", String, "The access token")
			Field(1, "room_id", String, "The id of the room")
			Field(2, "user_id", String, "The id of the user")

			Required("token", "room_id", "user_id")
		})

		Result(String)

		Error("invalid_argument", String)

		GRPC(func() {
			Response(CodeOK)
			Response("invalid_argument", CodeInvalidArgument)
		})
	})

	Method("stream_chat", func() {
		Description("Stream chat messages with bidirectional communication")

		Security(JWTAuth, func() {
			Scope("api:read")
		})

		Payload(func() {
			Token("token", String, "JWT token")
			Field(1, "room_id", String, "Room ID")
			Required("token", "room_id")
		})

		StreamingPayload(String)
		StreamingResult(EnrichedMessage)

		Error("unauthorized", String, "Unauthorized access")
		Error("permission-denied", String, "Permission denied")
		Error("internal_error", String, "Internal server error")

		GRPC(func() {
			Response(CodeOK)
			Response("unauthorized", CodeUnauthenticated)
			Response("permission-denied", CodePermissionDenied)
			Response("internal_error", CodeInternal)
		})
	})

	Method("get_profile", func() {
		Description("Get current user profile")

		Security(JWTAuth, func() {
			Scope("api:read")
		})

		Payload(func() {
			Token("token", String, "JWT token")
			Field(1, "user_id", String, "User ID")

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
		Description("Update current user profile")

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

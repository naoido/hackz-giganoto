package design

import . "goa.design/goa/v3/dsl"

var JWTAuth = JWTSecurity("jwt", func() {
	Description("JWT")
	Scope("api:read", "Read access to API resources")
	Scope("api:write", "Write access to API resources")
})

var Chat = Type("Chat", func() {
	Description("Chat message")

	Field(1, "user_id", String, "user_id")
	Field(2, "message", String, "Message content")
	Field(3, "id", String, "ID")
	Field(4, "created_at", Int64, "Created timestamp")
	Field(5, "updated_at", Int64, "Updated timestamp")
	Field(6, "room_id", String, "room")
	Required("user_id", "message", "id", "created_at", "updated_at", "room_id")
})

var _ = API("chat", func() {
	Title("Chat Service")
	Description("Real-time chat service")
	Version("1.0")

	Server("chat", func() {
		Host("localhost", func() {
			URI("grpc://localhost:50053")
		})
	})
})

var _ = Service("chat", func() {
	Description("Real-time chat service with bidirectional streaming")

	GRPC(func() {
		Package("chat.v1")
	})

	Error("unauthorized", String)
	Error("permission-denied", String)
	Error("internal", String)

	Method("create-room", func() {
		Description("Creates a new chat room")

		Security(JWTAuth, func() {
			Scope("api:write")
		})

		Payload(func() {
			Token("token", String, "The access token")
			Required("token")
		})

		Result(String)

		GRPC(func() {
			Response(CodeOK)
			Response("internal", CodeInvalidArgument)
		})
	})

	Method("history", func() {
		Description("Get all chat rooms history")

		Security(JWTAuth, func() {
			Scope("api:read")
		})

		Payload(func() {
			Token("token", String, "The access token")
			Field(1, "room_id", String, "The id of the room")
			Required("token", "room_id")
		})

		Result(ArrayOf(Chat))

		GRPC(func() {
			Response(CodeOK)
			Response("unauthorized", CodeInvalidArgument)
			Response("permission-denied", CodePermissionDenied)
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

		Error("unauthorized", String)

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

	Method("stream-room", func() {
		Description("Streams chat room events on a chat room")

		Security(JWTAuth, func() {
			Scope("api:read")
		})

		Payload(func() {
			Token("token", String, "The access token")
			Field(1, "room_id", String, "The room id")
			Required("token", "room_id")
		})

		StreamingPayload(String)
		StreamingResult(Chat)

		GRPC(func() {
			Response(CodeOK)
			Response("unauthorized", CodeInvalidArgument)
			Response("permission-denied", CodePermissionDenied)
		})
	})
})

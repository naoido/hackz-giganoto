package design

import (
	. "goa.design/goa/v3/dsl"
)

var JWTAuth = JWTSecurity("jwt", func() {
	Description("JWT")
	Scope("api:read", "Read access to API resources")
	Scope("api:write", "Write access to API resources")
})

var _ = API("wiki", func() {
	Title("Wiki Service")
	Description("Wiki management service")
	Version("1.0")

	Server("wiki", func() {
		Host("localhost", func() {
			URI("grpc://localhost:50053")
		})
	})
})

var _ = Service("wiki", func() {
	Description("Wikiサービス")

	GRPC(func() {
		Package("wiki.v1")
	})

	Security(JWTAuth, func() {
		Scope("api:read")
		Scope("api:write")
	})

	Method("show", func() {
		Description("Wiki取得")

		Security(JWTAuth, func() {
			Scope("api:read")
		})

		Payload(func() {
			Token("token", String, "JWT token")
			Field(1, "id", String, "WikiのユニークID")
			Required("id")
		})

		Result(Wiki)

		Error("unauthorized", String, "Unauthorized access")
		Error("not_found", String, "Wiki not found")
		Error("internal_error", String, "Internal server error")

		GRPC(func() {
			Response(CodeOK)
			Response("unauthorized", CodeUnauthenticated)
			Response("not_found", CodeNotFound)
			Response("internal_error", CodeInternal)
		})
	})

	Method("create", func() {
		Description("新しいWikiを作成")

		Security(JWTAuth, func() {
			Scope("api:write")
		})

		Payload(func() {
			Token("token", String, "JWT token")
			Field(1, "channel_id", String, "Wikiを作成するチャンネルのID")
			Field(2, "content", String, "Markdown形式のコンテンツ")
			Required("channel_id", "content")
		})

		Result(Wiki)

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

	Method("update", func() {
		Description("指定されたIDのWikiの内容を更新")

		Security(JWTAuth, func() {
			Scope("api:write")
		})

		Payload(func() {
			Token("token", String, "JWT token")
			Field(1, "id", String, "更新するWikiのユニークID")
			Field(2, "content", String, "更新後のMarkdownコンテンツ")
			Required("id", "content")
		})

		Result(Wiki)

		Error("unauthorized", String, "Unauthorized access")
		Error("bad_request", String, "Invalid request")
		Error("not_found", String, "Wiki not found")
		Error("internal_error", String, "Internal server error")

		GRPC(func() {
			Response(CodeOK)
			Response("unauthorized", CodeUnauthenticated)
			Response("bad_request", CodeInvalidArgument)
			Response("not_found", CodeNotFound)
			Response("internal_error", CodeInternal)
		})
	})

	Method("delete", func() {
		Description("指定されたIDのWikiを削除")

		Security(JWTAuth, func() {
			Scope("api:write")
		})

		Payload(func() {
			Token("token", String, "JWT token")
			Field(1, "id", String, "削除するWikiのユニークID")
			Required("id")
		})

		Error("unauthorized", String, "Unauthorized access")
		Error("bad_request", String, "Invalid request")
		Error("not_found", String, "Wiki not found")
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

var Wiki = ResultType("application/vnd.wiki", func() {
	Description("チャンネルWikiを表します。")
	Attributes(func() {
		Field(1, "id", String, "Wiki毎のユニークな値", func() {
			Example("wiki_12345abcde")
		})
		Field(2, "channel_id", String, "作成されたチャンネルのID", func() {
			Example("channel_open_67890fghij")
		})
		Field(3, "content", String, "Markdown形式のコンテンツ", func() {
			Example("# はじめに\n\nこれはテストです。")
		})
		Field(4, "created_at", String, "作成日時 (RFC3339 format)", func() {
			Format(FormatDateTime)
			Example("2025-06-21T18:40:43Z")
		})
		Field(5, "updated_at", String, "更新日時 (RFC3339 format)", func() {
			Format(FormatDateTime)
			Example("2025-06-21T18:40:43Z")
		})
	})
	Required("id", "channel_id", "content", "created_at", "updated_at")
})

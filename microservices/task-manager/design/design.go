package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = API("task-manager", func() {
	Title("Task Manager API")
	Description("GitHub Issue風のタスク管理システム")
	Version("1.0")
	Server("task-manager", func() {
		Host("localhost", func() {
			URI("grpc://localhost:8080")
		})
	})
})

// TaskResult はタスクの情報を返す型
var TaskResult = ResultType("application/vnd.task", func() {
	Field(1, "id", String, "タスクのユニークID")
	Field(2, "title", String, "タイトル")
	Field(3, "description", String, "説明")
	Field(4, "status", String, "ステータス (OPEN, IN_PROGRESS, DONE)")
	Field(5, "assignee_ids", ArrayOf(String), "担当者のユーザーIDリスト")
	Field(6, "label_ids", ArrayOf(String), "ラベルのIDリスト")
	Field(7, "created_at", String, "作成日時 (RFC3339)")
	Field(8, "updated_at", String, "更新日時 (RFC3339)")
	Required("id", "title", "status", "created_at", "updated_at")
})

// CommentResult はコメントの情報を返す型
var CommentResult = ResultType("application/vnd.comment", func() {
	Field(1, "id", String, "コメントのユニークID")
	Field(2, "task_id", String, "紐づくタスクのID")
	Field(3, "author_id", String, "投稿者のユーザーID")
	Field(4, "body", String, "コメント本文")
	Field(5, "created_at", String, "作成日時 (RFC3339)")
	Required("id", "task_id", "author_id", "body", "created_at")
})

// LabelResult はラベルの情報を返す型
var LabelResult = ResultType("application/vnd.label", func() {
	Field(1, "id", String, "ラベルのユニークID")
	Field(2, "name", String, "ラベル名")
	Field(3, "color", String, "色 (例: #d73a4a)")
	Required("id", "name", "color")
})

// TaskService はタスクを管理するサービス
var _ = Service("TaskService", func() {
	Description("タスクを管理するサービス")

	Method("List", func() {
		Description("タスクの一覧を取得します。フィルタリングも可能です。")
		Payload(func() {
			Field(1, "status", String, "ステータスでフィルタリング")
			Field(2, "assignee_id", String, "担当者IDでフィルタリング")
			Field(3, "label_id", String, "ラベルIDでフィルタリング")
		})
		Result(CollectionOf(TaskResult))
		GRPC(func() {})
	})

	Method("Get", func() {
		Description("指定されたIDのタスクを取得します。")
		Payload(func() {
			Field(1, "id", String, "タスクID")
			Required("id")
		})
		Result(TaskResult)
		Error("not_found", func() {
			Description("タスクが見つかりません")
		})
		GRPC(func() {})
	})

	Method("Create", func() {
		Description("新しいタスクを作成します。")
		Payload(func() {
			Field(1, "title", String, "タイトル")
			Field(2, "description", String, "説明")
			Field(3, "assignee_ids", ArrayOf(String), "担当者のユーザーIDリスト")
			Field(4, "label_ids", ArrayOf(String), "ラベルのIDリスト")
			Required("title")
		})
		Result(TaskResult)
		GRPC(func() {})
	})

	Method("Update", func() {
		Description("既存のタスクを更新します。")
		Payload(func() {
			Field(1, "id", String, "タスクID")
			Field(2, "title", String, "タイトル")
			Field(3, "description", String, "説明")
			Field(4, "status", String, "ステータス")
			Field(5, "assignee_ids", ArrayOf(String), "担当者のユーザーIDリスト")
			Field(6, "label_ids", ArrayOf(String), "ラベルのIDリスト")
			Required("id")
		})
		Result(TaskResult)
		Error("not_found", func() {
			Description("タスクが見つかりません")
		})
		GRPC(func() {})
	})

	Method("Delete", func() {
		Description("タスクを削除します。")
		Payload(func() {
			Field(1, "id", String, "タスクID")
			Required("id")
		})
		Error("not_found", func() {
			Description("タスクが見つかりません")
		})
		GRPC(func() {})
	})
})

// CommentService はコメントを管理するサービス
var _ = Service("CommentService", func() {
	Description("コメントを管理するサービス")

	Method("List", func() {
		Description("指定されたタスクのコメント一覧を取得します。")
		Payload(func() {
			Field(1, "task_id", String, "タスクID")
			Required("task_id")
		})
		Result(CollectionOf(CommentResult))
		GRPC(func() {})
	})

	Method("Get", func() {
		Description("指定されたIDのコメントを取得します。")
		Payload(func() {
			Field(1, "id", String, "コメントID")
			Required("id")
		})
		Result(CommentResult)
		Error("not_found", func() {
			Description("コメントが見つかりません")
		})
		GRPC(func() {})
	})

	Method("Create", func() {
		Description("新しいコメントを作成します。")
		Payload(func() {
			Field(1, "task_id", String, "紐づくタスクのID")
			Field(2, "author_id", String, "投稿者のユーザーID")
			Field(3, "body", String, "コメント本文")
			Required("task_id", "author_id", "body")
		})
		Result(CommentResult)
		GRPC(func() {})
	})

	Method("Update", func() {
		Description("既存のコメントを更新します。")
		Payload(func() {
			Field(1, "id", String, "コメントID")
			Field(2, "body", String, "コメント本文")
			Required("id", "body")
		})
		Result(CommentResult)
		Error("not_found", func() {
			Description("コメントが見つかりません")
		})
		GRPC(func() {})
	})

	Method("Delete", func() {
		Description("コメントを削除します。")
		Payload(func() {
			Field(1, "id", String, "コメントID")
			Required("id")
		})
		Error("not_found", func() {
			Description("コメントが見つかりません")
		})
		GRPC(func() {})
	})
})

// LabelService はラベルを管理するサービス
var _ = Service("LabelService", func() {
	Description("ラベルを管理するサービス")

	Method("List", func() {
		Description("ラベルの一覧を取得します。")
		Result(CollectionOf(LabelResult))
		GRPC(func() {})
	})

	Method("Get", func() {
		Description("指定されたIDのラベルを取得します。")
		Payload(func() {
			Field(1, "id", String, "ラベルID")
			Required("id")
		})
		Result(LabelResult)
		Error("not_found", func() {
			Description("ラベルが見つかりません")
		})
		GRPC(func() {})
	})

	Method("Create", func() {
		Description("新しいラベルを作成します。")
		Payload(func() {
			Field(1, "name", String, "ラベル名")
			Field(2, "color", String, "色 (例: #d73a4a)")
			Required("name", "color")
		})
		Result(LabelResult)
		GRPC(func() {})
	})

	Method("Update", func() {
		Description("既存のラベルを更新します。")
		Payload(func() {
			Field(1, "id", String, "ラベルID")
			Field(2, "name", String, "ラベル名")
			Field(3, "color", String, "色 (例: #d73a4a)")
			Required("id")
		})
		Result(LabelResult)
		Error("not_found", func() {
			Description("ラベルが見つかりません")
		})
		GRPC(func() {})
	})

	Method("Delete", func() {
		Description("ラベルを削除します。")
		Payload(func() {
			Field(1, "id", String, "ラベルID")
			Required("id")
		})
		Error("not_found", func() {
			Description("ラベルが見つかりません")
		})
		GRPC(func() {})
	})
})
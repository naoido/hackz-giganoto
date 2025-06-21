package comment_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"goa.design/clue/log"
	commentservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/comment_service"
	"object-t.com/hackz-giganoto/microservices/task-manager/repository"
)

// commentServiceImpl implements the CommentService interface
type commentServiceImpl struct {
	logger      func(context.Context, ...any)
	commentRepo repository.CommentRepository
}

// New returns the CommentService service implementation
func New(logger func(context.Context, ...any)) commentservice.Service {
	// モックリポジトリを作成（実際の実装では依存性注入を使用）
	commentRepo := NewMockCommentRepository()

	return &commentServiceImpl{
		logger:      logger,
		commentRepo: commentRepo,
	}
}

// List implements the "List" method of the "CommentService" service
func (s *commentServiceImpl) List(ctx context.Context, p *commentservice.ListPayload) (res commentservice.CommentCollection, err error) {
	log.Print(ctx, log.KV{K: "method", V: "commentService.List"}, log.KV{K: "taskId", V: p.TaskID})

	// リポジトリからコメント一覧を取得
	entities, err := s.commentRepo.List(ctx, p.TaskID)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to list comments"})
		return nil, err
	}

	// エンティティをGoaの結果型に変換
	res = make(commentservice.CommentCollection, len(entities))
	for i, entity := range entities {
		res[i] = repository.ConvertCommentEntityToResult(entity)
	}

	log.Print(ctx, log.KV{K: "count", V: len(res)}, log.KV{K: "message", V: "Listed comments"})
	return res, nil
}

// Get implements the "Get" method of the "CommentService" service
func (s *commentServiceImpl) Get(ctx context.Context, p *commentservice.GetPayload) (res *commentservice.Comment, err error) {
	log.Print(ctx, log.KV{K: "method", V: "commentService.Get"}, log.KV{K: "id", V: p.ID})

	// リポジトリからコメントを取得
	entity, err := s.commentRepo.Get(ctx, p.ID)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to get comment"})
		return nil, commentservice.MakeNotFound(err)
	}

	// エンティティをGoaの結果型に変換
	res = repository.ConvertCommentEntityToResult(entity)

	log.Print(ctx, log.KV{K: "message", V: "Got comment"}, log.KV{K: "id", V: res.ID})
	return res, nil
}

// Create implements the "Create" method of the "CommentService" service
func (s *commentServiceImpl) Create(ctx context.Context, p *commentservice.CreatePayload) (res *commentservice.Comment, err error) {
	log.Print(ctx, log.KV{K: "method", V: "commentService.Create"}, log.KV{K: "taskId", V: p.TaskID}, log.KV{K: "authorId", V: p.AuthorID})

	now := time.Now()
	entity := &repository.CommentEntity{
		ID:        uuid.New().String(),
		TaskID:    p.TaskID,
		AuthorID:  p.AuthorID,
		Body:      p.Body,
		CreatedAt: now,
	}

	// リポジトリにコメントを作成
	err = s.commentRepo.Create(ctx, entity)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to create comment"})
		return nil, err
	}

	// エンティティをGoaの結果型に変換
	res = repository.ConvertCommentEntityToResult(entity)

	log.Print(ctx, log.KV{K: "message", V: "Created comment"}, log.KV{K: "id", V: res.ID})
	return res, nil
}

// Update implements the "Update" method of the "CommentService" service
func (s *commentServiceImpl) Update(ctx context.Context, p *commentservice.UpdatePayload) (res *commentservice.Comment, err error) {
	log.Print(ctx, log.KV{K: "method", V: "commentService.Update"}, log.KV{K: "id", V: p.ID})

	// 既存のコメントを取得
	entity, err := s.commentRepo.Get(ctx, p.ID)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to get comment for update"})
		return nil, commentservice.MakeNotFound(err)
	}

	// 本文を更新
	entity.Body = p.Body

	// リポジトリでコメントを更新
	err = s.commentRepo.Update(ctx, entity)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to update comment"})
		return nil, err
	}

	// エンティティをGoaの結果型に変換
	res = repository.ConvertCommentEntityToResult(entity)

	log.Print(ctx, log.KV{K: "message", V: "Updated comment"}, log.KV{K: "id", V: res.ID})
	return res, nil
}

// Delete implements the "Delete" method of the "CommentService" service
func (s *commentServiceImpl) Delete(ctx context.Context, p *commentservice.DeletePayload) (err error) {
	log.Print(ctx, log.KV{K: "method", V: "commentService.Delete"}, log.KV{K: "id", V: p.ID})

	// コメントが存在するかチェック
	_, err = s.commentRepo.Get(ctx, p.ID)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to get comment for delete"})
		return commentservice.MakeNotFound(err)
	}

	// リポジトリからコメントを削除
	err = s.commentRepo.Delete(ctx, p.ID)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to delete comment"})
		return err
	}

	log.Print(ctx, log.KV{K: "message", V: "Deleted comment"}, log.KV{K: "id", V: p.ID})
	return nil
}

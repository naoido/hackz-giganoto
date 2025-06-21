package task_service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"goa.design/clue/log"
	taskservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/task_service"
	"object-t.com/hackz-giganoto/microservices/task-manager/repository"
)

// taskServiceImpl implements the TaskService interface
type taskServiceImpl struct {
	logger   func(context.Context, ...any)
	taskRepo repository.TaskRepository
}

// New returns the TaskService service implementation
func New(logger func(context.Context, ...any)) taskservice.Service {
	// モックリポジトリを作成（実際の実装では依存性注入を使用）
	taskRepo := NewMockTaskRepository()

	return &taskServiceImpl{
		logger:   logger,
		taskRepo: taskRepo,
	}
}

// List implements the "List" method of the "TaskService" service
func (s *taskServiceImpl) List(ctx context.Context, p *taskservice.ListPayload) (res taskservice.TaskCollection, err error) {
	s.logger(ctx, "taskService.List")

	// リポジトリからタスク一覧を取得
	entities, err := s.taskRepo.List(ctx, p.Status, p.AssigneeID, p.LabelID)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to list tasks"})
		return nil, err
	}

	// エンティティをGoaの結果型に変換
	res = make(taskservice.TaskCollection, len(entities))
	for i, entity := range entities {
		res[i] = repository.ConvertTaskEntityToResult(entity)
	}

	log.Print(ctx, log.KV{K: "count", V: len(res)}, log.KV{K: "message", V: "Listed tasks"})
	return res, nil
}

// Get implements the "Get" method of the "TaskService" service
func (s *taskServiceImpl) Get(ctx context.Context, p *taskservice.GetPayload) (res *taskservice.Task, err error) {
	log.Print(ctx, log.KV{K: "method", V: "taskService.Get"}, log.KV{K: "id", V: p.ID})

	// リポジトリからタスクを取得
	entity, err := s.taskRepo.Get(ctx, p.ID)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to get task"})
		return nil, taskservice.MakeNotFound(err)
	}

	// エンティティをGoaの結果型に変換
	res = repository.ConvertTaskEntityToResult(entity)

	log.Print(ctx, log.KV{K: "message", V: "Got task"}, log.KV{K: "title", V: res.Title})
	return res, nil
}

// Create implements the "Create" method of the "TaskService" service
func (s *taskServiceImpl) Create(ctx context.Context, p *taskservice.CreatePayload) (res *taskservice.Task, err error) {
	log.Print(ctx, log.KV{K: "method", V: "taskService.Create"}, log.KV{K: "title", V: p.Title})

	now := time.Now()
	entity := &repository.TaskEntity{
		ID:          uuid.New().String(),
		Title:       p.Title,
		Description: p.Description,
		Status:      "OPEN", // デフォルトステータス
		AssigneeIds: p.AssigneeIds,
		LabelIds:    p.LabelIds,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// nilスライスを空スライスに変換
	if entity.AssigneeIds == nil {
		entity.AssigneeIds = []string{}
	}
	if entity.LabelIds == nil {
		entity.LabelIds = []string{}
	}

	// リポジトリにタスクを作成
	err = s.taskRepo.Create(ctx, entity)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to create task"})
		return nil, err
	}

	// エンティティをGoaの結果型に変換
	res = repository.ConvertTaskEntityToResult(entity)

	log.Print(ctx, log.KV{K: "message", V: "Created task"}, log.KV{K: "id", V: res.ID})
	return res, nil
}

// Update implements the "Update" method of the "TaskService" service
func (s *taskServiceImpl) Update(ctx context.Context, p *taskservice.UpdatePayload) (res *taskservice.Task, err error) {
	log.Print(ctx, log.KV{K: "method", V: "taskService.Update"}, log.KV{K: "id", V: p.ID})

	// 既存のタスクを取得
	entity, err := s.taskRepo.Get(ctx, p.ID)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to get task for update"})
		return nil, taskservice.MakeNotFound(err)
	}

	// 更新フィールドをチェックして適用
	if p.Title != nil {
		entity.Title = *p.Title
	}
	if p.Description != nil {
		entity.Description = p.Description
	}
	if p.Status != nil {
		entity.Status = *p.Status
	}
	if p.AssigneeIds != nil {
		entity.AssigneeIds = p.AssigneeIds
	} else {
		entity.AssigneeIds = []string{}
	}
	if p.LabelIds != nil {
		entity.LabelIds = p.LabelIds
	} else {
		entity.LabelIds = []string{}
	}

	entity.UpdatedAt = time.Now()

	// リポジトリでタスクを更新
	err = s.taskRepo.Update(ctx, entity)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to update task"})
		return nil, err
	}

	// エンティティをGoaの結果型に変換
	res = repository.ConvertTaskEntityToResult(entity)

	log.Print(ctx, log.KV{K: "message", V: "Updated task"}, log.KV{K: "id", V: res.ID})
	return res, nil
}

// Delete implements the "Delete" method of the "TaskService" service
func (s *taskServiceImpl) Delete(ctx context.Context, p *taskservice.DeletePayload) (err error) {
	log.Print(ctx, log.KV{K: "method", V: "taskService.Delete"}, log.KV{K: "id", V: p.ID})

	// タスクが存在するかチェック
	_, err = s.taskRepo.Get(ctx, p.ID)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to get task for delete"})
		return taskservice.MakeNotFound(err)
	}

	// リポジトリからタスクを削除
	err = s.taskRepo.Delete(ctx, p.ID)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to delete task"})
		return err
	}

	log.Print(ctx, log.KV{K: "message", V: "Deleted task"}, log.KV{K: "id", V: p.ID})
	return nil
}

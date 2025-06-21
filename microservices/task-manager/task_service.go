package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	taskservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/task_service"
	"object-t.com/hackz-giganoto/microservices/task-manager/repository"
)

// taskServiceImpl implements the TaskService interface
type taskServiceImpl struct {
	logger   *log.Logger
	taskRepo repository.TaskRepository
}

// NewTaskService returns the TaskService service implementation
func NewTaskService(logger *log.Logger, taskRepo repository.TaskRepository) taskservice.Service {
	return &taskServiceImpl{
		logger:   logger,
		taskRepo: taskRepo,
	}
}

// List implements the "List" method of the "TaskService" service
func (s *taskServiceImpl) List(ctx context.Context, p *taskservice.ListPayload) (res taskservice.TaskCollection, err error) {
	s.logger.Print("taskService.List")

	// リポジトリからタスク一覧を取得
	entities, err := s.taskRepo.List(ctx, p.Status, p.AssigneeID, p.LabelID)
	if err != nil {
		s.logger.Printf("Failed to list tasks: %v", err)
		return nil, err
	}

	// エンティティをGoaの結果型に変換
	res = make(taskservice.TaskCollection, len(entities))
	for i, entity := range entities {
		res[i] = repository.ConvertTaskEntityToResult(entity)
	}

	s.logger.Printf("Listed %d tasks", len(res))
	return res, nil
}

// Get implements the "Get" method of the "TaskService" service
func (s *taskServiceImpl) Get(ctx context.Context, p *taskservice.GetPayload) (res *taskservice.Task, err error) {
	s.logger.Printf("taskService.Get: id=%s", p.ID)

	// リポジトリからタスクを取得
	entity, err := s.taskRepo.Get(ctx, p.ID)
	if err != nil {
		s.logger.Printf("Failed to get task: %v", err)
		return nil, taskservice.MakeNotFound(err)
	}

	// エンティティをGoaの結果型に変換
	res = repository.ConvertTaskEntityToResult(entity)

	s.logger.Printf("Got task: %s", res.Title)
	return res, nil
}

// Create implements the "Create" method of the "TaskService" service
func (s *taskServiceImpl) Create(ctx context.Context, p *taskservice.CreatePayload) (res *taskservice.Task, err error) {
	s.logger.Printf("taskService.Create: title=%s", p.Title)

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
		s.logger.Printf("Failed to create task: %v", err)
		return nil, err
	}

	// エンティティをGoaの結果型に変換
	res = repository.ConvertTaskEntityToResult(entity)

	s.logger.Printf("Created task: %s", res.ID)
	return res, nil
}

// Update implements the "Update" method of the "TaskService" service
func (s *taskServiceImpl) Update(ctx context.Context, p *taskservice.UpdatePayload) (res *taskservice.Task, err error) {
	s.logger.Printf("taskService.Update: id=%s", p.ID)

	// 既存のタスクを取得
	entity, err := s.taskRepo.Get(ctx, p.ID)
	if err != nil {
		s.logger.Printf("Failed to get task for update: %v", err)
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
		s.logger.Printf("Failed to update task: %v", err)
		return nil, err
	}

	// エンティティをGoaの結果型に変換
	res = repository.ConvertTaskEntityToResult(entity)

	s.logger.Printf("Updated task: %s", res.ID)
	return res, nil
}

// Delete implements the "Delete" method of the "TaskService" service
func (s *taskServiceImpl) Delete(ctx context.Context, p *taskservice.DeletePayload) (err error) {
	s.logger.Printf("taskService.Delete: id=%s", p.ID)

	// タスクが存在するかチェック
	_, err = s.taskRepo.Get(ctx, p.ID)
	if err != nil {
		s.logger.Printf("Failed to get task for delete: %v", err)
		return taskservice.MakeNotFound(err)
	}

	// リポジトリからタスクを削除
	err = s.taskRepo.Delete(ctx, p.ID)
	if err != nil {
		s.logger.Printf("Failed to delete task: %v", err)
		return err
	}

	s.logger.Printf("Deleted task: %s", p.ID)
	return nil
}

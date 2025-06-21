package task_service

import (
	"context"
	"errors"
	"sync"

	"object-t.com/hackz-giganoto/microservices/task-manager/repository"
)

// MockTaskRepository provides an in-memory implementation of TaskRepository
type MockTaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*repository.TaskEntity
}

func NewMockTaskRepository() repository.TaskRepository {
	return &MockTaskRepository{
		tasks: make(map[string]*repository.TaskEntity),
	}
}

func (r *MockTaskRepository) List(ctx context.Context, status, assigneeID, labelID *string) ([]*repository.TaskEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []*repository.TaskEntity
	for _, task := range r.tasks {
		// フィルタリング条件をチェック
		if status != nil && task.Status != *status {
			continue
		}
		if assigneeID != nil {
			found := false
			for _, id := range task.AssigneeIds {
				if id == *assigneeID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if labelID != nil {
			found := false
			for _, id := range task.LabelIds {
				if id == *labelID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		result = append(result, task)
	}
	return result, nil
}

func (r *MockTaskRepository) Get(ctx context.Context, id string) (*repository.TaskEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	task, exists := r.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (r *MockTaskRepository) Create(ctx context.Context, task *repository.TaskEntity) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.tasks[task.ID] = task
	return nil
}

func (r *MockTaskRepository) Update(ctx context.Context, task *repository.TaskEntity) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.tasks[task.ID]; !exists {
		return errors.New("task not found")
	}
	r.tasks[task.ID] = task
	return nil
}

func (r *MockTaskRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.tasks[id]; !exists {
		return errors.New("task not found")
	}
	delete(r.tasks, id)
	return nil
}
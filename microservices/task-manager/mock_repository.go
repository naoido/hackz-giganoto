package main

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

// MockCommentRepository provides an in-memory implementation of CommentRepository
type MockCommentRepository struct {
	mu       sync.RWMutex
	comments map[string]*repository.CommentEntity
}

func NewMockCommentRepository() repository.CommentRepository {
	return &MockCommentRepository{
		comments: make(map[string]*repository.CommentEntity),
	}
}

func (r *MockCommentRepository) List(ctx context.Context, taskID string) ([]*repository.CommentEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*repository.CommentEntity
	for _, comment := range r.comments {
		if comment.TaskID == taskID {
			result = append(result, comment)
		}
	}
	return result, nil
}

func (r *MockCommentRepository) Get(ctx context.Context, id string) (*repository.CommentEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, exists := r.comments[id]
	if !exists {
		return nil, errors.New("comment not found")
	}
	return comment, nil
}

func (r *MockCommentRepository) Create(ctx context.Context, comment *repository.CommentEntity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.comments[comment.ID] = comment
	return nil
}

func (r *MockCommentRepository) Update(ctx context.Context, comment *repository.CommentEntity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.comments[comment.ID]; !exists {
		return errors.New("comment not found")
	}
	r.comments[comment.ID] = comment
	return nil
}

func (r *MockCommentRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.comments[id]; !exists {
		return errors.New("comment not found")
	}
	delete(r.comments, id)
	return nil
}

// MockLabelRepository provides an in-memory implementation of LabelRepository
type MockLabelRepository struct {
	mu     sync.RWMutex
	labels map[string]*repository.LabelEntity
}

func NewMockLabelRepository() repository.LabelRepository {
	repo := &MockLabelRepository{
		labels: make(map[string]*repository.LabelEntity),
	}

	// デフォルトラベルを追加
	defaultLabels := []*repository.LabelEntity{
		{ID: "1", Name: "bug", Color: "#d73a4a"},
		{ID: "2", Name: "enhancement", Color: "#a2eeef"},
		{ID: "3", Name: "documentation", Color: "#0075ca"},
		{ID: "4", Name: "good first issue", Color: "#7057ff"},
	}

	for _, label := range defaultLabels {
		repo.labels[label.ID] = label
	}

	return repo
}

func (r *MockLabelRepository) List(ctx context.Context) ([]*repository.LabelEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*repository.LabelEntity
	for _, label := range r.labels {
		result = append(result, label)
	}
	return result, nil
}

func (r *MockLabelRepository) Get(ctx context.Context, id string) (*repository.LabelEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	label, exists := r.labels[id]
	if !exists {
		return nil, errors.New("label not found")
	}
	return label, nil
}

func (r *MockLabelRepository) Create(ctx context.Context, label *repository.LabelEntity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.labels[label.ID] = label
	return nil
}

func (r *MockLabelRepository) Update(ctx context.Context, label *repository.LabelEntity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.labels[label.ID]; !exists {
		return errors.New("label not found")
	}
	r.labels[label.ID] = label
	return nil
}

func (r *MockLabelRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.labels[id]; !exists {
		return errors.New("label not found")
	}
	delete(r.labels, id)
	return nil
}

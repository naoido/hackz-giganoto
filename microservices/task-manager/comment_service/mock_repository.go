package comment_service

import (
	"context"
	"errors"
	"sync"

	"object-t.com/hackz-giganoto/microservices/task-manager/repository"
)

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
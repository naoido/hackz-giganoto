package label_service

import (
	"context"
	"errors"
	"sync"

	"object-t.com/hackz-giganoto/microservices/task-manager/repository"
)

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
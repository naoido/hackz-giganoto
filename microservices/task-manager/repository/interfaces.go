package repository

import (
	"context"
	"time"

	commentservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/comment_service"
	labelservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/label_service"
	taskservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/task_service"
)

// TaskEntity represents a task entity in the database
type TaskEntity struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Status      string    `json:"status"`
	AssigneeIds []string  `json:"assignee_ids"`
	LabelIds    []string  `json:"label_ids"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CommentEntity represents a comment entity in the database
type CommentEntity struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	AuthorID  string    `json:"author_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// LabelEntity represents a label entity in the database
type LabelEntity struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// TaskRepository defines the interface for task data operations
type TaskRepository interface {
	// List retrieves tasks with optional filtering
	List(ctx context.Context, status, assigneeID, labelID *string) ([]*TaskEntity, error)

	// Get retrieves a task by ID
	Get(ctx context.Context, id string) (*TaskEntity, error)

	// Create creates a new task
	Create(ctx context.Context, task *TaskEntity) error

	// Update updates an existing task
	Update(ctx context.Context, task *TaskEntity) error

	// Delete deletes a task by ID
	Delete(ctx context.Context, id string) error
}

// CommentRepository defines the interface for comment data operations
type CommentRepository interface {
	// List retrieves comments for a specific task
	List(ctx context.Context, taskID string) ([]*CommentEntity, error)

	// Get retrieves a comment by ID
	Get(ctx context.Context, id string) (*CommentEntity, error)

	// Create creates a new comment
	Create(ctx context.Context, comment *CommentEntity) error

	// Update updates an existing comment
	Update(ctx context.Context, comment *CommentEntity) error

	// Delete deletes a comment by ID
	Delete(ctx context.Context, id string) error
}

// LabelRepository defines the interface for label data operations
type LabelRepository interface {
	// List retrieves all labels
	List(ctx context.Context) ([]*LabelEntity, error)

	// Get retrieves a label by ID
	Get(ctx context.Context, id string) (*LabelEntity, error)

	// Create creates a new label
	Create(ctx context.Context, label *LabelEntity) error

	// Update updates an existing label
	Update(ctx context.Context, label *LabelEntity) error

	// Delete deletes a label by ID
	Delete(ctx context.Context, id string) error
}

// ConvertTaskEntityToResult converts TaskEntity to Goa result type
func ConvertTaskEntityToResult(entity *TaskEntity) *taskservice.Task {
	return &taskservice.Task{
		ID:          entity.ID,
		Title:       entity.Title,
		Description: entity.Description,
		Status:      entity.Status,
		AssigneeIds: entity.AssigneeIds,
		LabelIds:    entity.LabelIds,
		CreatedAt:   entity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   entity.UpdatedAt.Format(time.RFC3339),
	}
}

// ConvertCommentEntityToResult converts CommentEntity to Goa result type
func ConvertCommentEntityToResult(entity *CommentEntity) *commentservice.Comment {
	return &commentservice.Comment{
		ID:        entity.ID,
		TaskID:    entity.TaskID,
		AuthorID:  entity.AuthorID,
		Body:      entity.Body,
		CreatedAt: entity.CreatedAt.Format(time.RFC3339),
	}
}

// ConvertLabelEntityToResult converts LabelEntity to Goa result type
func ConvertLabelEntityToResult(entity *LabelEntity) *labelservice.Label {
	return &labelservice.Label{
		ID:    entity.ID,
		Name:  entity.Name,
		Color: entity.Color,
	}
}

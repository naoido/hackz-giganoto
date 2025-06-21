package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	labelservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/label_service"
	"object-t.com/hackz-giganoto/microservices/task-manager/repository"
)

// labelServiceImpl implements the LabelService interface
type labelServiceImpl struct {
	logger    *log.Logger
	labelRepo repository.LabelRepository
}

// NewLabelService returns the LabelService service implementation
func NewLabelService(logger *log.Logger, labelRepo repository.LabelRepository) labelservice.Service {
	return &labelServiceImpl{
		logger:    logger,
		labelRepo: labelRepo,
	}
}

// List implements the "List" method of the "LabelService" service
func (s *labelServiceImpl) List(ctx context.Context) (res labelservice.LabelCollection, err error) {
	s.logger.Print("labelService.List")

	// リポジトリからラベル一覧を取得
	entities, err := s.labelRepo.List(ctx)
	if err != nil {
		s.logger.Printf("Failed to list labels: %v", err)
		return nil, err
	}

	// エンティティをGoaの結果型に変換
	res = make(labelservice.LabelCollection, len(entities))
	for i, entity := range entities {
		res[i] = repository.ConvertLabelEntityToResult(entity)
	}

	s.logger.Printf("Listed %d labels", len(res))
	return res, nil
}

// Get implements the "Get" method of the "LabelService" service
func (s *labelServiceImpl) Get(ctx context.Context, p *labelservice.GetPayload) (res *labelservice.Label, err error) {
	s.logger.Printf("labelService.Get: id=%s", p.ID)

	// リポジトリからラベルを取得
	entity, err := s.labelRepo.Get(ctx, p.ID)
	if err != nil {
		s.logger.Printf("Failed to get label: %v", err)
		return nil, labelservice.MakeNotFound(err)
	}

	// エンティティをGoaの結果型に変換
	res = repository.ConvertLabelEntityToResult(entity)

	s.logger.Printf("Got label: %s", res.Name)
	return res, nil
}

// Create implements the "Create" method of the "LabelService" service
func (s *labelServiceImpl) Create(ctx context.Context, p *labelservice.CreatePayload) (res *labelservice.Label, err error) {
	s.logger.Printf("labelService.Create: name=%s, color=%s", p.Name, p.Color)

	entity := &repository.LabelEntity{
		ID:    uuid.New().String(),
		Name:  p.Name,
		Color: p.Color,
	}

	// リポジトリにラベルを作成
	err = s.labelRepo.Create(ctx, entity)
	if err != nil {
		s.logger.Printf("Failed to create label: %v", err)
		return nil, err
	}

	// エンティティをGoaの結果型に変換
	res = repository.ConvertLabelEntityToResult(entity)

	s.logger.Printf("Created label: %s", res.ID)
	return res, nil
}

// Update implements the "Update" method of the "LabelService" service
func (s *labelServiceImpl) Update(ctx context.Context, p *labelservice.UpdatePayload) (res *labelservice.Label, err error) {
	s.logger.Printf("labelService.Update: id=%s", p.ID)

	// 既存のラベルを取得
	entity, err := s.labelRepo.Get(ctx, p.ID)
	if err != nil {
		s.logger.Printf("Failed to get label for update: %v", err)
		return nil, labelservice.MakeNotFound(err)
	}

	// 更新フィールドをチェックして適用
	if p.Name != nil {
		entity.Name = *p.Name
	}
	if p.Color != nil {
		entity.Color = *p.Color
	}

	// リポジトリでラベルを更新
	err = s.labelRepo.Update(ctx, entity)
	if err != nil {
		s.logger.Printf("Failed to update label: %v", err)
		return nil, err
	}

	// エンティティをGoaの結果型に変換
	res = repository.ConvertLabelEntityToResult(entity)

	s.logger.Printf("Updated label: %s", res.ID)
	return res, nil
}

// Delete implements the "Delete" method of the "LabelService" service
func (s *labelServiceImpl) Delete(ctx context.Context, p *labelservice.DeletePayload) (err error) {
	s.logger.Printf("labelService.Delete: id=%s", p.ID)

	// ラベルが存在するかチェック
	_, err = s.labelRepo.Get(ctx, p.ID)
	if err != nil {
		s.logger.Printf("Failed to get label for delete: %v", err)
		return labelservice.MakeNotFound(err)
	}

	// リポジトリからラベルを削除
	err = s.labelRepo.Delete(ctx, p.ID)
	if err != nil {
		s.logger.Printf("Failed to delete label: %v", err)
		return err
	}

	s.logger.Printf("Deleted label: %s", p.ID)
	return nil
}

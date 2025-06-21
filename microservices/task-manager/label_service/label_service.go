package label_service

import (
	"context"

	"goa.design/clue/log"
	"github.com/google/uuid"
	labelservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/label_service"
	"object-t.com/hackz-giganoto/microservices/task-manager/repository"
)

// labelServiceImpl implements the LabelService interface
type labelServiceImpl struct {
	logger    func(context.Context, ...any)
	labelRepo repository.LabelRepository
}

// New returns the LabelService service implementation
func New(logger func(context.Context, ...any)) labelservice.Service {
	// モックリポジトリを作成（実際の実装では依存性注入を使用）
	labelRepo := NewMockLabelRepository()
	
	return &labelServiceImpl{
		logger:    logger,
		labelRepo: labelRepo,
	}
}

// List implements the "List" method of the "LabelService" service
func (s *labelServiceImpl) List(ctx context.Context) (res labelservice.LabelCollection, err error) {
	log.Print(ctx, log.KV{K: "method", V: "labelService.List"})
	
	// リポジトリからラベル一覧を取得
	entities, err := s.labelRepo.List(ctx)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to list labels"})
		return nil, err
	}
	
	// エンティティをGoaの結果型に変換
	res = make(labelservice.LabelCollection, len(entities))
	for i, entity := range entities {
		res[i] = repository.ConvertLabelEntityToResult(entity)
	}
	
	log.Print(ctx, log.KV{K: "count", V: len(res)}, log.KV{K: "message", V: "Listed labels"})
	return res, nil
}

// Get implements the "Get" method of the "LabelService" service
func (s *labelServiceImpl) Get(ctx context.Context, p *labelservice.GetPayload) (res *labelservice.Label, err error) {
	log.Print(ctx, log.KV{K: "method", V: "labelService.Get"}, log.KV{K: "id", V: p.ID})
	
	// リポジトリからラベルを取得
	entity, err := s.labelRepo.Get(ctx, p.ID)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to get label"})
		return nil, labelservice.MakeNotFound(err)
	}
	
	// エンティティをGoaの結果型に変換
	res = repository.ConvertLabelEntityToResult(entity)
	
	log.Print(ctx, log.KV{K: "message", V: "Got label"}, log.KV{K: "name", V: res.Name})
	return res, nil
}

// Create implements the "Create" method of the "LabelService" service
func (s *labelServiceImpl) Create(ctx context.Context, p *labelservice.CreatePayload) (res *labelservice.Label, err error) {
	log.Print(ctx, log.KV{K: "method", V: "labelService.Create"}, log.KV{K: "name", V: p.Name}, log.KV{K: "color", V: p.Color})
	
	entity := &repository.LabelEntity{
		ID:    uuid.New().String(),
		Name:  p.Name,
		Color: p.Color,
	}
	
	// リポジトリにラベルを作成
	err = s.labelRepo.Create(ctx, entity)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to create label"})
		return nil, err
	}
	
	// エンティティをGoaの結果型に変換
	res = repository.ConvertLabelEntityToResult(entity)
	
	log.Print(ctx, log.KV{K: "message", V: "Created label"}, log.KV{K: "id", V: res.ID})
	return res, nil
}

// Update implements the "Update" method of the "LabelService" service
func (s *labelServiceImpl) Update(ctx context.Context, p *labelservice.UpdatePayload) (res *labelservice.Label, err error) {
	log.Print(ctx, log.KV{K: "method", V: "labelService.Update"}, log.KV{K: "id", V: p.ID})
	
	// 既存のラベルを取得
	entity, err := s.labelRepo.Get(ctx, p.ID)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to get label for update"})
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
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to update label"})
		return nil, err
	}
	
	// エンティティをGoaの結果型に変換
	res = repository.ConvertLabelEntityToResult(entity)
	
	log.Print(ctx, log.KV{K: "message", V: "Updated label"}, log.KV{K: "id", V: res.ID})
	return res, nil
}

// Delete implements the "Delete" method of the "LabelService" service
func (s *labelServiceImpl) Delete(ctx context.Context, p *labelservice.DeletePayload) (err error) {
	log.Print(ctx, log.KV{K: "method", V: "labelService.Delete"}, log.KV{K: "id", V: p.ID})
	
	// ラベルが存在するかチェック
	_, err = s.labelRepo.Get(ctx, p.ID)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to get label for delete"})
		return labelservice.MakeNotFound(err)
	}
	
	// リポジトリからラベルを削除
	err = s.labelRepo.Delete(ctx, p.ID)
	if err != nil {
		log.Error(ctx, err, log.KV{K: "error", V: "Failed to delete label"})
		return err
	}
	
	log.Print(ctx, log.KV{K: "message", V: "Deleted label"}, log.KV{K: "id", V: p.ID})
	return nil
}
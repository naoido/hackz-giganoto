package wikiapi

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"goa.design/clue/log"
	"goa.design/goa/v3/security"
	wiki "object-t.com/hackz-giganoto/microservices/wiki/gen/wiki"
)

type wikisrvc struct {
	store map[string]*wiki.Wiki
	mu    sync.RWMutex
}

func NewWiki() wiki.Service {
	return &wikisrvc{
		store: make(map[string]*wiki.Wiki),
	}
}

// Wiki取得
func (s *wikisrvc) Show(ctx context.Context, p *wiki.ShowPayload) (res *wiki.Wiki, err error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, wiki.Unauthorized("invalid token")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	log.Printf(ctx, "wiki.show (channel_id: %s, user_id: %s)", p.ID, userID)

	existingWiki, ok := s.store[p.ID]
	if !ok {
		return nil, wiki.NotFound("Wiki not found")
	}
	res = existingWiki
	return res, nil
}

// 新しいWikiを作成
func (s *wikisrvc) Create(ctx context.Context, p *wiki.CreatePayload) (res *wiki.Wiki, err error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, wiki.Unauthorized("invalid token")
	}

	if p.ChannelID == "" || p.Content == "" {
		return nil, wiki.BadRequest("channel_id and content are required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf(ctx, "wiki.create (channel_id: %s, user_id: %s)", p.ChannelID, userID)

	// 新しいWikiのデータを作成する
	currentTime := time.Now()
	newWiki := &wiki.Wiki{
		ID:        "wiki_" + uuid.NewString(),
		ChannelID: p.ChannelID,
		Content:   p.Content,
		CreatedAt: currentTime.Format(time.RFC3339),
		UpdatedAt: currentTime.Format(time.RFC3339),
	}
	s.store[newWiki.ID] = newWiki
	res = newWiki
	return res, nil
}

// 指定されたIDのWikiの内容を更新
func (s *wikisrvc) Update(ctx context.Context, p *wiki.UpdatePayload) (res *wiki.Wiki, err error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return nil, wiki.Unauthorized("invalid token")
	}

	if p.ID == "" || p.Content == "" {
		return nil, wiki.BadRequest("id and content are required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf(ctx, "wiki.update (id: %s, user_id: %s)", p.ID, userID)
	existingWiki, ok := s.store[p.ID]
	if !ok {
		return nil, wiki.NotFound("Wiki not found")
	}
	existingWiki.Content = p.Content
	existingWiki.UpdatedAt = time.Now().Format(time.RFC3339)
	res = existingWiki
	return res, nil
}

// 指定されたIDのWikiを削除
func (s *wikisrvc) Delete(ctx context.Context, p *wiki.DeletePayload) (err error) {
	userID, err := s.getUserIDFromContext(ctx)
	if err != nil {
		return wiki.Unauthorized("invalid token")
	}

	if p.ID == "" {
		return wiki.BadRequest("id is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf(ctx, "wiki.delete (id: %s, user_id: %s)", p.ID, userID)

	_, exists := s.store[p.ID]
	if !exists {
		return wiki.NotFound("Wiki not found")
	}

	delete(s.store, p.ID)
	return nil
}

// JWT認証
func (s *wikisrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	if token == "" {
		return ctx, fmt.Errorf("missing token")
	}

	userID := s.extractUserIDFromToken(token)
	if userID == "" {
		return ctx, fmt.Errorf("invalid token")
	}

	return context.WithValue(ctx, "userID", userID), nil
}

func (s *wikisrvc) getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

func (s *wikisrvc) extractUserIDFromToken(token string) string {
	return "user_" + token[:8]
}

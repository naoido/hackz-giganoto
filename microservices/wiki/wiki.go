package wikiapi

import (
	"context"
	"time"

	"github.com/google/uuid"
	"goa.design/clue/log"
	"goa.design/goa/v3/security"
	"github.com/redis/go-redis/v9"
	wiki "object-t.com/hackz-giganoto/microservices/wiki/gen/wiki"
	security2 "object-t.com/hackz-giganoto/pkg/security"
)

type wikisrvc struct {
	redis *redis.Client
}

func NewWiki(redisClient *redis.Client) wiki.Service {
    return &wikisrvc{
        redis: redisClient,
    }
}

// JWTAuth implements the authorization logic for the JWT security scheme.
func (s *wikisrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	claims, err := security2.ValidToken(token)
	if err != nil {
		return ctx, wiki.Unauthorized("invalid token")
	}
	return security2.HasPermission(ctx, claims, scheme)
}

// Wiki取得
func (s *wikisrvc) Show(ctx context.Context, p *wiki.ShowPayload) (res *wiki.Wiki, err error) {
	authInfo := security2.ContextAuthInfo(ctx)
	userID, ok := authInfo.Claims["sub"].(string)
	if !ok {
		return nil, wiki.Unauthorized("sub not a string")
	}

	log.Printf(ctx, "wiki.show (channel_id: %s, user_id: %s)", p.ID, userID)

	_, err = s.redis.Get(ctx, p.ID).Result()
	if err == nil {
		return nil, wiki.NotFound("Wiki not found")
	}
			
	res = &wiki.Wiki{
		ID:        p.ID,
	}
	return res, nil
}

// 新しいWikiを作成
func (s *wikisrvc) Create(ctx context.Context, p *wiki.CreatePayload) (res *wiki.Wiki, err error) {
	authInfo := security2.ContextAuthInfo(ctx)
	userID, ok := authInfo.Claims["sub"].(string)
	if !ok {
		return nil, wiki.Unauthorized("sub not a string")
	}

	uniqueID := uuid.NewString()

	_, err = s.redis.Set(ctx, p.ChannelID, p.Content, 10*time.Minute).Result()
	_, err = s.redis.Set(ctx, uniqueID, p.Content, 10*time.Minute).Result()

	log.Printf(ctx, "wiki.create (channel_id: %s, user_id: %s)", p.ChannelID, userID)

	
	res = &wiki.Wiki{
		ID:        uniqueID,
		ChannelID: p.ChannelID,
		Content:   p.Content,
	}
	return res, nil
}

// 指定されたIDのWikiの内容を更新
func (s *wikisrvc) Update(ctx context.Context, p *wiki.UpdatePayload) (res *wiki.Wiki, err error) {
	authInfo := security2.ContextAuthInfo(ctx)
	userID, ok := authInfo.Claims["sub"].(string)
	if !ok {
		return nil, wiki.Unauthorized("sub not a string")
	}
	
	_, err = s.redis.Get(ctx, p.ID).Result()
	if err != nil {
		return nil, wiki.NotFound("Wiki not found")
	}

	log.Printf(ctx, "wiki.update (id: %s, user_id: %s)", p.ID, userID)
	
	_, err = s.redis.Set(ctx, p.ID, p.Content, 10*time.Minute).Result()
	if err == nil {
		return nil, wiki.NotFound("Wiki not found")
	}
	
	res = &wiki.Wiki{
		ID:        p.ID,
		Content:  p.Content,
	}
	return res, nil
}

// 指定されたIDのWikiを削除
func (s *wikisrvc) Delete(ctx context.Context, p *wiki.DeletePayload) (err error) {
	authInfo := security2.ContextAuthInfo(ctx)
	userID, ok := authInfo.Claims["sub"].(string)
	if !ok {
		return nil
	}

	log.Printf(ctx, "wiki.delete (id: %s, user_id: %s)", p.ID, userID)

	_, err = s.redis.Del(ctx, p.ID).Result()
	if err == nil {
		return wiki.NotFound("Wiki not found")
	}
	return nil
}

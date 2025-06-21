package profileapi

import (
	"context"
	"github.com/redis/go-redis/v9"
	"goa.design/goa/v3/security"
	security2 "object-t.com/hackz-giganoto/pkg/security" 
	profile "object-t.com/hackz-giganoto/microservices/profile/gen/profile"
)

// profile service example implementation.
// The example methods log the requests and return zero values.
type profilesrvc struct {
	redis *redis.Client
}

func NewProfile(rdb *redis.Client) profile.Service {
	return &profilesrvc{redis: rdb}
}

// JWTAuth implements the authorization logic for the JWT security scheme.
func (s *profilesrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	claims, err := security2.ValidToken(token)
	if err != nil {
		return ctx, profile.InvalidToken("invalid token")
	}
	return security2.HasPermission(ctx, claims, scheme)
}

func (s *profilesrvc) CreateProfile(ctx context.Context, p *profile.CreateProfilePayload) (*profile.CreateProfileResult, error) {
	authInfo := security2.ContextAuthInfo(ctx)
	userID, ok := authInfo.Claims["sub"].(string)
	if !ok {
		return nil, profile.InvalidToken("sub not a string")
	}

	_, err := s.redis.Get(ctx, userID).Result()
	if err == redis.Nil {
		// 存在しなかった場合、新規登録
		if err := s.redis.Set(ctx, userID, p.Name, 0).Err(); err != nil {
			return nil, profile.InternalError("redis set failed")
		}
	} else if err != nil {
		return nil, profile.InternalError("redis get failed")
	}

	return &profile.CreateProfileResult{
		UserID: userID,
		Name:   p.Name,
	}, nil
}

func (s *profilesrvc) GetProfile(ctx context.Context, p *profile.GetProfilePayload) (*profile.GetProfileResult, error) {
	authInfo := security2.ContextAuthInfo(ctx)
	userID, ok := authInfo.Claims["sub"].(string)
	if !ok {
		return nil, profile.InvalidToken("sub not a string")
	}

	name, err := s.redis.Get(ctx, userID).Result()
	if err == redis.Nil {
		return nil, profile.NotFound("profile not found")
	}
	if err != nil {
		return nil, profile.InternalError("redis get failed")
	}

	return &profile.GetProfileResult{
		UserID: userID,
		Name:   name,
	}, nil
}

func (s *profilesrvc) UpdateProfile(ctx context.Context, p *profile.UpdateProfilePayload) (*profile.UpdateProfileResult, error) {
	authInfo := security2.ContextAuthInfo(ctx)
	userID, ok := authInfo.Claims["sub"].(string)
	if !ok {
		return nil, profile.InvalidToken("sub not a string")
	}

	if err := s.redis.Set(ctx, userID, p.Name, 0).Err(); err != nil {
		return nil, profile.InternalError("failed to update profile")
	}

	return &profile.UpdateProfileResult{
		UserID: userID,
		Name:   p.Name,
	}, nil
}

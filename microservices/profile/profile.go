package profileapi

import (
	"context"
	"github.com/redis/go-redis/v9"
	"object-t.com/hackz-giganoto/microservices/bff/gen/bff"
	security2 "object-t.com/hackz-giganoto/pkg/security"

	"goa.design/clue/log"
	"goa.design/goa/v3/security"
	profile "object-t.com/hackz-giganoto/microservices/profile/gen/profile"
)

// profile service example implementation.
// The example methods log the requests and return zero values.
type profilesrvc struct {
	redis *redis.Client
}

// NewProfile returns the profile service implementation.
func NewProfile(redis *redis.Client) profile.Service {
	return &profilesrvc{
		redis: redis,
	}
}

type Profile struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

// JWTAuth implements the authorization logic for service "profile" for the
// "jwt" security scheme.
func (s *profilesrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	claims, err := security2.ValidToken(token)
	if err != nil {
		log.Printf(ctx, "invalid token: %v", err)
		return ctx, bff.InvalidArgument("invalid token")
	}
	return security2.HasPermission(ctx, claims, scheme)
}

// Get user profile
func (s *profilesrvc) GetProfile(ctx context.Context, p *profile.GetProfilePayload) (res *profile.GetProfileResult, err error) {
	log.Printf(ctx, "profile.get_profile")

	profileKey := "profile:" + p.UserID
	userName := s.redis.Get(ctx, profileKey).Val()

	if userName == "" {
		return nil, profile.NotFound("user not found")
	}

	res = &profile.GetProfileResult{
		UserID: p.UserID,
		Name:   userName,
	}

	return
}

// Update user profile
func (s *profilesrvc) UpdateProfile(ctx context.Context, p *profile.UpdateProfilePayload) (res *profile.UpdateProfileResult, err error) {
	log.Printf(ctx, "profile.update_profile")
	userId := security2.ContextAuthInfo(ctx).Claims["user_id"].(string)
	profileKey := "profile:" + userId

	if s.redis.Set(ctx, profileKey, p.Name, 0).Err() != nil {
		return nil, profile.InternalError("Internal error")
	}

	return &profile.UpdateProfileResult{
		UserID: userId,
		Name:   p.Name,
	}, nil
}

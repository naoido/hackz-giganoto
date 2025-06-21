package profileapi

import (
	"context"

	"goa.design/clue/log"
	"goa.design/goa/v3/security"
	"gorm.io/gorm"
	security2 "object-t.com/hackz-giganoto/pkg/security"
	model "object-t.com/hackz-giganoto/microservices/model"
	profile "object-t.com/hackz-giganoto/microservices/profile/gen/profile"
)

// profile service example implementation.
// The example methods log the requests and return zero values.
type profilesrvc struct {
	db *gorm.DB
}

// NewProfile returns the profile service implementation.
func NewProfile(db *gorm.DB) profile.Service {
	return &profilesrvc{db: db}
}

// JWTAuth implements the authorization logic for service "profile" for the
// "jwt" security scheme.
func (s *profilesrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	claims, err := security2.ValidToken(token)
	if err != nil {
		log.Printf(ctx, "invalid token: %v", err)
		return ctx, profile.InvalidToken("invalid token")
	}
	return security2.HasPermission(ctx, claims, scheme)
}

// Create a new user profile
func (s *profilesrvc) CreateProfile(ctx context.Context, p *profile.CreateProfilePayload) (res *profile.CreateProfileResult, err error) {
	res = &profile.CreateProfileResult{}
	log.Printf(ctx, "profile.create_profile")

	var profileModel model.Profile

	profileModel.UserID = p.UserID
	profileModel.Name = p.Name

	if err := s.db.Create(profileModel).Error; err != nil {
		log.Printf(ctx, "failed to create profile: %v", err)
		return nil, profile.InternalError("failed to create profile")
	}

	res = &profile.CreateProfileResult{
		UserID:    profileModel.UserID,
		Name:      profileModel.Name,
	}

	log.Printf(ctx, "profile created: %s", res.UserID)
	return res, nil
}

// Get user profile
func (s *profilesrvc) GetProfile(ctx context.Context, p *profile.GetProfilePayload) (res *profile.GetProfileResult, err error) {
	res = &profile.GetProfileResult{}
	log.Printf(ctx, "profile.get_profile")

	var profileModel model.Profile

	authInfo := security2.ContextAuthInfo(ctx)
	userID, ok := authInfo.Claims["sub"].(string)
	if !ok {
		log.Printf(ctx, "invalid token: sub claim not found or not a string")
		return nil, profile.InvalidToken("invalid token")
	}

	if err := s.db.First(&profileModel, "user_id = ?", userID).Error; err != nil {
		log.Printf(ctx, "failed to get profile: %v", err)
		return nil, profile.NotFound("profile not found")
	}

	res = &profile.GetProfileResult{
		UserID: profileModel.UserID,
		Name:   profileModel.Name,
	}

	return res, nil
}

// Update user profile
func (s *profilesrvc) UpdateProfile(ctx context.Context, p *profile.UpdateProfilePayload) (res *profile.UpdateProfileResult, err error) {
	res = &profile.UpdateProfileResult{}
	log.Printf(ctx, "profile.update_profile")

	var profileModel model.Profile

	authInfo := security2.ContextAuthInfo(ctx)
	userID, ok := authInfo.Claims["sub"].(string)
	if !ok {
		log.Printf(ctx, "invalid token: sub claim not found or not a string")
		return nil, profile.InvalidToken("invalid token")
	}

	if err := s.db.First(&profileModel, "user_id = ?", userID).Error; err != nil {
		log.Printf(ctx, "failed to get profile for update: %v", err)
		return nil, profile.NotFound("profile not found")
	}

	if err := s.db.Model(&profileModel).Updates(model.Profile{Name: p.Name}).Error; err != nil {
		log.Printf(ctx, "failed to update profile: %v", err)
		return nil, profile.InternalError("failed to update profile")
	}

	res = &profile.UpdateProfileResult{
		UserID:    profileModel.UserID,
		Name:      profileModel.Name,
	}

	return res, nil
}

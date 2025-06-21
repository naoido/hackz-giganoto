package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"goa.design/goa/v3/security"

	profileservice "object-t.com/hackz-giganoto/microservices/chat/gen/profile"
)

type ProfileServiceImpl struct {
	profiles map[string]*ProfileData
	mu       sync.RWMutex
}

type ProfileData struct {
	UserID    string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProfileService() *ProfileServiceImpl {
	return &ProfileServiceImpl{
		profiles: make(map[string]*ProfileData),
	}
}

func (p *ProfileServiceImpl) GetProfile(ctx context.Context, payload *profileservice.GetProfilePayload) (*profileservice.GetProfileResult, error) {
	userID, err := p.getUserIDFromContext(ctx)
	if err != nil {
		return nil, profileservice.Unauthorized("invalid token")
	}

	p.mu.RLock()
	profile, exists := p.profiles[userID]
	p.mu.RUnlock()

	if !exists {
		return nil, profileservice.NotFound("profile not found")
	}

	createdAt := profile.CreatedAt.Format(time.RFC3339)
	updatedAt := profile.UpdatedAt.Format(time.RFC3339)

	return &profileservice.GetProfileResult{
		UserID:    profile.UserID,
		Name:      profile.Name,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}, nil
}

func (p *ProfileServiceImpl) UpdateProfile(ctx context.Context, payload *profileservice.UpdateProfilePayload) (*profileservice.UpdateProfileResult, error) {
	userID, err := p.getUserIDFromContext(ctx)
	if err != nil {
		return nil, profileservice.Unauthorized("invalid token")
	}

	if payload.Name == "" {
		return nil, profileservice.BadRequest("name is required")
	}

	now := time.Now()

	p.mu.Lock()
	profile, exists := p.profiles[userID]
	if !exists {
		profile = &ProfileData{
			UserID:    userID,
			CreatedAt: now,
		}
		p.profiles[userID] = profile
	}
	profile.Name = payload.Name
	profile.UpdatedAt = now
	p.mu.Unlock()

	createdAt := profile.CreatedAt.Format(time.RFC3339)
	updatedAt := profile.UpdatedAt.Format(time.RFC3339)

	return &profileservice.UpdateProfileResult{
		UserID:    profile.UserID,
		Name:      profile.Name,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}, nil
}

func (p *ProfileServiceImpl) GetProfileByID(ctx context.Context, payload *profileservice.GetProfileByIDPayload) (*profileservice.GetProfileByIDResult, error) {
	_, err := p.getUserIDFromContext(ctx)
	if err != nil {
		return nil, profileservice.Unauthorized("invalid token")
	}

	if payload.UserID == "" {
		return nil, profileservice.BadRequest("user_id is required")
	}

	p.mu.RLock()
	profile, exists := p.profiles[payload.UserID]
	p.mu.RUnlock()

	if !exists {
		return nil, profileservice.NotFound("profile not found")
	}

	createdAt := profile.CreatedAt.Format(time.RFC3339)
	updatedAt := profile.UpdatedAt.Format(time.RFC3339)

	return &profileservice.GetProfileByIDResult{
		UserID:            profile.UserID,
		Name:              profile.Name,
		CreationTimestamp: &createdAt,
		UpdatedAt:         &updatedAt,
	}, nil
}

func (p *ProfileServiceImpl) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	if token == "" {
		return ctx, fmt.Errorf("missing token")
	}
	
	userID := p.extractUserIDFromToken(token)
	if userID == "" {
		return ctx, fmt.Errorf("invalid token")
	}
	
	return context.WithValue(ctx, "userID", userID), nil
}

func (p *ProfileServiceImpl) getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

func (p *ProfileServiceImpl) extractUserIDFromToken(token string) string {
	return "user_" + token[:8]
}
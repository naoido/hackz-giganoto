package profileapi

import (
	"context"

	"goa.design/clue/log"
	profile "object-t.com/hackz-giganoto/microservices/profile/gen/profile"
)

// profile service example implementation.
// The example methods log the requests and return zero values.
type profilesrvc struct{}

// NewProfile returns the profile service implementation.
func NewProfile() profile.Service {
	return &profilesrvc{}
}

// Create implements create.
func (s *profilesrvc) Create(ctx context.Context, p *profile.ProfilePayload) (res *profile.GoaExampleProfile, err error) {
	res = &profile.GoaExampleProfile{}
	log.Printf(ctx, "profile.create")
	return
}

// Get implements get.
func (s *profilesrvc) Get(ctx context.Context, p *profile.GetPayload) (res *profile.GoaExampleProfile, err error) {
	res = &profile.GoaExampleProfile{}
	log.Printf(ctx, "profile.get")
	return
}

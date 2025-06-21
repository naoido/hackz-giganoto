package authapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"goa.design/clue/log"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	auth "object-t.com/hackz-giganoto/microservices/auth/gen/auth"
)

// GitHub user profile response
type GitHubUser struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// OAuth state entry
type stateEntry struct {
	CreatedAt time.Time
	ExpiresAt time.Time
}

// Token storage entry
type tokenEntry struct {
	UserID    string
	Login     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// auth service implementation with GitHub OAuth
type authsrvc struct {
	oauthConfig *oauth2.Config
	jwtSecret   []byte
	states      map[string]stateEntry
	tokens      map[string]tokenEntry
	statesMutex sync.RWMutex
	tokensMutex sync.RWMutex
}

// NewAuth returns the auth service implementation.
func NewAuth() auth.Service {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("secret") // Default for development
	}

	return &authsrvc{
		oauthConfig: &oauth2.Config{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     github.Endpoint,
		},
		jwtSecret: jwtSecret,
		states:    make(map[string]stateEntry),
		tokens:    make(map[string]tokenEntry),
	}
}

// Generate random state for OAuth CSRF protection
func (s *authsrvc) generateState() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Clean expired states and tokens
func (s *authsrvc) cleanExpired() {
	now := time.Now()

	s.statesMutex.Lock()
	for state, entry := range s.states {
		if now.After(entry.ExpiresAt) {
			delete(s.states, state)
		}
	}
	s.statesMutex.Unlock()

	s.tokensMutex.Lock()
	for token, entry := range s.tokens {
		if now.After(entry.ExpiresAt) {
			delete(s.tokens, token)
		}
	}
	s.tokensMutex.Unlock()
}

// Get GitHub OAuth authorization URL with state parameter
func (s *authsrvc) AuthURL(ctx context.Context) (res *auth.AuthURLResult, err error) {
	s.cleanExpired()

	state := s.generateState()

	// Store state with 10 minute expiration
	s.statesMutex.Lock()
	s.states[state] = stateEntry{
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	s.statesMutex.Unlock()

	authURL := s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	res = &auth.AuthURLResult{
		AuthURL: authURL,
		State:   state,
	}

	log.Info(ctx, log.KV{"auth.auth_url", fmt.Sprintf("generated state:  %s", state)})
	return
}

// Fetch GitHub user profile
func (s *authsrvc) fetchGitHubUser(ctx context.Context, token *oauth2.Token) (*GitHubUser, error) {
	client := s.oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", string(body))
	}

	var user GitHubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &user, nil
}

// Handle GitHub OAuth callback and return opaque token
func (s *authsrvc) OauthCallback(ctx context.Context, p *auth.OauthCallbackPayload) (res *auth.OauthCallbackResult, err error) {
	s.cleanExpired()

	// Validate state
	s.statesMutex.RLock()
	stateEntry, exists := s.states[p.State]
	s.statesMutex.RUnlock()

	if !exists || time.Now().After(stateEntry.ExpiresAt) {
		log.Print(ctx, log.KV{"auth.oauth_callback", fmt.Sprintf("ERROR: invalid state: %s", p.State)})
		return nil, auth.InvalidState("Invalid or expired state parameter")
	}

	// Remove used state
	s.statesMutex.Lock()
	delete(s.states, p.State)
	s.statesMutex.Unlock()

	// Exchange code for token
	token, err := s.oauthConfig.Exchange(ctx, p.Code)
	if err != nil {
		log.Print(ctx, log.KV{"auth.oauth_callback", "ERROR: failed to exchange code"}, log.KV{"error", err.Error()})
		return nil, auth.InvalidCode("Invalid authorization code")
	}

	// Fetch GitHub user profile
	githubUser, err := s.fetchGitHubUser(ctx, token)
	if err != nil {
		log.Print(ctx, log.KV{"auth.oauth_callback", "ERROR: failed to fetch GitHub user"}, log.KV{"error", err.Error()})
		return nil, auth.GithubError("Failed to fetch GitHub user profile")
	}

	// Generate opaque token
	opaqueToken := uuid.New().String()

	// Store token with 1 hour expiration
	s.tokensMutex.Lock()
	s.tokens[opaqueToken] = tokenEntry{
		UserID:    fmt.Sprintf("%d", githubUser.ID),
		Login:     githubUser.Login,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Minute),
	}
	s.tokensMutex.Unlock()

	res = &auth.OauthCallbackResult{
		AccessToken: opaqueToken,
		TokenType:   "Bearer",
		ExpiresIn:   3600, // 1 hour
		UserID:      fmt.Sprintf("%d", githubUser.ID),
	}

	log.Info(ctx, log.KV{"auth.oauth_callback", fmt.Sprintf("generated token for user %s (ID: %d)", githubUser.Login, githubUser.ID)})
	return
}

// Introspect opaque token and return internal JWT token for Kong Gateway
func (s *authsrvc) Introspect(ctx context.Context, p *auth.IntrospectPayload) (res *auth.IntrospectResult, err error) {
	s.cleanExpired()

	// Validate opaque token
	s.tokensMutex.RLock()
	tokenEntry, exists := s.tokens[p.Token]
	s.tokensMutex.RUnlock()

	if !exists || time.Now().After(tokenEntry.ExpiresAt) {
		log.Print(ctx, log.KV{"auth.introspect", "ERROR: invalid token"})
		return nil, auth.InvalidToken("Token is invalid or expired")
	}

	// Create JWT claims
	claims := jwt.MapClaims{
		"sub":    tokenEntry.UserID,
		"login":  tokenEntry.Login,
		"scopes": []string{"api:read", "api:write"},
		"iat":    time.Now().Unix(),
		"exp":    tokenEntry.ExpiresAt.Unix(),
	}

	log.Print(ctx, log.KV{"auth.introspect", "DEBUG: creating JWT for token validation"})

	// Create and sign JWT
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, err := jwtToken.SignedString(s.jwtSecret)
	if err != nil {
		log.Print(ctx, log.KV{"auth.introspect", "ERROR: failed to sign JWT"}, log.KV{"error", err.Error()})
		return nil, auth.InternalError("Failed to generate internal token")
	}

	exp := tokenEntry.ExpiresAt.Unix()
	res = &auth.IntrospectResult{
		JWT:    jwtString,
		Active: true,
		Exp:    &exp,
		Scopes: []string{"api:read", "api:write"},
	}

	log.Info(ctx, log.KV{"auth.introspect", fmt.Sprintf("validated token for user %s", tokenEntry.Login)})
	return
}

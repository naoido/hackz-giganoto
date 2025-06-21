package authapi

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	"crypto/rand"

	"github.com/redis/go-redis/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"goa.design/clue/log"
	auth "object-t.com/hackz-giganoto/microservices/auth/gen/auth"
)

// auth service example implementation.
// The example methods log the requests and return zero values.
type authsrvc struct {
	redis *redis.Client
}

// NewAuth returns the auth service implementation.
func NewAuth(redis *redis.Client) auth.Service {
	return &authsrvc{redis: redis}
}

// Introspect opaque token and return internal JWT token for Kong Gateway
func (s *authsrvc) Introspect(ctx context.Context, p *auth.IntrospectPayload) (res *auth.IntrospectResult, err error) {
	res = &auth.IntrospectResult{}
	log.Printf(ctx, "auth.introspect")

	userID, err := s.redis.Get(ctx, "token:"+p.Token).Result()
	if err == redis.Nil {
		log.Printf(ctx, "token not found: %s", p.Token)
		return nil, auth.InvalidToken("token not found")
	}

	claims := jwt.MapClaims{
		"sub":   userID,
		"scope": []string{"api:read", "api:write"},
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Minute).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte("secret"))
	if err != nil {
		return nil, auth.InternalError("failed to sign token")
	}

	res.JWT = token
	return res, nil
}

// Get GitHub OAuth authorization URL with state parameter
func (s *authsrvc) AuthURL(ctx context.Context) (res *auth.AuthURLResult, err error) {
	res = &auth.AuthURLResult{}
	log.Printf(ctx, "auth.auth_url")

	// Generate a random state string
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	state := hex.EncodeToString(b)

	// Store the state in Redis with an expiration
	err = s.redis.Set(ctx, "state:"+state, "true", 10*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	// Construct the GitHub OAuth URL
	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	redirectURI := os.Getenv("GITHUB_REDIRECT_URL")
	authURL := "https://github.com/login/oauth/authorize?client_id=" + githubClientID +
		"&redirect_uri=" + url.QueryEscape(redirectURI) +
		"&state=" + state

	res.AuthURL = authURL
	res.State = state

	return
}

// Handle GitHub OAuth callback and return opaque token
func (s *authsrvc) OauthCallback(ctx context.Context, p *auth.OauthCallbackPayload) (res *auth.OauthCallbackResult, err error) {
	res = &auth.OauthCallbackResult{}
	log.Printf(ctx, "auth.oauth_callback, state: %s", p.State)

	// 1. Validate state
	stateKey := "state:" + p.State
	val, err := s.redis.Get(ctx, stateKey).Result()
	if err == redis.Nil {
		log.Printf(ctx, "invalid or expired state: %s", p.State)
		return nil, auth.InvalidState("invalid or expired state")
	}
	if err != nil {
		log.Errorf(ctx, err, "failed to get state from redis")
		return nil, auth.InternalError("failed to get state from redis")
	}
	if val != "true" {
		log.Printf(ctx, "invalid state value: %s", val)
		return nil, auth.InvalidState("invalid state value")
	}
	// Delete state after use
	if err := s.redis.Del(ctx, stateKey).Err(); err != nil {
		// Log warning but continue
		log.Warnf(ctx, "failed to delete state from redis: %v", err)
	}

	// 2. Exchange code for access token
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	redirectURI := os.Getenv("GITHUB_REDIRECT_URL")

	tokenURL := "https://github.com/login/oauth/access_token"
	reqBody, _ := json.Marshal(map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          p.Code,
		"redirect_uri":  redirectURI,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Errorf(ctx, err, "failed to create request to github")
		return nil, auth.InternalError("failed to create request to github")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf(ctx, err, "failed to request access token from github")
		return nil, auth.GithubError("failed to request access token from github")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Errorf(ctx, nil, "github returned non-200 status for access token: %s", string(bodyBytes))
		return nil, auth.InvalidCode("github returned non-200 status for access token")
	}

	var tokenRes struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		log.Errorf(ctx, err, "failed to decode access token response from github")
		return nil, auth.InternalError("failed to decode access token response")
	}
	if tokenRes.AccessToken == "" {
		log.Errorf(ctx, nil, "access token is empty in github response")
		return nil, auth.InvalidCode("access token is empty in github response")
	}

	// 3. Get user info
	userURL := "https://api.github.com/user"
	req, err = http.NewRequestWithContext(ctx, "GET", userURL, nil)
	if err != nil {
		log.Errorf(ctx, err, "failed to create request for user info")
		return nil, auth.InternalError("failed to create request for user info")
	}
	req.Header.Set("Authorization", "Bearer "+tokenRes.AccessToken)

	resp, err = client.Do(req)
	if err != nil {
		log.Errorf(ctx, err, "failed to request user info from github")
		return nil, auth.GithubError("failed to request user info from github")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Errorf(ctx, nil, "github returned non-200 status for user info: %s", string(bodyBytes))
		return nil, auth.GithubError("github returned non-200 status for user info")
	}

	var userRes struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userRes); err != nil {
		log.Errorf(ctx, err, "failed to decode user info response from github")
		return nil, auth.InternalError("failed to decode user info response")
	}

	// 4. Create Opaque Token and store info in Redis
	opaqueToken := uuid.NewString()

	expiration := 1 * time.Hour
	tokenKey := "token:" + opaqueToken
	if err := s.redis.Set(ctx, tokenKey, userRes.ID, expiration).Err(); err != nil {
		log.Errorf(ctx, err, "failed to store opaque token in redis")
		return nil, auth.InternalError("failed to store opaque token in redis")
	}

	// 5. Return result
	res.AccessToken = opaqueToken
	res.TokenType = "Bearer"
	res.ExpiresIn = int64(expiration.Seconds())
	res.UserID = strconv.FormatInt(userRes.ID, 10)

	log.Printf(ctx, "successfully issued opaque token for user %s", res.UserID)

	return res, nil
}

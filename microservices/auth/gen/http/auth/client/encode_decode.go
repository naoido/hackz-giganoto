// Code generated by goa v3.21.1, DO NOT EDIT.
//
// auth HTTP client encoders and decoders
//
// Command:
// $ goa gen object-t.com/hackz-giganoto/microservices/auth/design

package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	goahttp "goa.design/goa/v3/http"
	auth "object-t.com/hackz-giganoto/microservices/auth/gen/auth"
)

// BuildIntrospectRequest instantiates a HTTP request object with method and
// path set to call the "auth" service "introspect" endpoint
func (c *Client) BuildIntrospectRequest(ctx context.Context, v any) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: IntrospectAuthPath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("auth", "introspect", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeIntrospectRequest returns an encoder for requests sent to the auth
// introspect server.
func EncodeIntrospectRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, any) error {
	return func(req *http.Request, v any) error {
		p, ok := v.(*auth.IntrospectPayload)
		if !ok {
			return goahttp.ErrInvalidType("auth", "introspect", "*auth.IntrospectPayload", v)
		}
		body := NewIntrospectRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("auth", "introspect", err)
		}
		return nil
	}
}

// DecodeIntrospectResponse returns a decoder for responses returned by the
// auth introspect endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeIntrospectResponse may return the following errors:
//   - "internal_error" (type auth.InternalError): http.StatusInternalServerError
//   - "invalid_token" (type auth.InvalidToken): http.StatusUnauthorized
//   - error: internal error
func DecodeIntrospectResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (any, error) {
	return func(resp *http.Response) (any, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body IntrospectResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("auth", "introspect", err)
			}
			err = ValidateIntrospectResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("auth", "introspect", err)
			}
			res := NewIntrospectResultOK(&body)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body string
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("auth", "introspect", err)
			}
			return nil, NewIntrospectInternalError(body)
		case http.StatusUnauthorized:
			var (
				body string
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("auth", "introspect", err)
			}
			return nil, NewIntrospectInvalidToken(body)
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("auth", "introspect", resp.StatusCode, string(body))
		}
	}
}

// BuildAuthURLRequest instantiates a HTTP request object with method and path
// set to call the "auth" service "auth_url" endpoint
func (c *Client) BuildAuthURLRequest(ctx context.Context, v any) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: AuthURLAuthPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("auth", "auth_url", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// DecodeAuthURLResponse returns a decoder for responses returned by the auth
// auth_url endpoint. restoreBody controls whether the response body should be
// restored after having been read.
// DecodeAuthURLResponse may return the following errors:
//   - "internal_error" (type auth.InternalError): http.StatusInternalServerError
//   - error: internal error
func DecodeAuthURLResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (any, error) {
	return func(resp *http.Response) (any, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body AuthURLResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("auth", "auth_url", err)
			}
			err = ValidateAuthURLResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("auth", "auth_url", err)
			}
			res := NewAuthURLResultOK(&body)
			return res, nil
		case http.StatusInternalServerError:
			var (
				body string
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("auth", "auth_url", err)
			}
			return nil, NewAuthURLInternalError(body)
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("auth", "auth_url", resp.StatusCode, string(body))
		}
	}
}

// BuildOauthCallbackRequest instantiates a HTTP request object with method and
// path set to call the "auth" service "oauth_callback" endpoint
func (c *Client) BuildOauthCallbackRequest(ctx context.Context, v any) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: OauthCallbackAuthPath()}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("auth", "oauth_callback", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeOauthCallbackRequest returns an encoder for requests sent to the auth
// oauth_callback server.
func EncodeOauthCallbackRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, any) error {
	return func(req *http.Request, v any) error {
		p, ok := v.(*auth.OauthCallbackPayload)
		if !ok {
			return goahttp.ErrInvalidType("auth", "oauth_callback", "*auth.OauthCallbackPayload", v)
		}
		values := req.URL.Query()
		values.Add("code", p.Code)
		values.Add("state", p.State)
		req.URL.RawQuery = values.Encode()
		return nil
	}
}

// DecodeOauthCallbackResponse returns a decoder for responses returned by the
// auth oauth_callback endpoint. restoreBody controls whether the response body
// should be restored after having been read.
// DecodeOauthCallbackResponse may return the following errors:
//   - "github_error" (type auth.GithubError): http.StatusBadGateway
//   - "internal_error" (type auth.InternalError): http.StatusInternalServerError
//   - "invalid_code" (type auth.InvalidCode): http.StatusBadRequest
//   - "invalid_state" (type auth.InvalidState): http.StatusBadRequest
//   - error: internal error
func DecodeOauthCallbackResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (any, error) {
	return func(resp *http.Response) (any, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body OauthCallbackResponseBody
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("auth", "oauth_callback", err)
			}
			err = ValidateOauthCallbackResponseBody(&body)
			if err != nil {
				return nil, goahttp.ErrValidationError("auth", "oauth_callback", err)
			}
			res := NewOauthCallbackResultOK(&body)
			return res, nil
		case http.StatusBadGateway:
			var (
				body string
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("auth", "oauth_callback", err)
			}
			return nil, NewOauthCallbackGithubError(body)
		case http.StatusInternalServerError:
			var (
				body string
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("auth", "oauth_callback", err)
			}
			return nil, NewOauthCallbackInternalError(body)
		case http.StatusBadRequest:
			en := resp.Header.Get("goa-error")
			switch en {
			case "invalid_code":
				var (
					body string
					err  error
				)
				err = decoder(resp).Decode(&body)
				if err != nil {
					return nil, goahttp.ErrDecodingError("auth", "oauth_callback", err)
				}
				return nil, NewOauthCallbackInvalidCode(body)
			case "invalid_state":
				var (
					body string
					err  error
				)
				err = decoder(resp).Decode(&body)
				if err != nil {
					return nil, goahttp.ErrDecodingError("auth", "oauth_callback", err)
				}
				return nil, NewOauthCallbackInvalidState(body)
			default:
				body, _ := io.ReadAll(resp.Body)
				return nil, goahttp.ErrInvalidResponse("auth", "oauth_callback", resp.StatusCode, string(body))
			}
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("auth", "oauth_callback", resp.StatusCode, string(body))
		}
	}
}

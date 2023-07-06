package transport

import (
	"context"
	"net/http"
)

type (
	TokenProvider func(ctx context.Context) (string, error)

	Logger interface {
		Log(req *http.Request, res *http.Response) error
	}

	UserAgentRoundTripper struct {
		userAgent string
		next      http.RoundTripper
	}

	AuthorizationRoundTripper struct {
		provider      TokenProvider
		authorization string
		next          http.RoundTripper
	}
)

func NewUserAgentRoundTripper(userAgent string, next http.RoundTripper) http.RoundTripper {
	return &UserAgentRoundTripper{
		userAgent: userAgent,
		next:      next,
	}
}

func NewAuthorizationRoundTripper(provider TokenProvider, next http.RoundTripper) http.RoundTripper {
	return &AuthorizationRoundTripper{
		provider: provider,
		next:     next,
	}
}

func (a UserAgentRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	defer closeBody(r)

	newRequest := cloneRequest(r)
	newRequest.Header.Add("User-Agent", a.userAgent)
	return a.next.RoundTrip(newRequest)
}

func (p *AuthorizationRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	defer closeBody(r)

	if p.authorization == "" {
		token, err := p.provider(r.Context())
		if err != nil {
			return p.next.RoundTrip(r)
		}

		p.authorization = token
	}

	newRequest := cloneRequest(r)
	newRequest.Header.Add("Authorization", p.authorization)
	return p.next.RoundTrip(newRequest)
}

func cloneRequest(request *http.Request) *http.Request {
	newRequest := new(http.Request)
	*newRequest = *request

	newRequest.Header = make(http.Header, len(request.Header))
	for k, v := range request.Header {
		newRequest.Header[k] = append([]string(nil), v...)
	}

	return newRequest
}

func closeBody(r *http.Request) {
	if r.Body != nil {
		_ = r.Body.Close()
	}
}

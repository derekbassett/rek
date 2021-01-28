package rek

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Option func(tripper http.RoundTripper) http.RoundTripper
type Options []Option

type SessionOption func(*Session)
type RequestOption func(*Request)

type RoundTripFunc func(*http.Request) (*http.Response, error)

func (rt RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

// SetHeader to the request.
func SetHeader(key, value string) RequestOption {
	return func(req *Request) {
		req.Header.Set(key, value)
	}
}

// Add a timeout to the request.
func Timeout(timeout time.Duration) SessionOption {
	return func(s *Session) {
		s.Timeout = timeout
	}
}

// Add a basic auth username and password to the request.
func BasicAuth(username, password string) RequestOption {
	return func(req *Request) {
		if username != "" && password != "" {
			req.SetBasicAuth(username, password)
		}
	}
}

// Add a User-Agent header to the request.
func UserAgent(userAgent string) RequestOption {
	return func(req *Request) {
		req.Header.Set("User-Agent", userAgent)
	}
}

// Add cookies to the request.
func Cookies(cookies []*http.Cookie) RequestOption {
	return func(req *Request) {
		for _, c := range cookies {
			req.AddCookie(c)
		}
	}
}

// Add a cookie jar to the request.
func CookieJar(jar http.CookieJar) SessionOption {
	return func(s *Session) {
		s.Jar = jar
	}
}

// Bearer Add a bearer header of the form "Authorization: Bearer ..."
func Bearer(bearer string) RequestOption {
	return func(req *Request) {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearer))
	}
}

// Turn redirects off.
func DisallowRedirects() SessionOption {
	return func(s *Session) {
		s.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
}

// Applies an Accept header to the request.
func Accept(accept string) RequestOption {
	return func(req *Request) {
		req.Header.Set("Accept", accept)
	}
}

// Adds an API key to the request.
func ApiKey(key string) RequestOption {
	return func(req *Request) {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", key))
	}
}

// Pass a context into the HTTP request (allows for request cancellation, for example).
func WithContext(ctx context.Context) RequestOption {
	return func(req *Request) {
		req.Request = req.Request.WithContext(ctx)
	}
}

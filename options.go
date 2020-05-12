package rek

import (
	"context"
	"net/http"
	"time"
)

type options struct {
	headers           map[string]string
	timeout           time.Duration
	username          string
	password          string
	userAgent         string
	cookies           []*http.Cookie
	cookieJar         *http.CookieJar
	bearer            string
	disallowRedirects bool
	accept            string
	apiKey            string
	ctx               context.Context
}

type option func(*options)

// Add headers to the request.
func Headers(headers map[string]string) option {
	return func(opts *options) {
		opts.headers = headers
	}
}

// Add a timeout to the request.
func Timeout(timeout time.Duration) option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}

// Add a basic auth username and password to the request.
func BasicAuth(username, password string) option {
	return func(opts *options) {
		opts.username = username
		opts.password = password
	}
}

// Add a User-Agent header to the request.
func UserAgent(agent string) option {
	return func(opts *options) {
		opts.userAgent = agent
	}
}

// Add cookies to the request.
func Cookies(cookies []*http.Cookie) option {
	return func(opts *options) {
		opts.cookies = cookies
	}
}

// Add a cookie jar to the request.
func CookieJar(jar http.CookieJar) option {
	return func(opts *options) {
		opts.cookieJar = &jar
	}
}

// Add a bearer header of the form "Authorization: Bearer ..."
func Bearer(bearer string) option {
	return func(opts *options) {
		opts.bearer = bearer
	}
}

// Turn redirects off.
func DisallowRedirects() option {
	return func(opts *options) {
		opts.disallowRedirects = true
	}
}

// Applies an Accept header to the request.
func Accept(accept string) option {
	return func(opts *options) {
		opts.accept = accept
	}
}

// Adds an API key to the request.
func ApiKey(key string) option {
	return func(opts *options) {
		opts.apiKey = key
	}
}

// Pass a context into the HTTP request (allows for request cancellation, for example).
func WithContext(ctx context.Context) option {
	return func(opts *options) {
		opts.ctx = ctx
	}
}

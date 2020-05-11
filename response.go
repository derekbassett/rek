package rek

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// A struct containing the relevant response information returned by a rek request.
type Response struct {
	Status string // e.g. "200 OK"

	// StatusCode of the response (200, 404, etc.)
	StatusCode int

	// Content is the response body as raw bytes.
	Content []byte

	// Headers associated with the response.
	Headers http.Header

	// Encoding the response's encoding.
	Encoding []string

	// ContentLength records the length of the associated content. The
	// value -1 indicates that the length is unknown. Unless Request.Method
	// is "HEAD", values >= 0 indicate that the given number of bytes may
	// be read from Body.
	ContentLength int64

	cookies []*http.Cookie
	res     *http.Response
}

// The response body as a string.
func (r *Response) Text() string {
	return string(r.Content)
}

// Marshal a JSON response body.
func (r *Response) Json(v interface{}) error {
	return json.NewDecoder(bytes.NewBuffer(r.Content)).Decode(v)
}

// The Content-Type header for the request (if any).
func (r *Response) ContentType() string {
	if r.Headers == nil {
		return ""
	}
	return r.Headers.Get("Content-Type")
}

// The raw *http.Response returned by the operation.
func (r *Response) Raw() *http.Response {
	return r.res
}

// The cookies associated with the response.
func (r *Response) Cookies() []*http.Cookie {
	return r.cookies
}

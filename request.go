package rek

import (
	"net/http"
)

var DefaultRequest = &Request{}

// GET request
func Get(url string, opts ...option) (*Response, error) {
	return DefaultRequest.Do(http.MethodGet, url, nil, opts...)
}

// POST request
func Post(url string, reader *BodyReader, opts ...option) (*Response, error) {
	return DefaultRequest.Do(http.MethodPost, url, reader, opts...)
}

// PUT request
func Put(url string, reader *BodyReader, opts ...option) (*Response, error) {
	return DefaultRequest.Do(http.MethodPut, url, reader, opts...)
}

// DELETE request
func Delete(url string, opts ...option) (*Response, error) {
	return DefaultRequest.Do(http.MethodDelete, url, nil, opts...)
}

// PATCH request
func Patch(url string, reader *BodyReader, opts ...option) (*Response, error) {
	return DefaultRequest.Do(http.MethodPatch, url, reader, opts...)
}

// HEAD request
func Head(url string, opts ...option) (*Response, error) {
	return DefaultRequest.Do(http.MethodHead, url, nil, opts...)
}

type Request struct {
	Endpoint string
	Method   string
	Body     *BodyReader
	Opts     []option
}

func (r *Request) Do(method, endpoint string, body *BodyReader, opts ...option) (*Response, error) {
	session := NewSession()
	defer session.Close()
	return session.Request(method, endpoint, body, opts...)
}

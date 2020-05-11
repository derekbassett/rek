package rek

import (
	"net/http"
)

var DefaultRequest = &Request{}

// Content-Type MIME of the most common data formats.
const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain; charset=utf-8"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
)

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
	Endpoint          string
	Method            string
	ContentTypeReader ContentTypeReader
	Opts              []option
}

func (r *Request) Do(method, endpoint string, contentTypeReader ContentTypeReader, opts ...option) (*Response, error) {
	session := NewSession()
	defer session.Close()
	return session.Request(method, endpoint, contentTypeReader, opts...)
}

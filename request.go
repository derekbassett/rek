package rek

import (
	"io"
	"net/http"
)

func NewRequest(method, endpoint string, bodyReader *BodyReader, opts ...Option) (*Request, error) {
	var body io.Reader
	var contentType string

	if bodyReader != nil {
		body = bodyReader
		contentType = bodyReader.ContentType
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	if bodyReader != nil && contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return &Request{
		Request:  req,
		Endpoint: endpoint,
		Method:   method,
		Body:     bodyReader,
		Opts:     opts,
	}, nil
}

type Request struct {
	*http.Request
	Endpoint string
	Method   string
	Body     *BodyReader
	Opts     Options
}

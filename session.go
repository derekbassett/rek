package rek

import (
	"net/http"
	"net/url"
)

type Session struct {
	Header http.Header
}

func NewSession() *Session {
	return &Session{
	}
}

// GET request
func (s *Session) Get(url string, opts ...option) (*Response, error) {
	return s.Request(http.MethodGet, url, opts...)
}

// POST request
func (s *Session) Post(url string, opts ...option) (*Response, error) {
	return s.Request(http.MethodPost, url, opts...)
}

// PUT request
func (s *Session) Put(url string, opts ...option) (*Response, error) {
	return s.Request(http.MethodPut, url, opts...)
}

// DELETE request
func (s *Session) Delete(url string, opts ...option) (*Response, error) {
	return s.Request(http.MethodDelete, url, opts...)
}

// PATCH request
func (s *Session) Patch(url string, opts ...option) (*Response, error) {
	return s.Request(http.MethodPatch, url, opts...)
}

// HEAD request
func (s *Session) Head(url string, opts ...option) (*Response, error) {
	options, err := buildOptions(opts...)

	cl := makeClient(options)

	res, err := cl.Head(url)
	if err != nil {
		return nil, err
	}

	return makeResponse(res)
}

func (s *Session) Send(request PreparedRequest) (*Response, error) {
	u, err := url.Parse(request.Endpoint)
	if err != nil {
		return nil, err
	}

	options, err := buildOptions(request.Opts...)
	if err != nil {
		return nil, err
	}

	cl := makeClient(options)

	req, err := makeRequest(request.Method, u.String(), options)

	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}

	resp, err := makeResponse(res)
	if err != nil {
		return nil, err
	}

	if options.callback != nil {
		options.callback(resp)
	}

	return resp, nil
}

func (s *Session) PreparedRequest(request Request) PreparedRequest {
	return PreparedRequest{
		Endpoint: request.Endpoint,
		Method: request.Method,
		Opts: request.Opts,
	}
}

func (s *Session) Request(method, endpoint string, opts ...option) (*Response, error) {
	req := Request{
		Endpoint: endpoint,
		Method: method,
		Opts: opts,
	}
	prep := s.PreparedRequest(req)
	return s.Send(prep)
}

func (s *Session) Close() error {
	return nil
}

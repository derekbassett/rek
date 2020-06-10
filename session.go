package rek

import (
	"net/http"
)

//go:generate gtrace

type defaultTrace struct {
}

func (defaultTrace) OnTrace() func(error) {
	return func(error) {
		return
	}
}

var DefaultSession = &Session{}
var DefaultTrace = &defaultTrace{}

// GET request
func Get(url string, opts ...Option) (*Response, error) {
	return Do(http.MethodGet, url, nil, opts...)
}

// POST request
func Post(url string, reader *BodyReader, opts ...Option) (*Response, error) {
	return Do(http.MethodPost, url, reader, opts...)
}

// PUT request
func Put(url string, reader *BodyReader, opts ...Option) (*Response, error) {
	return Do(http.MethodPut, url, reader, opts...)
}

// DELETE request
func Delete(url string, opts ...Option) (*Response, error) {
	return Do(http.MethodDelete, url, nil, opts...)
}

// PATCH request
func Patch(url string, reader *BodyReader, opts ...Option) (*Response, error) {
	return Do(http.MethodPatch, url, reader, opts...)
}

// HEAD request
func Head(url string, opts ...Option) (*Response, error) {
	return Do(http.MethodHead, url, nil, opts...)
}

func Do(method, endpoint string, body *BodyReader, opts ...Option) (*Response, error) {
	req, err := NewRequest(method, endpoint, body, opts...)
	if err != nil {
		return nil, err
	}
	return DefaultSession.Send(req)
}

type Tracer interface {
	OnTrace() func(error)
}

type SessionTrace struct {
	OnTrace func() func(error)
}

type Session struct {
	http.Client
	options Options
	Trace   Tracer
}

// Get request
func (s *Session) Get(url string, opts ...Option) (*Response, error) {
	return s.Request(http.MethodGet, url, nil, opts...)
}

// Post request
func (s *Session) Post(url string, bodyReader *BodyReader, opts ...Option) (*Response, error) {
	return s.Request(http.MethodPost, url, bodyReader, opts...)
}

// Put request
func (s *Session) Put(url string, bodyReader *BodyReader, opts ...Option) (*Response, error) {
	return s.Request(http.MethodPut, url, bodyReader, opts...)
}

// Delete request
func (s *Session) Delete(url string, opts ...Option) (*Response, error) {
	return s.Request(http.MethodDelete, url, nil, opts...)
}

// Patch request
func (s *Session) Patch(url string, bodyReader *BodyReader, opts ...Option) (*Response, error) {
	return s.Request(http.MethodPatch, url, bodyReader, opts...)
}

// Head request
func (s *Session) Head(url string, opts ...Option) (*Response, error) {
	return s.Request(http.MethodHead, url, nil, opts...)
}

func (s *Session) Request(method, endpoint string, body *BodyReader, opts ...Option) (*Response, error) {
	req, err := NewRequest(method, endpoint, body, opts...)
	if err != nil {
		return nil, err
	}
	return s.Send(req)
}

func (s *Session) Send(request *Request) (*Response, error) {
	done := s.trace().OnTrace()
	resp, err := s.doSend(request)
	done(err)
	return resp, err
}

func (s *Session) trace() Tracer {
	if s.Trace != nil {
		return s.Trace
	}
	return DefaultTrace
}

func (s *Session) doSend(request *Request) (*Response, error) {
	var opts Options
	opts = append(opts, s.options...)
	opts = append(opts, request.Opts...)
	var send http.RoundTripper = RoundTripFunc(func(request *http.Request) (*http.Response, error) {
		return s.Do(request)
	})
	for i := range opts {
		i = len(opts) - i - 1
		option := opts[i]
		if option == nil {
			panic("options cannot be nil")
		}
		send = option(send)
	}
	s.Transport = s.transport()
	res, err := send.RoundTrip(request.Request)
	if err != nil {
		return nil, err
	}

	resp, err := ReadResponse(res)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Session) transport() http.RoundTripper {
	if s.Transport != nil {
		return s.Transport
	}
	return http.DefaultTransport
}

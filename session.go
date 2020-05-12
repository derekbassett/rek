package rek

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Session struct {
	http.Client
}

// Get request
func (s *Session) Get(url string, opts ...option) (*Response, error) {
	return s.Request(http.MethodGet, url, nil, opts...)
}

// Post request
func (s *Session) Post(url string, bodyReader *BodyReader, opts ...option) (*Response, error) {
	return s.Request(http.MethodPost, url, bodyReader, opts...)
}

// Put request
func (s *Session) Put(url string, bodyReader *BodyReader, opts ...option) (*Response, error) {
	return s.Request(http.MethodPut, url, bodyReader, opts...)
}

// Delete request
func (s *Session) Delete(url string, opts ...option) (*Response, error) {
	return s.Request(http.MethodDelete, url, nil, opts...)
}

// Patch request
func (s *Session) Patch(url string, bodyReader *BodyReader, opts ...option) (*Response, error) {
	return s.Request(http.MethodPatch, url, bodyReader, opts...)
}

// Head request
func (s *Session) Head(url string, opts ...option) (*Response, error) {
	return s.Request(http.MethodHead, url, nil, opts...)
}

func (s *Session) Request(method, endpoint string, body *BodyReader, opts ...option) (*Response, error) {
	req := Request{
		Endpoint: endpoint,
		Method:   method,
		Body:     body,
		Opts:     opts,
	}
	return s.Send(req)
}

func (s *Session) Send(request Request) (*Response, error) {
	req, err := makeRequest(request.Method, request.Endpoint, request.Body)
	if err != nil {
		return nil, err
	}
	options, err := buildOptions(request.Opts...)
	if err != nil {
		return nil, err
	}

	req = applyRequestOptions(req, options)
	s.applyClientOptions(options)
	s.Transport = s.transport()
	res, err := s.Do(req)
	if err != nil {
		return nil, err
	}

	resp, err := makeResponse(res)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Session) Close() error {
	return nil
}

func (s *Session) transport() http.RoundTripper {
	if s.Transport != nil {
		return s.Transport
	}
	return http.DefaultTransport
}

func buildOptions(opts ...option) (*options, error) {
	os := &options{
		headers: nil,
		timeout: 10 * time.Second,
	}

	for _, opt := range opts {
		opt(os)
	}

	return os, nil
}

func makeRequest(method, endpoint string, bodyReader *BodyReader) (*http.Request, error) {
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
	return req, nil
}

func applyRequestOptions(req *http.Request, opts *options) *http.Request {
	if opts.ctx != nil {
		req = req.WithContext(opts.ctx)
	}

	if opts.headers != nil {
		for k, v := range opts.headers {
			req.Header.Set(k, v)
		}
	}

	if opts.userAgent != "" {
		req.Header.Set("User-Agent", opts.userAgent)
	}

	if opts.accept != "" {
		req.Header.Set("Accept", opts.accept)
	}

	if opts.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", opts.apiKey))
	}

	if opts.bearer != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", opts.bearer))
	}

	if opts.username != "" && opts.password != "" {
		req.SetBasicAuth(opts.username, opts.password)
	}

	if opts.cookies != nil {
		for _, c := range opts.cookies {
			req.AddCookie(c)
		}
	}

	return req
}

func (s *Session) applyClientOptions(opts *options) {
	if opts.cookieJar != nil {
		s.Jar = *opts.cookieJar
	}

	if opts.timeout != 0 {
		s.Timeout = opts.timeout
	}

	if opts.disallowRedirects {
		s.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
}

func makeResponse(res *http.Response) (*Response, error) {
	resp := &Response{
		Status:        res.Status,
		StatusCode:    res.StatusCode,
		ContentLength: res.ContentLength,
		Encoding:      res.TransferEncoding,
		res:           res,
	}

	if res.Header != nil {
		headers := make(map[string][]string, len(res.Header))

		for k, v := range res.Header {
			headers[k] = v
		}

		resp.Headers = headers
	}

	if res.Body != nil {
		defer res.Body.Close()

		bs, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		resp.Content = bs
	}

	if res.Cookies() != nil {
		resp.cookies = res.Cookies()
	}

	return resp, nil
}

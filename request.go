package rek

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func makeRequest(method, endpoint string, opts *options) (*http.Request, error) {
	var body io.Reader
	var contentType string
	var req *http.Request
	var err error

	if opts.data != nil {
		data, err := getData(opts)
		if err != nil {
			return nil, err
		}

		body = data
	}

	if opts.jsonObj != nil {
		js, err := getJson(opts)
		if err != nil {
			return nil, err
		}

		body = js
	}

	if opts.file != nil {
		b, ct, err := buildMultipartBody(opts)
		if err != nil {
			return nil, err
		}

		contentType = ct
		body = b
	}

	if opts.formData != nil {
		form := url.Values{}

		for k, v := range opts.formData {
			form.Set(k, v)
		}

		body = strings.NewReader(form.Encode())
	}

	req, err = http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	if opts.ctx != nil {
		req = req.WithContext(opts.ctx)
	}

	setHeaders(req, opts)

	if opts.file != nil {
		req.Header.Set("Content-Type", contentType)
	}

	if opts.bearer != "" {
		bearerHeader := fmt.Sprintf("Bearer %s", opts.bearer)
		req.Header.Set("Authorization", bearerHeader)
	}

	setBasicAuth(req, opts)

	setCookies(req, opts)

	if opts.reqModifier != nil {
		opts.reqModifier(req)
	}

	return req, nil
}

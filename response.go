package rek

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

func makeResponse(res *http.Response) (*Response, error) {
	resp := &Response{
		Status:        res.Status,
		StatusCode:    res.StatusCode,
		ContentLength: res.ContentLength,
		res:           res,
	}

	if res.Header != nil {
		headers := make(map[string][]string)

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

	if res.TransferEncoding != nil {
		resp.Encoding = res.TransferEncoding
	}

	if res.Cookies() != nil {
		resp.cookies = res.Cookies()
	}

	return resp, nil
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

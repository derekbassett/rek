package rek

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
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

	// The cookies associated with the response.
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

// Marshal a XML response body.
func (r *Response) Xml(v interface{}) error {
	return xml.NewDecoder(bytes.NewBuffer(r.Content)).Decode(v)
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

func (r *Response) Cookies() []*http.Cookie {
	return r.cookies
}

func ReadResponse(res *http.Response) (*Response, error) {
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
		// According to the comments
		// | The http Client and Transport guarantee that Body is always
		// | non-nil, even on responses without a body or responses with
		// | a zero-length body.
		var (
			save io.ReadCloser
			err  error
		)
		save, res.Body, err = drainBody(res.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to drain body from the response: %w", err)
		}
		resp.Content, err = ioutil.ReadAll(save)
		if err != nil {
			return nil, fmt.Errorf("unable to read data from the response: %w", err)
		}
	}

	if res.Cookies() != nil {
		resp.cookies = res.Cookies()
	}

	return resp, nil
}

// drainBody takes a readcloser and drains it of data supplying to identical readclosers.  Both readcloser have
// the exact same data.
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

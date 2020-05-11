package rek

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDo(t *testing.T) {
	tt := map[string]struct {
		method       string
		path         string
		requestBody  string
		responseBody string
		statusCode   int
	}{
		"POST happy path": {
			method:       http.MethodPost,
			path:         "/v1/test",
			requestBody:  "Hello World",
			responseBody: "Complete",
			statusCode:   http.StatusOK,
		},
		"GET happy path": {
			method:       http.MethodGet,
			path:         "/v1/test",
			requestBody:  "",
			responseBody: "Complete",
			statusCode:   http.StatusOK,
		},
		"PUT happy path": {
			method:       http.MethodPut,
			path:         "/v1/test",
			requestBody:  "Hello",
			responseBody: "World",
			statusCode:   http.StatusOK,
		},
		"POST status code": {
			method:       http.MethodPost,
			path:         "/v1",
			requestBody:  "Hello",
			responseBody: "World",
			statusCode:   http.StatusInternalServerError,
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, req *http.Request) {
				if tc.method != req.Method {
					t.Errorf("method is wrong, got %q want %q", req.Method, tc.method)
				}
				if !strings.Contains(req.URL.Path, tc.path) {
					t.Errorf("path does not contain expected path, got %q want %q", req.URL.Path, tc.path)
				}
				b, err := ioutil.ReadAll(req.Body)
				if err != nil {
					t.Fatal(err)
				}
				if got, want := string(bytes.TrimSpace(b)), tc.requestBody; got != want {
					t.Errorf("request body is wrong, got %q want %q", got, want)
				}
				w.WriteHeader(tc.statusCode)
				if tc.responseBody != "" {
					_, _ = w.Write([]byte(tc.responseBody))
				}
			}
			srv := httptest.NewServer(http.HandlerFunc(handler))
			defer srv.Close()

			endpoint := fmt.Sprintf("%s%s", srv.URL, tc.path)
			resp, err := DefaultRequest.Do(tc.method, endpoint, String(tc.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			if got, want := resp.StatusCode, tc.statusCode; got != want {
				t.Errorf("status code is wrong, got %q want %q", got, want)
			}
			if got, want := int(resp.ContentLength), len(tc.responseBody); got != want {
				t.Errorf("content length is wrong, got %d want %d", got, want)
			}
			if got, want := resp.Content, []byte(tc.responseBody); !bytes.Equal(got, want) {
				t.Errorf("content is wrong, got %s want %s", got, want)
			}
		})
	}
}

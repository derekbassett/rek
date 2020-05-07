package rek

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCall(t *testing.T) {
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
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, req *http.Request) {
				assert.Equal(t, tc.method, req.Method)
				assert.Contains(t, req.URL.Path, tc.path)
				b, err := ioutil.ReadAll(req.Body)
				if tc.requestBody != "" && assert.NoError(t, err) {
					got := string(bytes.TrimSpace(b))
					assert.Equal(t, tc.requestBody, got)
				}
				w.WriteHeader(tc.statusCode)
				if tc.responseBody != "" {
					w.Write([]byte(tc.responseBody))
				}
			}

			srv := httptest.NewServer(http.HandlerFunc(handler))
			defer srv.Close()
			endpoint := fmt.Sprintf("%s%s", srv.URL, tc.path)
			resp, err := call(tc.method, endpoint, String(tc.requestBody))

			if assert.NoError(t, err) {
				assert.Equal(t, tc.statusCode, resp.StatusCode())
				assert.Equal(t, len(tc.responseBody), int(resp.ContentLength()))
				assert.Equal(t, []byte(tc.responseBody), resp.Content())
			}
		})
	}
}

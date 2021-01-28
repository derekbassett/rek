package rek

import (
	"bytes"
	"fmt"
	"go/token"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestResponse_ContentType(t *testing.T) {
	tests := []struct {
		Raw  *http.Response
		Resp Response
		Body string
	}{
		{
			Raw:  &http.Response{StatusCode: 200},
			Resp: Response{StatusCode: 200},
			Body: "",
		},
	}
	for i, tt := range tests {
		resp, err := ReadResponse(tt.Raw)
		if err != nil {
			t.Errorf("#%d: %v", i, err)
			continue
		}
		rContent := resp.Content
		diff(t, fmt.Sprintf("#%d Response", i), resp, &tt.Resp)
		var bout bytes.Buffer
		if rContent != nil {
			_, err = io.Copy(&bout, bytes.NewReader(rContent))
			if err != nil {
				t.Errorf("#%d: %v", i, err)
				continue
			}
		}
		body := bout.String()
		if body != tt.Body {
			t.Errorf("#%d: Body = %q want %q", i, body, tt.Body)
		}
	}
}

func diff(t *testing.T, prefix string, have, want interface{}) {
	hv := reflect.ValueOf(have).Elem()
	wv := reflect.ValueOf(want).Elem()
	if hv.Type() != wv.Type() {
		t.Errorf("%s: type mismatch %v want %v", prefix, hv.Type(), wv.Type())
	}
	for i := 0; i < hv.NumField(); i++ {
		name := hv.Type().Field(i).Name
		if !token.IsExported(name) {
			continue
		}
		hf := hv.Field(i).Interface()
		wf := wv.Field(i).Interface()
		if !reflect.DeepEqual(hf, wf) {
			t.Errorf("%s: %s = %v want %v", prefix, name, hf, wf)
		}
	}
}

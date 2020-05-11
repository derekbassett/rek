package rek

import (
	"errors"
	"net/http"
	"testing"
)

type recordingTransport struct {
	req *http.Request
}

func (t *recordingTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	t.req = req
	return nil, errors.New("dummy impl")
}

func TestGetRequestFormat(t *testing.T) {
	tr := &recordingTransport{}
	session := &Session{Transport: tr}
	url := "http://dummy.faketld/"
	_, _ = session.Get(url)
	if tr.req.Method != "GET" {
		t.Errorf("expected method %q; got %q", "GET", tr.req.Method)
	}
	if tr.req.URL.String() != url {
		t.Errorf("expected URL %q; got %q", url, tr.req.URL.String())
	}
}

func TestPutRequestFormat(t *testing.T) {
	tr := &recordingTransport{}
	session := &Session{Transport: tr}
	url := "http://dummy.faketld/"
	_, _ = session.Put(url, String("Reader"))
	if tr.req.Method != "PUT" {
		t.Errorf("expected method %q; got %q", "PUT", tr.req.Method)
	}
	if tr.req.URL.String() != url {
		t.Errorf("expected method %q; got %q", url, tr.req.URL.String())
	}
}

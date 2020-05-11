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
	session.Get(url) //nolint
	if tr.req.Method != http.MethodGet {
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
	session.Put(url, String("Reader")) //nolint
	if tr.req.Method != http.MethodPut {
		t.Errorf("expected method %q; got %q", "PUT", tr.req.Method)
	}
	if tr.req.URL.String() != url {
		t.Errorf("expected method %q; got %q", url, tr.req.URL.String())
	}
}

func TestPatchRequest(t *testing.T) {
	tr := &recordingTransport{}
	session := &Session{Transport: tr}
	url := "http://dummy.faketld/"
	session.Patch(url, String("Reader")) //nolint errcheck
	if tr.req.Method != http.MethodPatch {
		t.Errorf("expected method %q; got %q", http.MethodPatch, tr.req.Method)
	}
	if tr.req.URL.String() != url {
		t.Errorf("expected method %q; got %q", url, tr.req.URL.String())
	}
}

func TestPostRequest(t *testing.T) {
	tr := &recordingTransport{}
	session := &Session{Transport: tr}
	url := "http://dummy.faketld/"
	session.Post(url, String("Reader")) //nolint errcheck
	if tr.req.Method != http.MethodPost {
		t.Errorf("expected method %q; got %q", http.MethodPost, tr.req.Method)
	}
	if tr.req.URL.String() != url {
		t.Errorf("expected method %q; got %q", url, tr.req.URL.String())
	}
}

func TestDeleteRequest(t *testing.T) {
	tr := &recordingTransport{}
	session := &Session{Transport: tr}
	url := "http://dummy.faketld/"
	session.Delete(url) //nolint errcheck
	if tr.req.Method != "DELETE" {
		t.Errorf("expected method %q; got %q", "DELETE", tr.req.Method)
	}
	if tr.req.URL.String() != url {
		t.Errorf("expected method %q; got %q", url, tr.req.URL.String())
	}
}

func TestHeadRequest(t *testing.T) {
	tr := &recordingTransport{}
	session := &Session{Transport: tr}
	url := "http://dummy.faketld/"
	session.Head(url) //nolint errcheck
	if tr.req.Method != http.MethodHead {
		t.Errorf("expected method %q; got %q", http.MethodHead, tr.req.Method)
	}
	if tr.req.URL.String() != url {
		t.Errorf("expected method %q; got %q", url, tr.req.URL.String())
	}
}

func TestRequest(t *testing.T) {
	tr := &recordingTransport{}
	session := &Session{Transport: tr}
	url := "http://dummy.faketld/"
	session.Request(http.MethodOptions, url, nil) //nolint errcheck
	if tr.req.Method != http.MethodOptions {
		t.Errorf("expected method %q; got %q", http.MethodOptions, tr.req.Method)
	}
	if tr.req.URL.String() != url {
		t.Errorf("expected method %q; got %q", url, tr.req.URL.String())
	}
}

func TestWithTransport(t *testing.T) {
	t.Fatal("not implemented")
}

func TestWithClient(t *testing.T) {
	t.Fatal("not implemented")
}

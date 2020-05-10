package rek

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"net/url"
	"strings"
)

type ContentTypeReader interface {
	io.Reader
	ContentType() string
}

type BodyReader struct {
	io.Reader
	contentType string
}

//
func Must(reader *BodyReader, err error) *BodyReader {
	if err != nil {
		panic(err)
	}
	return reader
}

func (b *BodyReader) ContentType() string {
	return b.contentType
}

// Add any interface{} that can be serialized to a []byte and apply a "Content-Type: application/octet-stream" header.
func EncodeBinary(data interface{}) (*BodyReader, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(data)
	if err != nil {
		return nil, err
	}
	return &BodyReader{
		bytes.NewReader(buf.Bytes()),
		"application/octet-stream",
	}, nil
}

// Add a String that can be sent
func String(text string) *BodyReader {
	return &BodyReader{
		strings.NewReader(text),
		"text/plain",
	}
}

// Add any interface{} that can be marshaled as JSON to the request body and apply a "Content-Type:
// application/json;charset=utf-8" header.
func EncodeJson(v interface{}) (*BodyReader, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(v)
	if err != nil {
		return nil, err
	}
	return &BodyReader{
		bytes.NewReader(buf.Bytes()),
		"application/json; charset=utf-8",
	}, nil
}

// Add key/value form data to the request body and apply a "Content-Type: application/x-www-form-urlencoded" header.
func EncodeFormData(formData map[string]string) (*BodyReader, error) {
	form := url.Values{}
	for k, v := range formData {
		form.Set(k, v)
	}
	return &BodyReader{
		strings.NewReader(form.Encode()),
		"application/x-www-form-urlencoded",
	}, nil
}

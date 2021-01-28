package rek

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

const defaultFileFieldName = "file"

// Create a multipart file upload request.
func EncodeFile(fieldName, filepath string, params map[string]string) (*BodyReader, error) {
	file := &file{
		FieldName: fieldName,
		Filepath:  filepath,
		Params:    params,
	}
	b, ct, err := buildMultipartBody(file)
	if err != nil {
		return nil, err
	}
	return &BodyReader{
		b,
		ct,
	}, nil
}

type file struct {
	FieldName string
	Filepath  string
	Params    map[string]string
}

func (f *file) build() *file {
	if f.FieldName == "" {
		f.FieldName = defaultFileFieldName
	}

	return f
}

func buildMultipartBody(file *file) (io.Reader, string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	f := file.build()

	data, err := os.Open(f.Filepath)
	if err != nil {
		return nil, "", err
	}

	part, err := writer.CreateFormFile(f.FieldName, filepath.Base(f.Filepath))
	if err != nil {
		return nil, "", err
	}

	if _, err := io.Copy(part, data); err != nil {
		return nil, "", err
	}

	if f.Params != nil {
		for k, v := range f.Params {
			if err := writer.WriteField(k, v); err != nil {
				return nil, "", err
			}
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return &body, writer.FormDataContentType(), nil
}

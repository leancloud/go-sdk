package leancloud

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"testing"
)

func beforeTestFileRef() (*File, error) {
	filename, err := generateTempFile("go-sdk-file-upload-*.txt")
	if err != nil {
		return nil, err
	}

	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	_, name := filepath.Split(filename)
	file := &File{
		Name: name,
		MIME: mime.TypeByExtension(filepath.Ext(name)),
	}

	if err := c.Files.Upload(file, fd); err != nil {
		return nil, err
	}

	if err := checkUpload(file.URL); err != nil {
		return nil, err
	}

	return file, nil
}

func TestFileRefGet(t *testing.T) {
	file, err := beforeTestFileRef()
	if err != nil {
		t.Fatal(err)
	}

	fetch, err := client.File(file.ID).Get()
	if err != nil {
		t.Fatal(err)
	}

	if fetch.ID != file.ID {
		t.Fatal(fmt.Errorf("failed to get File object"))
	}
}

func TestFileRefDestroy(t *testing.T) {
	file, err := beforeTestFileRef()
	if err != nil {
		t.Fatal(err)
	}

	if err := client.File(file.ID).Destroy(UseMasterKey(true)); err != nil {
		t.Fatal(err)
	}

	resp, err := client.request(ServiceAPI, MethodGet, fmt.Sprint("/1.1/files/", file.ID), nil)
	if err != nil {
		if resp.StatusCode != 404 {
			t.Fatal(err)
		}
	}
}

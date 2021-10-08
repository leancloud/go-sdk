package leancloud

import (
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"testing"

	"github.com/levigross/grequests"
)

func generateTempFile(pattern string) (string, error) {
	content := []byte("temporary file's content")
	tmpfile, err := ioutil.TempFile("", "go-sdk-file-upload-*.txt")
	if err != nil {
		return "", err
	}

	if _, err := tmpfile.Write(content); err != nil {
		return "", err
	}

	name := tmpfile.Name()

	if err := tmpfile.Close(); err != nil {
		return "", err
	}

	return name, nil
}

func checkUpload(url string) error {
	resp, err := grequests.Get(url, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unable to get file with url: %v", url)
	}

	return nil
}
func TestFilesUpload(t *testing.T) {
	filename, err := generateTempFile("go-sdk-file-upload-*.txt")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(filename)

	fd, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	_, name := filepath.Split(filename)
	file := &File{
		Name: name,
		MIME: mime.TypeByExtension(filepath.Ext(name)),
	}

	if err := testC.Files.Upload(file, fd); err != nil {
		t.Fatal(err)
	}

	if err := checkUpload(file.URL); err != nil {
		t.Fatal(err)
	}
}

func TestFilesUploadWithOwner(t *testing.T) {
	filename, err := generateTempFile("go-sdk-file-upload-*.txt")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(filename)

	fd, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	_, name := filepath.Split(filename)
	file := &File{
		Name: name,
		MIME: mime.TypeByExtension(filepath.Ext(name)),
	}

	user, err := testC.Users.LogIn(testUsername, testPassword)
	if err != nil {
		t.Fatal(err)
	}

	if err := testC.Files.Upload(file, fd, UseUser(user)); err != nil {
		t.Fatal(err)
	}

	if err := checkUpload(file.URL); err != nil {
		t.Fatal(err)
	}
}
func TestFilesUploadFromURL(t *testing.T) {
	file := &File{
		Name: "go-sdk-file-upload.txt",
		MIME: "text/plain",
		URL:  "https://example.com/assets/go-sdk-file-upload.txt",
	}

	if err := testC.Files.UploadFromURL(file); err != nil {
		t.Fatal(err)
	}

	if file.ID == "" {
		t.Fatal("unable to create _File object")
	}
}

func TestFilesUploadFromFile(t *testing.T) {
	filename, err := generateTempFile("go-sdk-file-upload-*.txt")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(filename)

	file := &File{
		Meatadata: map[string]interface{}{
			"comment": "This is a comment of Metadata",
		},
	}

	if err := testC.Files.UploadFromLocalFile(file, filename); err != nil {
		t.Fatal(err)
	}

	if err := checkUpload(file.URL); err != nil {
		t.Fatal(err)
	}
}

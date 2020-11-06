package leancloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/levigross/grequests"
)

type File struct {
	Object
	Key       string                 `json:"key"`
	Name      string                 `json:"name"`
	Provider  string                 `json:"provider"`
	Bucket    string                 `json:"bucket"`
	MIME      string                 `json:"mime_type"`
	URL       string                 `json:"url"`
	Size      int64                  `json:"size"`
	Meatadata map[string]interface{} `json:"metadata"`
}

// GetMap export raw hashmap of File object
func (file *File) GetMap() map[string]interface{} {
	return nil
}

// Get export the value by given key in File object
func (file *File) Get(key string) interface{} {
	return file.fields[key]
}

func (file *File) fetchToken(client *Client, authOptions ...AuthOption) (string, string, error) {
	reqJSON := encodeFile(file, false)

	path := "/1.1/fileTokens"
	options := client.getRequestOptions()
	options.JSON = reqJSON

	resp, err := client.request(ServiceAPI, MethodPost, path, options, authOptions...)
	if err != nil {
		return "", "", err
	}

	respJSON := make(map[string]interface{})
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return "", "", err
	}

	objectID, ok := respJSON["objectId"].(string)
	if !ok {
		return "", "", fmt.Errorf("unexpected error when parse objectId from response: want type string but %v", reflect.TypeOf(respJSON["objectId"]))
	}
	file.ID = objectID

	createdAt, ok := respJSON["createdAt"].(string)
	if !ok {
		return "", "", fmt.Errorf("unexpected error when parse createdAt from response: want type string but %v", reflect.TypeOf(respJSON["createdAt"]))
	}
	decodedCreatedAt, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return "", "", fmt.Errorf("unexpected error when parse createdAt from response: %v", err)
	}
	file.CreatedAt = decodedCreatedAt

	key, ok := respJSON["key"].(string)
	if !ok {
		return "", "", fmt.Errorf("unexpected error when parse key from response: want type string but %v", reflect.TypeOf(respJSON["key"]))
	}
	file.Key = key

	url, ok := respJSON["url"].(string)
	if !ok {
		return "", "", fmt.Errorf("unexpected error when parse url from response: want type string but %v", reflect.TypeOf(respJSON["url"]))
	}
	file.URL = url

	token, ok := respJSON["token"].(string)
	if !ok {
		return "", "", fmt.Errorf("unexpected error when parse token from response: want type string but %v", reflect.TypeOf(respJSON["token"]))
	}

	bucket, ok := respJSON["bucket"].(string)
	if !ok {
		return "", "", fmt.Errorf("unexpected error when parse bucket from response: want type string but %v", reflect.TypeOf(respJSON["bucket"]))
	}
	file.Bucket = bucket

	uploadURL, ok := respJSON["upload_url"].(string)
	if !ok {
		return "", "", fmt.Errorf("unexpected error when parse upload_url from response: want type but %v", reflect.TypeOf(respJSON["upload_url"]))
	}

	return token, uploadURL, nil
}

func (file *File) fileCallback(result bool, token string, client *Client, authOptions ...AuthOption) error {
	path := "/1.1/fileCallback"
	options := client.getRequestOptions()
	options.JSON = map[string]interface{}{
		"result": result,
		"token":  token,
	}

	_, err := client.request(ServiceAPI, MethodPost, path, options, authOptions...)
	if err != nil {
		return err
	}

	return nil
}

func (file *File) uploadQiniu(token, uploadURL string, reader io.ReadSeeker) error {
	out, in := io.Pipe()
	part := multipart.NewWriter(in)
	done := make(chan error)

	go func() {
		if err := part.WriteField("key", file.Key); err != nil {
			in.Close()
			done <- err
			return
		}

		if err := part.WriteField("token", token); err != nil {
			in.Close()
			done <- err
			return
		}

		writer, err := part.CreateFormFile("file", file.Name)
		if err != nil {
			in.Close()
			done <- err
			return
		}

		_, err = io.Copy(writer, reader)
		if err != nil {
			in.Close()
			done <- err
			return
		}

		in.Close()
		done <- nil
	}()

	options := grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type": part.FormDataContentType(),
		},
		RequestBody: out,
	}

	resp, err := grequests.Post(uploadURL, &options)
	if err != nil {
		return fmt.Errorf("unexpected error when upload file to Qiniu: %v", err)
	}

	err = <-done
	if err != nil {
		return fmt.Errorf("unexpected error when upload file to Qiniu: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected error when upload file to Qiniu: %s", string(resp.Bytes()))
	}

	return nil
}

func (file *File) uploadS3(token, uploadURL string, reader io.ReadSeeker) error {
	options := grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type":   file.MIME,
			"Cache-Control":  "public, max-age=31536000",
			"Content-Length": fmt.Sprint("%d", file.Size),
		},
		RequestBody: reader,
	}

	resp, err := grequests.Put(uploadURL, &options)
	if err != nil {
		return fmt.Errorf("unexpected error when upload file to AWS S3: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected error when upload file to AWS S3: %v", string(resp.Bytes()))
	}

	return nil
}

func (file *File) uploadCOS(token, uploadURL string, reader io.ReadSeeker) error {
	out, in := io.Pipe()
	part := multipart.NewWriter(in)
	done := make(chan error)

	go func() {
		if err := part.WriteField("op", "upload"); err != nil {
			in.Close()
			done <- err
			return
		}
		writer, err := part.CreateFormFile("fileContent", file.Name)
		if err != nil {
			in.Close()
			done <- err
			return
		}
		_, err = io.Copy(writer, reader)
		if err != nil {
			in.Close()
			done <- err
			return
		}
		if err := part.Close(); err != nil {
			in.Close()
			done <- err
			return
		}
		in.Close()
		done <- nil
	}()

	reqBody, err := ioutil.ReadAll(out)
	if err != nil {
		return err
	}

	options := grequests.RequestOptions{
		Headers: map[string]string{
			"Content-Type":   part.FormDataContentType(),
			"Content-Length": fmt.Sprint(int64(len(reqBody))),
		},
		RequestBody: bytes.NewBuffer(reqBody),
	}

	resp, err := grequests.Post(fmt.Sprint(uploadURL, "?sign", url.QueryEscape(token)), &options)
	if err != nil {
		return fmt.Errorf("unexpected error when upload file to COS: %v", string(resp.Bytes()))
	}

	err = <-done
	if err != nil {
		return fmt.Errorf("unexpected error when upload file to COS: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected error when upload file to COS: %v", string(resp.Bytes()))
	}

	return nil
}

func getSeekerSize(seeker io.Seeker) (int64, error) {
	size, err := seeker.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	_, err = seeker.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return size, nil
}

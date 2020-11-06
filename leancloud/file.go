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

func (file *File) fetchOwner(authOptions ...AuthOption) (*User, error) {
	options := client.getRequestOptions()
	for _, authOption := range authOptions {
		authOption.apply(client, options)
	}

	if options.Headers["X-LC-Session"] == "" {
		return nil, nil
	}

	user, err := client.Users.Become(options.Headers["X-LC-Session"])
	if err != nil {
		return nil, err
	}

	return user, nil
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
		return "", "", fmt.Errorf("unexpected error when parse upload_url from response: want type string but %v", reflect.TypeOf(respJSON["upload_url"]))
	}

	provider, ok := respJSON["provider"].(string)
	if !ok {
		return "", "", fmt.Errorf("unexpected error when parse provider from response: want type string but %v", reflect.TypeOf(respJSON["provider"]))
	}
	file.Provider = provider

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
		if err := part.Close(); err != nil {
			in.Close()
			done <- err
			return
		}
		in.Close()
		done <- nil
	}()

	req, err := http.NewRequest("POST", "https://up.qbox.me/", out)
	req.Header.Set("Content-Type", part.FormDataContentType())
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected error when upload file to Qiniu: %v", err)
	}

	err = <-done
	if err != nil {
		return fmt.Errorf("unexpected error when upload file to Qiniu: %v", err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unexpected error when upload file to Qiniu: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected error when upload file to Qiniu: %v", string(content))
	}

	return nil
}

func (file *File) uploadS3(token, uploadURL string, reader io.ReadSeeker) error {
	req, err := http.NewRequest("PUT", uploadURL, reader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", file.MIME)
	req.Header.Set("Cache-Control", "public, max-age=31536000")
	req.ContentLength = file.Size

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected error when upload file to AWS S3: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("unexpected error when upload file to AWS S3: %v", err)
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected error when upload file to AWS S3: %v", string(body))
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

	body, err := ioutil.ReadAll(out)
	if err != nil {
		return fmt.Errorf("unexpected error when upload file to COS: %v", err)
	}

	req, err := http.NewRequest("POST", uploadURL+"?sign="+url.QueryEscape(token), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("unexpected error when upload file to COS: %v", err)
	}

	req.Header.Set("Content-Type", part.FormDataContentType())
	req.ContentLength = int64(len(body))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected error when upload file to COS: %v", err)
	}

	err = <-done
	if err != nil {
		return fmt.Errorf("unexpected error when upload file to COS: %v", err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unexpected error when upload file to COS: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected error when upload file to COS: %v", string(content))
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

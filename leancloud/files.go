package leancloud

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

type Files struct {
	c *Client
}

// Upload transfer the file to cloud storage and create an File object in _File class
func (ref *Files) Upload(file *File, reader io.ReadSeeker, authOptions ...AuthOption) error {
	size, err := getSeekerSize(reader)
	if err != nil {
		return fmt.Errorf("unexpected error when get length of file: %v", err)
	}

	owner, err := file.fetchOwner(authOptions...)
	if err != nil {
		return fmt.Errorf("unexpected error when fetch owner: %v", err)
	}

	if file.Size == 0 {
		file.Size = size
	}

	if reflect.ValueOf(file.Meatadata).IsNil() {
		file.Meatadata = make(map[string]interface{})
	}
	file.Meatadata["size"] = file.Size
	if owner != nil {
		file.Meatadata["owner"] = owner.ID
	} else {
		file.Meatadata["owner"] = "unknown"
	}

	token, uploadURL, err := file.fetchToken(ref.c, authOptions...)
	if err != nil {
		return err
	}

	switch file.Provider {
	case "qiniu":
		if err := file.uploadQiniu(token, "https://up.qbox.me/", reader); err != nil {
			if err := file.fileCallback(false, token, ref.c, authOptions...); err != nil {
				return err
			}
			return err
		}
	case "s3":
		if err := file.uploadS3(token, uploadURL, reader); err != nil {
			if err := file.fileCallback(false, token, ref.c, authOptions...); err != nil {
				return err
			}
			return err
		}
	case "qcloud":
		if err := file.uploadCOS(token, uploadURL, reader); err != nil {
			if err := file.fileCallback(false, token, ref.c, authOptions...); err != nil {
				return err
			}
			return err
		}
	}

	if err := file.fileCallback(true, token, ref.c, authOptions...); err != nil {
		return err
	}

	return nil
}

// UploadFromURL create an object of file in _File class with given file's url
func (ref *Files) UploadFromURL(file *File, authOptions ...AuthOption) error {
	if reflect.ValueOf(file.Meatadata).IsNil() {
		file.Meatadata = make(map[string]interface{})
	}
	owner, err := file.fetchOwner(authOptions...)
	if err != nil {
		return fmt.Errorf("unexpected error when fetch owner: %v", err)
	}
	file.Meatadata["__source"] = "external"
	if owner != nil {
		file.Meatadata["owner"] = owner.ID
	} else {
		file.Meatadata["owner"] = "unknown"
	}

	path := "/1.1/files"
	options := ref.c.getRequestOptions()
	options.JSON = encodeFile(file, false)

	resp, err := ref.c.request(methodPost, path, options, authOptions...)
	if err != nil {
		return err
	}

	respJSON := make(map[string]interface{})
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return err
	}

	objectID, ok := respJSON["objectId"].(string)
	if !ok {
		return fmt.Errorf("unexpected error when fetch objectId: want type string but %v", reflect.TypeOf(respJSON["objectId"]))
	}
	file.ID = objectID

	createdAt, ok := respJSON["createdAt"].(string)
	if !ok {
		return fmt.Errorf("unexpected error when fetch createdAt: want type string but %v", reflect.TypeOf(respJSON["createdAt"]))
	}
	decodedCreatedAt, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return fmt.Errorf("unexpected error when parse createdAt: %v", err)
	}
	file.CreatedAt = decodedCreatedAt

	return nil
}

// UploadFromFile transfer the file given by path to cloud storage and create an object in _File class
//
// After uploading it will return an File object
func (ref *Files) UploadFromFile(file *File, path string, authOptions ...AuthOption) error {
	if file.Name == "" {
		_, name := filepath.Split(path)
		file.Name = name
	}

	if file.MIME == "" {
		mime := mime.TypeByExtension(filepath.Ext(path))
		file.MIME = mime
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unexpected error when open %s: %v", path, err)
	}

	if err := ref.c.Files.Upload(file, f, authOptions...); err != nil {
		return err
	}

	return nil
}

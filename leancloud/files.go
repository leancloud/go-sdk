package leancloud

import (
	"fmt"
	"io"
)

type Files struct {
	c *Client
}

func (ref *Files) Upload(name, mime string, reader io.ReadSeeker, authOptions ...AuthOption) (*File, error) {
	size, err := getSeekerSize(reader)
	if err != nil {
		return nil, fmt.Errorf("unexpected error when get length of file: %v", err)
	}

	file := &File{
		Name: name,
		MIME: mime,
		Size: size,
		Meatadata: map[string]interface{}{
			"size":  size,
			"owner": "unknown",
		},
	}

	token, uploadURL, err := file.fetchToken(ref.c, authOptions...)
	if err != nil {
		return nil, err
	}

	switch file.Provider {
	case "qiniu":
		if err := file.uploadQiniu(token, "https://up.qbox.me/", reader); err != nil {
			if err := file.fileCallback(false, token, ref.c, authOptions...); err != nil {
				return nil, err
			}
			return nil, err
		}
	case "s3":
		if err := file.uploadS3(token, uploadURL, reader); err != nil {
			if err := file.fileCallback(false, token, ref.c, authOptions...); err != nil {
				return nil, err
			}
			return nil, err
		}
	case "qcloud":
		if err := file.uploadCOS(token, uploadURL, reader); err != nil {
			if err := file.fileCallback(false, token, ref.c, authOptions...); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	if err := file.fileCallback(true, token, ref.c, authOptions...); err != nil {
		return nil, err
	}

	return file, nil
}

func (ref *Files) UploadFromURL(name, mime, url string, authOptions ...AuthOption) (*File, error) {
	return nil, nil
}

func (ref *Files) UploadFromFile(path string, authOptions ...AuthOption) (*File, error) {

}

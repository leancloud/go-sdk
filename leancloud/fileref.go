package leancloud

import (
	"fmt"
	"reflect"
)

// FileRef refer to a File object in _File class
type FileRef struct {
	c     *Client
	class string
	ID    string
}

// File construct an new reference to a _File object by given objectId
func (client *Client) File(id string) *FileRef {
	return &FileRef{
		c:     client,
		class: "files",
		ID:    id,
	}
}

// Get fetch the referred _File object
func (ref *FileRef) Get(authOptions ...AuthOption) (*File, error) {
	decodedFile, err := objectGet(ref, authOptions...)
	if err != nil {
		return nil, err
	}

	file, ok := decodedFile.(*File)
	if !ok {
		return nil, fmt.Errorf("unexpected error when parse Object from response: want type *File but %v", reflect.TypeOf(decodedFile))
	}

	return file, nil
}

// Destroy delete the referred _File object
func (ref *FileRef) Destroy(authOptions ...AuthOption) error {
	if err := objectDestroy(ref, authOptions...); err != nil {
		return err
	}

	return nil
}

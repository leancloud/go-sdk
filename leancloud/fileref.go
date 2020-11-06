package leancloud

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
	return nil, nil
}

// Destroy delete the referred _File object
func (ref *FileRef) Destroy(authOptions ...AuthOption) error {
	return nil
}

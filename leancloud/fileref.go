package leancloud

type FileRef struct {
	c     *Client
	class string
	ID    string
}

func (client *Client) File(id string) *FileRef {
	return &FileRef{
		c:     client,
		class: "files",
		ID:    id,
	}
}

func (ref *FileRef) Get(authOptions ...AuthOption) (*File, error) {
	return nil, nil
}

func (ref *FileRef) Destroy(authOptions ...AuthOption) error {
	return nil
}

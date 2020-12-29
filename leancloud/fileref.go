package leancloud

// FileRef refer to a File object in _File class
type FileRef struct {
	c     *Client
	class string
	ID    string
}

// Get fetch the referred _File object
func (ref *FileRef) Get(file *File, authOptions ...AuthOption) error {
	err := objectGet(ref, file, authOptions...)
	if err != nil {
		return err
	}

	return nil
}

// Destroy delete the referred _File object
func (ref *FileRef) Destroy(authOptions ...AuthOption) error {
	if err := objectDestroy(ref, authOptions...); err != nil {
		return err
	}

	return nil
}

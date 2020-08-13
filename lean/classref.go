package lean

type ClassRef struct {
	c    *Client
	Name string
}

func (ref *ClassRef) Object(id string) (*ObjectRef, error) {
	// TODO
	return nil, nil
}

func (ref *ClassRef) Create(data interface{}) (*ObjectRef, error) {
	// TODO
	return nil, nil
}

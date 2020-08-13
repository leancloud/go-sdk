package lean

type ObjectRef struct {
	c     *Client
	class string
	ID    string
}

func (r *ObjectRef) Get() (*Object, error) {
	// TODO
	return nil, nil
}

func (r *ObjectRef) Set(data interface{}) error {
	// TODO
	return nil
}

func (r *ObjectRef) Update() error {
	// TODO
	return nil
}

func (r *ObjectRef) Delete() error {
	// TODO
	return nil
}

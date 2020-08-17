package lean

type ClassRef struct {
	c    *Client
	Name string
}

func (client *Client) Class(name string) *ClassRef {
	return &ClassRef{
		c:    client,
		Name: name,
	}
}

func (r *ClassRef) Object(id string) *ObjectRef {
	return &ObjectRef{
		c:     r.c,
		class: r.Name,
		ID:    id,
	}
}

func (r *ClassRef) Create(data interface{}) (*ObjectRef, error) {

	return nil, nil
}

func (r *ClassRef) NewQuery() *Query {
	return &Query{
		c:        r.c,
		classRef: r,
	}
}

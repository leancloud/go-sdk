package leancloud

type Class struct {
	c    *Client
	Name string
}

func (client *Client) Class(name string) *Class {
	return &Class{
		c:    client,
		Name: name,
	}
}

func (r *Class) Object(id string) *ObjectRef {
	return &ObjectRef{
		c:     r.c,
		class: r.Name,
		ID:    id,
	}
}

func (r *Class) Create(data interface{}, auth ...AuthOption) (*ObjectRef, error) {
	// TODO
	return nil, nil
}

func (r *Class) NewQuery() *Query {
	return &Query{
		c:        r.c,
		classRef: r,
	}
}

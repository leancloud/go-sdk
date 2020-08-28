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

func (ref *Class) Object(id string) *ObjectRef {
	return &ObjectRef{
		c:     ref.c,
		class: ref.Name,
		ID:    id,
	}
}

func (ref *Class) Create(data interface{}, authOptions ...AuthOption) (*ObjectRef, error) {
	objectRef := new(ObjectRef)
	if err := objectCreate(ref, data, objectRef, authOptions...); err != nil {
		return nil, err
	}

	return objectRef, nil
}

func (ref *Class) NewQuery() *Query {
	return &Query{
		c:     ref.c,
		class: ref,
	}
}

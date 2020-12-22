package leancloud

type Class struct {
	c    *Client
	Name string
}

func (ref *Class) ID(id string) *ObjectRef {
	return &ObjectRef{
		c:     ref.c,
		class: ref.Name,
		ID:    id,
	}
}

func (ref *Class) Create(object interface{}, authOptions ...AuthOption) (*ObjectRef, error) {
	newRef, err := objectCreate(ref, object, authOptions...)
	if err != nil {
		return nil, err
	}

	objectRef, _ := newRef.(*ObjectRef)
	return objectRef, nil
}

func (ref *Class) NewQuery() *Query {
	return &Query{
		c:     ref.c,
		class: ref,
		where: make(map[string]interface{}),
	}
}

package leancloud

type Class struct {
	c    *Client
	Name string
}

// ID constructs reference with objectId and className
func (ref *Class) ID(id string) *ObjectRef {
	return &ObjectRef{
		c:     ref.c,
		class: ref.Name,
		ID:    id,
	}
}

// Create write the Object to the Storage from the custom structure/bare Object/map.
func (ref *Class) Create(object interface{}, authOptions ...AuthOption) (*ObjectRef, error) {
	newRef, err := objectCreate(ref, object, authOptions...)
	if err != nil {
		return nil, err
	}

	objectRef, _ := newRef.(*ObjectRef)

	return objectRef, nil
}

// NewQuery constructs a new Query for general Class
func (ref *Class) NewQuery() *Query {
	return &Query{
		c:     ref.c,
		class: ref,
		where: make(map[string]interface{}),
	}
}

package leancloud

import (
	"fmt"
	"reflect"
)

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
	decodedRef, err := objectCreate(ref, data, authOptions...)
	if err != nil {
		return nil, err
	}

	objectRef, ok := decodedRef.(*ObjectRef)
	if !ok {
		return nil, fmt.Errorf("unexpected error when parse ObjectRef from response: want type *ObjectRef but %v", reflect.TypeOf(decodedRef))
	}

	return objectRef, nil
}

func (ref *Class) NewQuery() *Query {
	return &Query{
		c:     ref.c,
		class: ref,
		where: make(map[string]interface{}),
	}
}

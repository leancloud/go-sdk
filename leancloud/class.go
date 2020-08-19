package leancloud

import (
	"encoding/json"
	"errors"
	"fmt"
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
	method := methodPost
	path := fmt.Sprint("/1.1/classes/", ref.Name)

	options := ref.c.getRequestOptions()
	options.JSON = data

	resp, err := ref.c.request(ServiceAPI, method, path, options, authOptions...)
	if err != nil {
		return nil, err
	}

	respJSON := map[string]interface{}{}
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return nil, err
	}

	objectID, ok := respJSON["objectId"].(string)
	if !ok {
		return nil, errors.New("unable to fetch objectId from response")
	}
	return &ObjectRef{
		c:     ref.c,
		class: ref.Name,
		ID:    objectID,
	}, nil
}

func (ref *Class) NewQuery() *Query {
	return &Query{
		c:        ref.c,
		classRef: ref,
	}
}

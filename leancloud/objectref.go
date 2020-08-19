package leancloud

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ObjectRef struct {
	c     *Client
	class string
	ID    string
}

func (client *Client) Object(name, id string) *ObjectRef {
	return &ObjectRef{
		c:     client,
		class: name,
		ID:    id,
	}
}

func (r *ObjectRef) Get(authOption ...AuthOption) (*Object, error) {
	method := methodGet
	path := fmt.Sprint("/1.1/classes/", r.class, "/", r.ID)
	options := r.c.getRequestOptions()
	resp, err := r.c.request(ServiceAPI, method, path, options, authOption...)
	if err != nil {
		return nil, err
	}
	resBody := make(map[string]interface{})

	if err := json.Unmarshal(resp.Bytes(), &resBody); err != nil {
		return nil, err
	}

	return &Object{
		fields: resBody,
	}, nil
}

func (r *ObjectRef) Set(field string, value interface{}, authOption ...AuthOption) error {
	method := methodPut
	path := fmt.Sprint("/1.1/classes/", r.class)

	if r.ID != "" {
		method = methodPost
		path = fmt.Sprint(path, "/", r.ID)
	}

	options := r.c.getRequestOptions()
	options.JSON = map[string]interface{}{
		field: value,
	}

	resp, err := r.c.request(ServiceAPI, method, path, options, authOption...)

	if err != nil {
		return err
	}

	respJSON := map[string]interface{}{}
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return err
	}

	if r.ID == "" {
		objectID, ok := respJSON["objectId"].(string)
		if !ok {
			return errors.New("unable to fetch object id from response")
		}
		r.ID = objectID
	}

	return nil
}

func (r *ObjectRef) Update(data map[string]interface{}, authOption ...AuthOption) error {
	method := methodPut
	path := fmt.Sprint("/1.1/classes/", r.class)

	if r.ID != "" {
		method = methodPut
		path = fmt.Sprint(path, "/", r.ID)
	}

	options := r.c.getRequestOptions()
	options.JSON = data

	resp, err := r.c.request(ServiceAPI, method, path, options, authOption...)

	if err != nil {
		return err
	}

	respJSON := map[string]interface{}{}
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return err
	}

	if r.ID == "" {
		objectID, ok := respJSON["objectId"].(string)
		if !ok {
			return errors.New("unable to fetch object id from response")
		}
		r.ID = objectID
	}

	return nil
}

func (r *ObjectRef) UpdateWithQuery(data map[string]interface{}, query *Query, authOption ...AuthOption) error {
	// TODO
	return nil
}

func (r *ObjectRef) Destroy(authOption ...AuthOption) error {
	if r.ID == "" {
		return errors.New("cannot destroy nonexist object")
	}
	method := methodDelete
	path := fmt.Sprint("/1.1/classes/", r.class, "/", r.ID)
	options := r.c.getRequestOptions()

	_, err := r.c.request(ServiceAPI, method, path, options, authOption...)
	if err != nil {
		return err
	}

	return nil
}

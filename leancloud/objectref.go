package leancloud

import (
	"encoding/json"
	"fmt"
	"time"
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

func (ref *ObjectRef) Get(authOptions ...AuthOption) (*Object, error) {
	object := new(Object)
	if err := objectGet(ref, object, authOptions...); err != nil {
		return nil, err
	}

	return object, nil
}

func (ref *ObjectRef) Set(field string, value interface{}, authOptions ...AuthOption) error {
	if ref.ID == "" {
		return fmt.Errorf("no reference to object")
	}

	if err := objectSet(ref, field, value, authOptions...); err != nil {
		return err
	}

	return nil
}

func (ref *ObjectRef) Update(data map[string]interface{}, authOptions ...AuthOption) error {
	if ref.ID == "" {
		return fmt.Errorf("no reference to object")
	}

	if err := objectUpdate(ref, data, authOptions...); err != nil {
		return err
	}

	return nil
}

func (ref *ObjectRef) UpdateWithQuery(data map[string]interface{}, query *Query, authOptions ...AuthOption) error {
	// TODO
	return nil
}

func (ref *ObjectRef) Destroy(authOptions ...AuthOption) error {
	if ref.ID == "" {
		return fmt.Errorf("no reference to object")
	}

	if err := objectDestroy(ref); err != nil {
		return err
	}

	return nil
}

func objectGet(ref interface{}, object interface{}, authOptions ...AuthOption) error {
	path := "/1.1/"
	id := ""
	var c *Client

	switch ref.(type) {
	case *ObjectRef:
		objectRef := ref.(*ObjectRef)
		path = fmt.Sprint(path, "classes/", objectRef.class, "/", objectRef.ID)
		id = objectRef.ID
		c = objectRef.c
		break
	case *UserRef:
		userRef := ref.(*UserRef)
		path = fmt.Sprint(path, "users/", userRef.ID)
		id = userRef.ID
		c = userRef.c
		break
	}

	resp, err := c.request(ServiceAPI, methodGet, path, c.getRequestOptions(), authOptions...)
	if err != nil {
		return err
	}

	respJSON := make(map[string]interface{})

	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return err
	}

	createdAt, err := time.Parse(time.RFC3339, respJSON["createdAt"].(string))
	if err != nil {
		return fmt.Errorf("unable to parse createdAt from response %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, respJSON["updatedAt"].(string))
	if err != nil {
		return fmt.Errorf("unable to parse updatedAt from response %w", err)
	}

	switch ref.(type) {
	case *ObjectRef:
		object := object.(*Object)
		object.ID = id
		object.CreatedAt = createdAt
		object.UpdatedAt = updatedAt
		object.fields = respJSON
		break
	case *UserRef:
		sessionToken := respJSON["sessionToken"].(string)
		user := object.(*User)
		user.ID = id
		user.CreatedAt = createdAt
		user.UpdatedAt = updatedAt
		user.sessionToken = sessionToken
		user.fields = respJSON
		break
	}

	return nil
}

func objectSet(ref interface{}, field string, data interface{}, authOptions ...AuthOption) error {
	path := "/1.1/"
	var c *Client

	switch ref.(type) {
	case *ObjectRef:
		objectRef := ref.(*ObjectRef)
		path = fmt.Sprint(path, "classes/", objectRef.class, "/", objectRef.ID)
		c = objectRef.c
		break
	case *UserRef:
		userRef := ref.(*UserRef)
		path = fmt.Sprint(path, "classes/users/", userRef.ID)
		c = userRef.c
		break
	}

	options := c.getRequestOptions()
	options.JSON = encodeObject(map[string]interface{}{
		field: data,
	})

	resp, err := c.request(ServiceAPI, methodPut, path, options, authOptions...)
	if err != nil {
		return err
	}

	respJSON := make(map[string]interface{})
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return err
	}

	return nil
}

func objectUpdate(ref interface{}, data map[string]interface{}, authOptions ...AuthOption) error {
	path := "/1.1/"
	var c *Client

	switch ref.(type) {
	case *ObjectRef:
		objectRef := ref.(*ObjectRef)
		path = fmt.Sprint(path, "classes/", objectRef.class, "/", objectRef.ID)
		c = objectRef.c
		break
	case *UserRef:
		userRef := ref.(*UserRef)
		path = fmt.Sprint(path, "classes/users/", userRef.ID)
		c = userRef.c
		break
	}

	options := c.getRequestOptions()
	options.JSON = encodeObject(data)

	resp, err := c.request(ServiceAPI, methodPut, path, options, authOptions...)
	if err != nil {
		return err
	}

	respJSON := make(map[string]interface{})
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return err
	}

	return nil
}

func objectDestroy(ref interface{}, authOptions ...AuthOption) error {
	path := "/1.1/"
	var c *Client

	switch ref.(type) {
	case *ObjectRef:
		objectRef := ref.(*ObjectRef)
		path = fmt.Sprint(path, "classes/", objectRef.class, "/", objectRef.ID)
		c = objectRef.c
		break
	case *UserRef:
		userRef := ref.(*UserRef)
		path = fmt.Sprint(path, "classes/users/", userRef.ID)
		c = userRef.c
		break
	}

	resp, err := c.request(ServiceAPI, methodDelete, path, c.getRequestOptions(), authOptions...)
	if err != nil {
		return err
	}

	respJSON := make(map[string]interface{})
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return err
	}

	return nil
}

package leancloud

import (
	"encoding/json"
	"errors"
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
	path := fmt.Sprint("/1.1/classes/", ref.class, "/", ref.ID)

	resp, err := ref.c.request(ServiceAPI, methodGet, path, ref.c.getRequestOptions(), authOptions...)
	if err != nil {
		return nil, err
	}
	resBody := make(map[string]interface{})

	if err := json.Unmarshal(resp.Bytes(), &resBody); err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, resBody["createdAt"].(string))
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339, resBody["updatedAt"].(string))
	if err != nil {
		return nil, err
	}

	return &Object{
		ID:        resBody["objectId"].(string),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		fields:    resBody,
	}, nil
}

func (ref *ObjectRef) Set(field string, value interface{}, authOptions ...AuthOption) error {
	if ref.ID == "" {
		return errors.New("no reference to object")
	}

	path := fmt.Sprint("/1.1/classes/", ref.class, "/", ref.ID)
	options := ref.c.getRequestOptions()

	options.JSON = encodeObject(map[string]interface{}{
		field: value,
	})

	resp, err := ref.c.request(ServiceAPI, methodPut, path, options, authOptions...)

	if err != nil {
		return err
	}

	respJSON := map[string]interface{}{}
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return err
	}

	return nil
}

func (ref *ObjectRef) Update(data map[string]interface{}, authOptions ...AuthOption) error {
	path := fmt.Sprint("/1.1/classes/", ref.class, "/", ref.ID)

	options := ref.c.getRequestOptions()
	options.JSON = encodeObject(data)

	resp, err := ref.c.request(ServiceAPI, methodPut, path, options, authOptions...)

	if err != nil {
		return err
	}

	respJSON := map[string]interface{}{}
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return err
	}

	if ref.ID == "" {
		objectID, ok := respJSON["objectId"].(string)
		if !ok {
			return errors.New("unable to fetch object id from response")
		}
		ref.ID = objectID
	}

	return nil
}

func (ref *ObjectRef) UpdateWithQuery(data map[string]interface{}, query *Query, authOptions ...AuthOption) error {
	// TODO
	return nil
}

func (ref *ObjectRef) Destroy(authOptions ...AuthOption) error {
	if ref.ID == "" {
		return errors.New("cannot destroy nonexist object")
	}
	path := fmt.Sprint("/1.1/classes/", ref.class, "/", ref.ID)

	_, err := ref.c.request(ServiceAPI, methodDelete, path, ref.c.getRequestOptions(), authOptions...)
	if err != nil {
		return err
	}

	return nil
}

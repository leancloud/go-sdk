package leancloud

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/levigross/grequests"
)

type ObjectRef struct {
	c     *Client
	class string
	ID    string
}

func (client *Client) Object(object interface{}) *ObjectRef {
	if meta := extractObjectMeta(object); meta != nil {
		return meta.ref.(*ObjectRef)
	}

	return nil
}

// Get fetchs object from backend
func (ref *ObjectRef) Get(object interface{}, authOptions ...AuthOption) error {
	if ref == nil || ref.ID == "" || ref.class == "" {
		return nil
	}

	if err := objectGet(ref, object, authOptions...); err != nil {
		return err
	}

	return nil
}

// Set manipulate
func (ref *ObjectRef) Set(key string, value interface{}, authOptions ...AuthOption) error {
	if ref == nil || ref.ID == "" || ref.class == "" {
		return nil
	}

	if err := objectSet(ref, key, value, authOptions...); err != nil {
		return err
	}

	return nil
}

func (ref *ObjectRef) Update(diff interface{}, authOptions ...AuthOption) error {
	if ref == nil || ref.ID == "" || ref.class == "" {
		return nil
	}

	if err := objectUpdate(ref, diff, authOptions...); err != nil {
		return err
	}

	return nil
}

func (ref *ObjectRef) UpdateWithQuery(data map[string]interface{}, query *Query, authOptions ...AuthOption) error {
	// TODO
	return nil
}

func (ref *ObjectRef) Destroy(authOptions ...AuthOption) error {
	if ref == nil || ref.ID == "" || ref.class == "" {
		return nil
	}

	if err := objectDestroy(ref, authOptions...); err != nil {
		return err
	}

	return nil
}

func objectCreate(class interface{}, object interface{}, authOptions ...AuthOption) (interface{}, error) {
	path := "/1.1/"
	var c *Client
	var options *grequests.RequestOptions

	switch v := class.(type) {
	case *Class:
		path = fmt.Sprint(path, "classes/", v.Name)
		c = v.c
		options = c.getRequestOptions()
		switch reflect.ValueOf(object).Kind() {
		case reflect.Map:
			options.JSON = encodeMap(object, false)
		case reflect.Struct:
			options.JSON = encodeMap(object, false)
		default:
			return nil, fmt.Errorf("")
		}
		break
	case *Users:
		path = fmt.Sprint(path, "users")
		c = v.c
		options = c.getRequestOptions()
		switch reflect.ValueOf(object).Kind() {
		case reflect.Map:
			options.JSON = encodeMap(object, false)
		case reflect.Struct:
			options.JSON = encodeUser(object, false, false)
		default:
			return nil, fmt.Errorf("")
		}
		break
	}

	resp, err := c.request(ServiceAPI, MethodPost, path, options, authOptions...)
	if err != nil {
		return nil, err
	}

	respJSON := make(map[string]interface{})
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return nil, err
	}
	switch v := class.(type) {
	case *Class:
		objectID, ok := respJSON["objectId"].(string)
		if !ok {
			return nil, fmt.Errorf("unexpected error when parse objectId from response: want type string but %v", reflect.TypeOf(respJSON["objectId"]))
		}
		return &ObjectRef{
			ID:    objectID,
			class: v.Name,
			c:     c,
		}, nil
	case *Users:
		return decodeUser(respJSON)
	}

	return nil, nil

}

func objectGet(ref interface{}, object interface{}, authOptions ...AuthOption) error {
	path := "/1.1/"
	var c *Client

	switch v := ref.(type) {
	case *ObjectRef:
		path = fmt.Sprint(path, "classes/", v.class, "/", v.ID)
		c = v.c
		break
	case *UserRef:
		path = fmt.Sprint(path, "users/", v.ID)
		c = v.c
		break
	case *FileRef:
		path = fmt.Sprint(path, "files/", v.ID)
		c = v.c
	}

	resp, err := c.request(ServiceAPI, MethodGet, path, c.getRequestOptions(), authOptions...)
	if err != nil {
		return err
	}

	respJSON := make(map[string]interface{})
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return err
	}

	switch v := ref.(type) {
	case *ObjectRef:
		decodedObject, err := decodeObject(respJSON)
		if err != nil {
			return err
		}
		decodedObject.ref = v
		return bind(reflect.ValueOf(decodedObject), reflect.ValueOf(object))
	case *UserRef:
		decodedUser, err := decodeUser(respJSON)
		if err != nil {
			return err
		}
		decodedUser.ref = v
		return bind(reflect.ValueOf(decodedUser), reflect.ValueOf(object))
	case *FileRef:
		decodedFile, err := decodeFile(respJSON)
		if err != nil {
			return err
		}
		decodedFile.ref = v
		object = decodedFile
	}

	return nil
}

func objectSet(ref interface{}, key string, value interface{}, authOptions ...AuthOption) error {
	path := "/1.1/"
	var c *Client

	switch v := ref.(type) {
	case *ObjectRef:
		path = fmt.Sprint(path, "classes/", v.class, "/", v.ID)
		c = v.c
		break
	case *UserRef:
		path = fmt.Sprint(path, "users/", v.ID)
		c = v.c
		break
	}

	options := c.getRequestOptions()
	options.JSON = encode(map[string]interface{}{key: value}, true)

	_, err := c.request(ServiceAPI, MethodPut, path, options, authOptions...)
	if err != nil {
		return err
	}

	return nil
}

func objectUpdate(ref interface{}, diff interface{}, authOptions ...AuthOption) error {
	path := "/1.1/"
	var c *Client
	var options *grequests.RequestOptions

	switch v := ref.(type) {
	case *ObjectRef:
		path = fmt.Sprint(path, "classes/", v.class, "/", v.ID)
		c = v.c
		options = c.getRequestOptions()
		switch reflect.ValueOf(diff).Kind() {
		case reflect.Map:
			options.JSON = encodeMap(diff, true)
		case reflect.Struct:
			options.JSON = encodeObject(diff, false, true)
		default:
			return fmt.Errorf("")
		}
		break
	case *UserRef:
		path = fmt.Sprint(path, "users/", v.ID)
		c = v.c
		options = c.getRequestOptions()
		switch reflect.ValueOf(diff).Kind() {
		case reflect.Map:
			options.JSON = encodeMap(diff, true)
		case reflect.Struct:
			options.JSON = encodeUser(diff, false, true)
		default:
			return fmt.Errorf("")
		}
		break
	}

	_, err := c.request(ServiceAPI, MethodPut, path, options, authOptions...)
	if err != nil {
		return err
	}

	return nil
}

func objectDestroy(ref interface{}, authOptions ...AuthOption) error {
	path := "/1.1/"
	var c *Client

	switch v := ref.(type) {
	case *ObjectRef:
		path = fmt.Sprint(path, "classes/", v.class, "/", v.ID)
		c = v.c
	case *UserRef:
		path = fmt.Sprint(path, "users/", v.ID)
		c = v.c
	case *FileRef:
		path = fmt.Sprint(path, "files/", v.ID)
		c = v.c
	}

	resp, err := c.request(ServiceAPI, MethodDelete, path, c.getRequestOptions(), authOptions...)
	if err != nil {
		return err
	}

	respJSON := make(map[string]interface{})
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return err
	}

	return nil
}

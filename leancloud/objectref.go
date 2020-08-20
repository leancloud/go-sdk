package leancloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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

	options.JSON = objectSerialize(map[string]interface{}{
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

	if ref.ID == "" {
		objectID, ok := respJSON["objectId"].(string)
		if !ok {
			return errors.New("unable to fetch object id from response")
		}
		ref.ID = objectID
	}

	return nil
}

func (ref *ObjectRef) Update(data interface{}, authOptions ...AuthOption) error {
	method := methodPut
	path := fmt.Sprint("/1.1/classes/", ref.class)

	if ref.ID != "" {
		method = methodPut
		path = fmt.Sprint(path, "/", ref.ID)
	}

	options := ref.c.getRequestOptions()
	options.JSON = objectSerialize(data)

	resp, err := ref.c.request(ServiceAPI, method, path, options, authOptions...)

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

func objectSerialize(object interface{}) map[string]interface{} {
	mapObject := make(map[string]interface{})
	if reflect.TypeOf(object).Kind() == reflect.Struct {
		v := reflect.ValueOf(object)
		s := reflect.TypeOf(object)
		for i := 0; i < v.NumField(); i++ {
			switch v.Field(i).Type() {
			case reflect.TypeOf(time.Time{}):
				date := v.Field(i).Interface().(time.Time)
				mapObject[s.Field(i).Tag.Get("json")] = map[string]interface{}{
					"__type": "Date",
					"iso":    fmt.Sprint(date.In(time.FixedZone("UTC", 0)).Format("2006-01-02T15:04:05.000Z")),
				}
				break
			default:
				mapObject[s.Field(i).Tag.Get("json")] = v.Field(i).Interface()
			}
		}
	}
	if reflect.TypeOf(object).Kind() == reflect.Map {
		iter := reflect.ValueOf(object).MapRange()
		for iter.Next() {
			switch iter.Value().Elem().Type() {
			case reflect.TypeOf(time.Time{}):
				date := iter.Value().Interface().(time.Time)
				mapObject[iter.Key().String()] = map[string]interface{}{
					"__type": "Date",
					"iso":    fmt.Sprint(date.In(time.FixedZone("UTC", 0)).Format("2006-01-02T15:04:05.000Z")),
				}
				break
			default:
				mapObject[iter.Key().String()] = iter.Value().Interface()
			}
		}
	}

	return mapObject
}

package lean

import (
	"encoding/json"
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

func (r *ObjectRef) Get(auth ...AuthOption) (*Object, error) {
	method := methodGet
	path := fmt.Sprint("/1.1/classes/", r.class, "/", r.ID)
	options := r.c.getRequestOptions()
	resp, err := r.c.request(ServiceAPI, method, path, options, auth...)
	if err != nil {
		return nil, err
	}
	resBody := make(map[string]interface{})

	if err := json.Unmarshal(resp.Bytes(), &resBody); err != nil {
		return nil, err
	}

	return &Object{
		fields: &resBody,
	}, nil
}

func (r *ObjectRef) Set(field string, value interface{}, authOption ...AuthOption) error {
	// TODO
	return nil
}

func (r *ObjectRef) Update(data map[string]interface{}, authOption ...AuthOption) error {
	// TODO
	return nil
}

func (r *ObjectRef) UpdateWithQuery(data map[string]interface{}, query *Query, authOption ...AuthOption) error {
	return nil
}

func (r *ObjectRef) Delete() error {
	// TODO
	return nil
}

func OpIncrement(amount int) map[string]interface{} {
	op := make(map[string]interface{})

	op["__op"] = "Increment"
	op["amount"] = amount

	return op
}

func OpDecrement(amount int) map[string]interface{} {
	op := make(map[string]interface{})

	op["__op"] = "Decrement"
	op["amount"] = amount

	return op
}

func OpAdd(objects interface{}) map[string]interface{} {
	op := make(map[string]interface{})

	op["__op"] = "Add"
	op["objects"] = objects

	return op
}

func OpAddUnique(objects interface{}) map[string]interface{} {
	op := make(map[string]interface{})

	op["__op"] = "AddUnique"
	op["objects"] = objects

	return op
}

func OpAddRelation() {
	// TODO after Pointer implementation
}

func OpRemove(objects interface{}) map[string]interface{} {
	op := make(map[string]interface{})

	op["__op"] = "Remove"
	op["objects"] = objects

	return op
}

func OpRemoveRelation() {
	// TODO after Pointer implementation
}

func OpDelete(delete bool) map[string]interface{} {
	op := make(map[string]interface{})

	op["__op"] = "Delete"
	op["delete"] = delete

	return op
}

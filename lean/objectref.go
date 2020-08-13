package lean

type ObjectRef struct {
	c     *Client
	class string
	ID    string
}

func (client *Client) Object(name, id string) *ObjectRef {
	// TODO
	return nil
}

func (r *ObjectRef) Get() (*Object, error) {
	// TODO
	return nil, nil
}

func (r *ObjectRef) Set(data interface{}) error {
	// TODO
	return nil
}

func (r *ObjectRef) Update(data []map[string]interface{}) error {
	// TODO
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

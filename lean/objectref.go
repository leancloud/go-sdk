package lean

type ObjectRef struct {
	c     *Client
	class string
	ID    string
}

type OpType string

const (
	BitAnd         OpType = "BitAnd"
	BitOr          OpType = "BitOr"
	BitXor         OpType = "BitXor"
	Increment      OpType = "Increment"
	Decrement      OpType = "Decrement"
	Add            OpType = "Add"
	AddUnique      OpType = "AddUnique"
	AddRelation    OpType = "AddRelation"
	Remove         OpType = "Remove"
	RemoveRelation OpType = "RemoveRelation"
	Delete         OpType = "Delete"
)

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

func Op(op OpType, data interface{}) map[string]interface{} {
	opObject := make(map[string]interface{})
	opObject["__op"] = op

	switch op {
	case Increment:
	case Decrement:
		opObject["amount"] = data
		break
	case BitAnd:
	case BitOr:
	case BitXor:
		opObject["value"] = data
	case Add:
	case AddUnique:
	case AddRelation:
	case Remove:
	case RemoveRelation:
		opObject["objects"] = data
		break
	case Delete:
		opObject["delete"] = data
		break
	}

	return opObject
}

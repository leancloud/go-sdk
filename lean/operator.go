package lean

type Operator interface {
	MarshalJSON() ([]byte, error)
}

type op struct {
	Name    string      `json:"__op"`
	Amount  *int        `json:"amount,omitempty"`
	Objects interface{} `json:",omitempty"`
}

func (op op) MarshalJSON() ([]byte, error) {
	// TODO
	return nil, nil
}

func OpIncrement(amount int) Operator {
	return op{
		Name:   "Increment",
		Amount: &amount,
	}
}

func OpAdd(objects interface{}) Operator {
	return op{
		Name:    "Add",
		Objects: objects,
	}
}

func OpRemove(objects interface{}) Operator {
	return op{
		Name:    "Remove",
		Objects: objects,
	}
}

func OpDelete() Operator {
	return op{
		Name: "Delete",
		Objects: map[string]bool{
			"delete": true,
		},
	}
}

package lean

type Operator interface {
	MarshalJSON() ([]byte, error)
}

type Op struct {
	Name    string      `json:"__op"`
	Amount  *int        `json:"amount,omitempty"`
	Objects interface{} `json:",omitempty"`
}

func (op Op) MarshalJSON() ([]byte, error) {
	// TODO
	return nil, nil
}

func OpIncrement(amount int) Operator {
	if amount > 0 {
		return Op{
			Name:   "Increment",
			Amount: &amount,
		}
	}

	positive := -amount
	return Op{
		Name:   "Decrement",
		Amount: &positive,
	}
}

func OpAdd(objects interface{}) Operator {
	return Op{
		Name:    "Add",
		Objects: objects,
	}
}

func OpRemove(objects interface{}) Operator {
	return Op{
		Name:    "Remove",
		Objects: objects,
	}
}

func OpDelete() Operator {
	return Op{
		Name: "Delete",
		Objects: map[string]bool{
			"delete": true,
		},
	}
}

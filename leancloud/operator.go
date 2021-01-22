package leancloud

type Op struct {
	name    string
	objects interface{}
}

func OpIncrement(amount interface{}) Op {
	return Op{
		name:    "Increment",
		objects: amount,
	}
}

func OpDecrement(amount interface{}) Op {
	return Op{
		name:    "Decrement",
		objects: amount,
	}
}
func OpAdd(objects interface{}) Op {
	return Op{
		name:    "Add",
		objects: objects,
	}
}

func OpAddUnique(objects interface{}) Op {
	return Op{
		name:    "AddUnique",
		objects: objects,
	}
}
func OpRemove(objects interface{}) Op {
	return Op{
		name:    "Remove",
		objects: objects,
	}
}

func OpDelete() Op {
	return Op{
		name: "Delete",
	}
}

func OpAddRelation(objects interface{}) Op {
	// TODO
	return Op{}
}

func OpRemoveRelation(objects interface{}) Op {
	// TODO
	return Op{}
}

func OpBitAnd(value interface{}) Op {
	return Op{
		name:    "BitAnd",
		objects: value,
	}
}

func OpBitOr(value interface{}) Op {
	return Op{
		name:    "BitAnd",
		objects: value,
	}
}

func OpBitXor(value interface{}) Op {
	return Op{
		name:    "BitAnd",
		objects: value,
	}
}

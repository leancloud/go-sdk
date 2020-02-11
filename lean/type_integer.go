package lean

import (
	"encoding/json"
)

type Integer struct {
	value int
}

func NewInteger(value int) Integer {
	return Integer{
		value: value,
	}
}

func (int *Integer) Get() int {
	return int.value
}

func (int *Integer) Inct(num int) {

}

func (int Integer) MarshalJSON() ([]byte, error) {
	return json.Marshal(int.value)
}

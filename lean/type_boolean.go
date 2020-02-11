package lean

import (
	"encoding/json"
	"fmt"
)

type Boolean struct {
	value *bool
}

func NewBoolean(value bool) Boolean {
	return Boolean{
		value: &value,
	}
}

func (bool *Boolean) Get() bool {
	if bool.value == nil {
		return false
	}

	return *bool.value
}

func (bool Boolean) String() string {
	return fmt.Sprint(bool.Get())
}

func (bool *Boolean) IsNull() bool {
	return bool.value == nil
}

func (bool Boolean) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool.value)
}

func (boolean *Boolean) UnmarshalJSON(bytes []byte) error {
	var b bool

	err := json.Unmarshal(bytes, &b)

	if err != nil {
		return err
	}

	boolean.value = &b

	return nil
}

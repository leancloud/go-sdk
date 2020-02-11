package lean

import (
	"encoding/json"
)

type String struct {
	value *string
}

func NewString(value string) String {
	return String{
		value: &value,
	}
}

func (string *String) Get() string {
	if string.value == nil {
		return ""
	}

	return *string.value
}

func (string String) String() string {
	return string.Get()
}

func (string *String) IsNull() bool {
	return string.value == nil
}

func (string String) MarshalJSON() ([]byte, error) {
	return json.Marshal(string.value)
}

func (str *String) UnmarshalJSON(bytes []byte) error {
	var s string

	err := json.Unmarshal(bytes, &s)

	if err != nil {
		return err
	}

	str.value = &s

	return nil
}

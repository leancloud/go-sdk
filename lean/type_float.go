package lean

import (
	"encoding/json"
)

type Float struct {
	value float64
}

func NewFloat(value float64) Float {
	return Float{
		value: value,
	}
}

func (float *Float) Get() float64 {
	return float.value
}

func (float Float) MarshalJSON() ([]byte, error) {
	return json.Marshal(float.value)
}

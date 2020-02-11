package lean

import (
	"encoding/json"
	"time"
)

type Date struct {
	value time.Time
}

func NewDate(value time.Time) Date {
	return Date{
		value: value,
	}
}

func (date *Date) Get() time.Time {
	return date.value
}

func (date Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"__type": "Date",
		"iso":    date.value.In(time.FixedZone("UTC", 0)).Format("2006-01-02T15:04:05.000Z"),
	})
}

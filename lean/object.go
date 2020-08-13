package lean

import "time"

type Object struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Fields    *map[string]interface{}
}

func (ss *Object) Data() map[string]interface{} {
	// TODO
	return nil
}

func (ss *Object) Struct(p interface{}) error {
	// TODO
	return nil
}

func (ss *Object) Get(filed string) interface{} {
	// TODO
	return nil
}

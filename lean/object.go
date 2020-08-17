package lean

import "time"

type Object struct {
	ID        string    `json:"objectId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	fields    *map[string]interface{}
}

func (ss *Object) GetMap() map[string]interface{} {
	// TODO
	return nil
}

func (ss *Object) ToStruct(p interface{}) error {
	// TODO
	return nil
}

func (ss *Object) Get(filed string) interface{} {
	// TODO
	return nil
}

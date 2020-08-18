package leancloud

import "time"

type Object struct {
	ID        string    `json:"objectId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	fields    map[string]interface{}
}

func (object *Object) GetMap() map[string]interface{} {
	// TODO
	return nil
}

func (object *Object) ToStruct(p interface{}) error {
	// TODO
	return nil
}

func (object *Object) Get(filed string) interface{} {
	// TODO
	return nil
}

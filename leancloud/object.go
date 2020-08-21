package leancloud

import (
	"time"
)

type Object struct {
	ID        string    `json:"objectId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	fields    map[string]interface{}
}

func (object *Object) GetMap() map[string]interface{} {
	return decodeFields(object.fields)
}

func (object *Object) ToStruct(p interface{}) {
	decodeObject(object.fields, p)
}

func (object *Object) Get(field string) interface{} {
	return object.fields[field]
}

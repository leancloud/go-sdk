package leancloud

import (
	"time"
)

type Object struct {
	ID        string    `json:"objectId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	fields    map[string]interface{}
	isPointer bool
}

func (object *Object) GetMap() map[string]interface{} {
	return object.fields
}

func (object *Object) ToStruct(p interface{}) {
	_ = decodeObject(object.fields, p)
}

func (object *Object) Get(field string) interface{} {
	return object.fields[field]
}

func (object *Object) IsPointer() bool {
	return object.isPointer
}

func (object *Object) Included() bool {
	if object.isPointer {
		if len(object.fields) != 0 {
			return true
		}
	}

	return false
}
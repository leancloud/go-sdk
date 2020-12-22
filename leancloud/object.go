package leancloud

import (
	"reflect"
	"time"
)

type Object struct {
	ID         string    `json:"objectId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	fields     map[string]interface{}
	isPointer  bool
	isIncluded bool
	ref        interface{}
}

func (object *Object) Raw() map[string]interface{} {
	return object.fields
}

func (object *Object) Get(key string) interface{} {
	return object.fields[key]
}

func (object *Object) Int(key string) int64 {
	return reflect.ValueOf(object.fields[key]).Int()
}

func (object *Object) String(key string) string {
	return reflect.ValueOf(object.fields[key]).String()
}

func (object *Object) Float(key string) float64 {
	return reflect.ValueOf(object.fields[key]).Float()
}

func (object *Object) Bool(key string) bool {
	return reflect.ValueOf(object.fields[key]).Bool()
}

func (object *Object) GeoPoint(key string) *GeoPoint {
	pointPtr, ok := object.fields[key].(*GeoPoint)
	if !ok {
		point, ok := object.fields[key].(GeoPoint)
		if !ok {
			return nil
		}
		return &point
	}
	return pointPtr
}

func (object *Object) Date(key string) *time.Time {
	datePtr, ok := object.fields[key].(*time.Time)
	if !ok {
		date, ok := object.fields[key].(time.Time)
		if !ok {
			return nil
		}
		return &date
	}
	return datePtr
}

func (object *Object) File(key string) *File {
	filePtr, ok := object.fields[key].(*File)
	if !ok {
		file, ok := object.fields[key].(File)
		if !ok {
			return nil
		}
		return &file
	}
	return filePtr
}

func (object *Object) Bytes(key string) []byte {
	bytes, ok := object.fields[key].([]byte)
	if !ok {
		return nil
	}
	return bytes
}

func (object *Object) ACL() *ACL {
	aclPtr, ok := object.fields["ACL"].(*ACL)
	if !ok {
		acl, ok := object.fields["ACL"].(ACL)
		if !ok {
			return nil
		}
		return &acl
	}
	return aclPtr
}

func (object *Object) IsPointer() bool {
	return object.isPointer
}

func (object *Object) Included() bool {
	return object.isIncluded
}

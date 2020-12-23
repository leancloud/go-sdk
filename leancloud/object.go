package leancloud

import (
	"reflect"
	"time"
)

// Object contains full data of Object.
// Also Object could be metadata for custom structure
type Object struct {
	ID         string    `json:"objectId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	fields     map[string]interface{}
	isPointer  bool
	isIncluded bool
	ref        interface{}
}

// Raw returns raw data of Object in form of map
func (object *Object) Raw() map[string]interface{} {
	return object.fields
}

// Get returns value of the key
func (object *Object) Get(key string) interface{} {
	return object.fields[key]
}

// Int returns int64 value of the key
func (object *Object) Int(key string) int64 {
	return reflect.ValueOf(object.fields[key]).Int()
}

// String returns string value of the key
func (object *Object) String(key string) string {
	return reflect.ValueOf(object.fields[key]).String()
}

// Float returns float64 value of the key
func (object *Object) Float(key string) float64 {
	return reflect.ValueOf(object.fields[key]).Float()
}

// Bool returns boolean value of the key
func (object *Object) Bool(key string) bool {
	return reflect.ValueOf(object.fields[key]).Bool()
}

// GeoPoint returns GeoPoint value of the key
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

// Date returns time.Time value of the key
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

// File returns File value of the key
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

// Bytes returns []byte value of the key
func (object *Object) Bytes(key string) []byte {
	bytes, ok := object.fields[key].([]byte)
	if !ok {
		return nil
	}
	return bytes
}

// ACL returns ACL value
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

// IsPointer shows whether the Object is a Pointer
func (object *Object) IsPointer() bool {
	return object.isPointer
}

// Included shows whether the Object is a included Pointer
func (object *Object) Included() bool {
	return object.isIncluded
}

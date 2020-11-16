package leancloud

import "reflect"

type User struct {
	Object
	sessionToken string
}

func (user *User) GetMap() map[string]interface{} {
	return user.fields
}

func (user *User) ToStruct(p interface{}) error {
	fv := reflect.ValueOf(user.fields)
	pv := reflect.Indirect(reflect.ValueOf(p))
	return bind(fv, pv)
}

func (user *User) Get(field string) interface{} {
	return user.fields[field]
}

func (user *User) GetSessionToken() string {
	return user.sessionToken
}

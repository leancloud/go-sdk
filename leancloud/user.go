package leancloud

type User struct {
	Object
	sessionToken string
}

func (user *User) GetMap() map[string]interface{} {
	return user.fields
}

func (user *User) ToStruct(p interface{}) error {
	return transform(user.fields, p)

}

func (user *User) Get(field string) interface{} {
	return user.fields[field]
}

func (user *User) GetSessionToken() string {
	return user.sessionToken
}

package lean

type User struct {
	Object
	sessionToken string
}

func (user *User) GetMap() map[string]interface{} {
	// TODO
	return nil
}

func (user *User) ToStruct(p interface{}) error {
	// TODO
	return nil
}

func (user *User) Get(field string) interface{} {
	// TODO
	return nil
}

func (user *User) GetSessionToken() string {
	return user.sessionToken
}

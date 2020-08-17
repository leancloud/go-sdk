package lean

type User struct {
	Fields map[string]interface{}
}

func (user *User) Map() map[string]interface{} {
	// TODO
	return nil
}

func (user *User) Struct(p interface{}) error {
	// TODO
	return nil
}

func (user *User) Get(field string) interface{} {
	// TODO
	return nil
}

package leancloud

type User struct {
	Object
	sessionToken string
}

func (user *User) SessionToken() string {
	return user.sessionToken
}

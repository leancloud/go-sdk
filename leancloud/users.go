package leancloud

type Users struct {
	c *Client
}

func (r *Users) LogIn(username, password string) (*User, error) {
	// TODO
	return nil, nil
}

func (r *Users) SignUp(username, password string) (*User, error) {
	// TODO
	return nil, nil
}

func (c *Client) NewUserQuery() *UserQuery {
	// TODO
	return nil
}

func (ref *Users) Become(sessionToken string) (*User, error) {
	// TODO
	return nil, nil
}

package lean

type UsersRef struct {
	c *Client
}

func (r *UsersRef) Login(username, password string) (*User, error) {
	// TODO
	return nil, nil
}

func (r *UsersRef) Signup(username, password string) (*User, error) {
	// TODO
	return nil, nil
}

func (c *Client) NewUserQuery() *UserQuery {
	// TODO
	return nil
}

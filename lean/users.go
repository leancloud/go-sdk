package lean

type UsersRef struct {
	c *Client
}

func (r *UsersRef) Login() (*User, error) {
	// TODO
	return nil, nil
}

func (r *UsersRef) Signup() (*User, error) {
	// TODO
	return nil, nil
}

func (c *Client) NewUserQuery() *UserQuery {
	// TODO
	return nil
}

package leancloud

type RoleRef struct {
	c     *Client
	class string
	ID    string
}

func (client *Client) Role(id string) *RoleRef {
	return &RoleRef{
		c:     client,
		class: "roles",
		ID:    id,
	}
}

func (ref *RoleRef) Get(authOption ...AuthOption) (*Role, error) {
	return nil, nil
}

func (ref *RoleRef) Set(field string, value interface{}, authOptions ...AuthOption) error {
	return nil
}

func (ref *RoleRef) Update(data map[string]interface{}, authOptions ...AuthOption) error {
	return nil
}

func (ref *RoleRef) UpdateWithQuery(data map[string]interface{}, authOptions ...AuthOption) error {
	// TODO
	return nil
}

func (ref *RoleRef) Destroy(authOptions ...AuthOption) error {
	return nil
}

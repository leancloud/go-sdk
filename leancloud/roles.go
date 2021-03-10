package leancloud

type Roles struct {
	c *Client
}

func (ref *Roles) NewQuery() *Query {
	return &Query{
		class: &Class{
			Name: "_Role",
			c:    ref.c,
		},
		c:     ref.c,
		where: make(map[string]interface{}),
	}
}

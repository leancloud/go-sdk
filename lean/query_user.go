package lean

type UserQuery struct {
	c        *Client
	classRef *Class
	where    map[string]string
	order    []string
	limit    int
	skip     int
}

func (r *UserQuery) Find() ([]User, error) {
	// TODO
	return nil, nil
}

func (r *UserQuery) First() (*User, error) {
	// TODO
	return nil, nil
}

func (r *UserQuery) Count() (int, error) {
	// TODO
	return 0, nil
}

func (r *UserQuery) Skip() *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) Limit() *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) Order() *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) EqualTo(key string, value string) *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) NotEqualTo(key string, value string) *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) SizeEqualTo() *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) Greater(key string, value string) *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) GreaterEqual() *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) Less(key string, value string) *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) LessEqual() *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) In(key string, value string) *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) NotIn(key string, value string) *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) All() *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) Exists() *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) Select() *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) Exclude() *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) Regexp() *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) Contains() *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) ContainsAll() *ObjectQuery {
	// TODO
	return nil
}

func (r *UserQuery) StartsWith() *ObjectQuery {
	// TODO
	return nil
}

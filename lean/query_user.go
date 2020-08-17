package lean

type UserQuery struct {
	c        *Client
	classRef *Class
	where    map[string]string
	order    []string
	limit    int
	skip     int
}

func (r *UserQuery) Find(auth ...AuthOption) ([]Object, error) {
	// TODO
	return nil, nil
}

func (r *UserQuery) First() (*Object, error) {
	// TODO
	return nil, nil
}

func (r *UserQuery) Count() (int, error) {
	// TODO
	return 0, nil
}

func (r *UserQuery) Skip(count int) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) Limit(limit int) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) Order(keys ...string) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) EqualTo(key string, value interface{}) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) NotEqualTo(key string, value interface{}) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) SizeEqualTo(key string, count int) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) GreaterThan(key string, value interface{}) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) GreaterThanEqualTo(key string, value interface{}) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) LessThan(key string, value interface{}) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) LessThanEqualTo(key string, value interface{}) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) In(key string, data interface{}) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) NotIn(key string, data interface{}) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) Regexp(expr, options string) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) Contains(key, substring string) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) ContainsAll(key string, objects interface{}) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) StartsWith(key, prefix string) *UserQuery {
	// TODO
	return nil
}

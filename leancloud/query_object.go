package leancloud

type Query struct {
	c        *Client
	classRef *Class
	where    map[string]string
	order    []string
	limit    int
	skip     int
}

func (r *Query) Find(auth ...AuthOption) ([]Object, error) {
	// TODO
	return nil, nil
}

func (r *Query) First() (*Object, error) {
	// TODO
	return nil, nil
}

func (r *Query) Count() (int, error) {
	// TODO
	return 0, nil
}

func (r *Query) Skip(count int) *Query {
	// TODO
	return nil
}

func (r *Query) Limit(limit int) *Query {
	// TODO
	return nil
}

func (r *Query) Order(keys ...string) *Query {
	// TODO
	return nil
}

func (r *Query) EqualTo(key string, value interface{}) *User {
	// TODO
	return nil
}

func (r *Query) NotEqualTo(key string, value interface{}) *Query {
	// TODO
	return nil
}

func (r *Query) SizeEqualTo(key string, count int) *Query {
	// TODO
	return nil
}

func (r *Query) GreaterThan(key string, value interface{}) *Query {
	// TODO
	return nil
}

func (r *Query) GreaterThanOrEqualTo(key string, value interface{}) *Query {
	// TODO
	return nil
}

func (r *Query) LessThan(key string, value interface{}) *Query {
	// TODO
	return nil
}

func (r *Query) LessThanOrEqualTo(key string, value interface{}) *Query {
	// TODO
	return nil
}

func (r *Query) In(key string, data interface{}) *Query {
	// TODO
	return nil
}

func (r *Query) NotIn(key string, data interface{}) *Query {
	// TODO
	return nil
}

func (r *Query) Regexp(expr, options string) *Query {
	// TODO
	return nil
}

func (r *Query) Contains(key, substring string) *Query {
	// TODO
	return nil
}

func (r *Query) ContainsAll(key string, objects interface{}) *Query {
	// TODO
	return nil
}

func (r *Query) StartsWith(key, prefix string) *Query {
	// TODO
	return nil
}

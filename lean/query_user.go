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

func (r *UserQuery) EqualToInt(key string, value int) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) NotEqualToInt(key string, value int) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) EqualToFloat(key string, value float64) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) NotEqualToFloat(key string, value float64) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) SizeEqualTo(key string, count int) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) GreaterThanInt(key string, value int) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) GreaterThanFloat(key string, value float64) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) GreaterThanEqualToInt(key string, value int) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) GreaterThanEqualToFloat(key string, value float64) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) LessThanInt(key string, value int) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) LessThanFloat(key string, value float64) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) LessThanEqualToInt(key string, value int) *UserQuery {
	// TODO
	return nil
}

func (r *UserQuery) LessThanEqualToFloat(key string, value float64) *UserQuery {
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

package lean

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

func (r *Query) EqualToInt(key string, value int) *Query {
	// TODO
	return nil
}

func (r *Query) NotEqualToInt(key string, value int) *Query {
	// TODO
	return nil
}

func (r *Query) EqualToFloat(key string, value float64) *Query {
	// TODO
	return nil
}

func (r *Query) NotEqualToFloat(key string, value float64) *Query {
	// TODO
	return nil
}

func (r *Query) SizeEqualTo(key string, count int) *Query {
	// TODO
	return nil
}

func (r *Query) GreaterThanInt(key string, value int) *Query {
	// TODO
	return nil
}

func (r *Query) GreaterThanFloat(key string, value float64) *Query {
	// TODO
	return nil
}

func (r *Query) GreaterThanEqualToInt(key string, value int) *Query {
	// TODO
	return nil
}

func (r *Query) GreaterThanEqualToFloat(key string, value float64) *Query {
	// TODO
	return nil
}

func (r *Query) LessThanInt(key string, value int) *Query {
	// TODO
	return nil
}

func (r *Query) LessThanFloat(key string, value float64) *Query {
	// TODO
	return nil
}

func (r *Query) LessThanEqualToInt(key string, value int) *Query {
	// TODO
	return nil
}

func (r *Query) LessThanEqualToFloat(key string, value float64) *Query {
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

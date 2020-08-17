package lean

type Query struct {
	c        *Client
	classRef *ClassRef
	where    map[string]string
	order    []string
	limit    int
	skip     int
}

func (r *Query) Find() (*[]ObjectRef, error) {
	// TODO
	return nil, nil
}

func (r *Query) First() (*ObjectRef, error) {
	// TODO
	return nil, nil
}

func (r *Query) Count() (int, error) {
	// TODO
	return 0, nil
}

func (r *Query) Or() *Query {
	// TODO
	return nil
}

func (r *Query) And() *Query {
	// TODO
	return nil
}

func (r *Query) Skip() *Query {
	// TODO
	return nil
}

func (r *Query) Limit() *Query {
	// TODO
	return nil
}

func (r *Query) Order() *Query {
	// TODO
	return nil
}

func (r *Query) EqualTo(key string, value string) *Query {
	// TODO
	return nil
}

func (r *Query) NotEqualTo(key string, value string) *Query {
	// TODO
	return nil
}

func (r *Query) SizeEqualTo() *Query {
	// TODO
	return nil
}

func (r *Query) Greater(key string, value string) *Query {
	// TODO
	return nil
}

func (r *Query) GreaterEqual() *Query {
	// TODO
	return nil
}

func (r *Query) Less(key string, value string) *Query {
	// TODO
	return nil
}

func (r *Query) LessEqual() *Query {
	// TODO
	return nil
}

func (r *Query) In(key string, value string) *Query {
	// TODO
	return nil
}

func (r *Query) NotIn(key string, value string) *Query {
	// TODO
	return nil
}

func (r *Query) All() *Query {
	// TODO
	return nil
}

func (r *Query) Exists() *Query {
	// TODO
	return nil
}

func (r *Query) Select() *Query {
	// TODO
	return nil
}

func (r *Query) Exclude() *Query {
	// TODO
	return nil
}

func (r *Query) Regexp() *Query {
	// TODO
	return nil
}

func (r *Query) Contains() *Query {
	// TODO
	return nil
}

func (r *Query) ContainsAll() *Query {
	// TODO
	return nil
}

func (r *Query) StartsWith() *Query {
	// TODO
	return nil
}

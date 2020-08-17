package lean

type ObjectQuery struct {
	c        *Client
	classRef *Class
	where    map[string]string
	order    []string
	limit    int
	skip     int
}

func (r *ObjectQuery) Find() ([]Object, error) {
	// TODO
	return nil, nil
}

func (r *ObjectQuery) First() (*Object, error) {
	// TODO
	return nil, nil
}

func (r *ObjectQuery) Count() (int, error) {
	// TODO
	return 0, nil
}

func (r *ObjectQuery) Skip() *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) Limit() *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) Order() *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) EqualTo(key string, value string) *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) NotEqualTo(key string, value string) *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) SizeEqualTo() *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) Greater(key string, value string) *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) GreaterEqual() *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) Less(key string, value string) *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) LessEqual() *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) In(key string, value string) *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) NotIn(key string, value string) *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) All() *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) Exists() *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) Select() *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) Exclude() *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) Regexp() *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) Contains() *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) ContainsAll() *ObjectQuery {
	// TODO
	return nil
}

func (r *ObjectQuery) StartsWith() *ObjectQuery {
	// TODO
	return nil
}

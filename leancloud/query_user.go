package leancloud

import "fmt"

type UserQuery struct {
	c     *Client
	class *Class
	where map[string]interface{}
	order []string
	limit int
	skip  int
}

func (q *UserQuery) Find(authOptions ...AuthOption) ([]User, error) {
	respUsers, err := objectQuery(q, false, false, authOptions...)
	if err != nil {
		return nil, err
	}

	users, ok := respUsers.([]User)
	if !ok {
		return nil, fmt.Errorf("unable to parse users from response")
	}

	return users, nil
}

func (q *UserQuery) First(authOptions ...AuthOption) (*User, error) {
	respUsers, err := objectQuery(q, false, true, authOptions...)
	if err != nil {
		return nil, err
	}

	users, ok := respUsers.([]User)
	if !ok {
		return nil, fmt.Errorf("unable to parse user from response")
	}

	if len(users) > 1 {
		return nil, fmt.Errorf("wrong count of response")
	}

	return &users[0], nil
}

func (q *UserQuery) Count(authOptions ...AuthOption) (int, error) {
	resp, err := objectQuery(q, true, false, authOptions...)
	if err != nil {
		return 0, err
	}

	count, ok := resp.(float64)
	if !ok {
		return 0, fmt.Errorf("unable to parse count from response")
	}

	return int(count), nil
}

func (q *UserQuery) Skip(count int) *UserQuery {
	q.skip = count
	return q
}

func (q *UserQuery) Limit(limit int) *UserQuery {
	q.limit = limit
	return q
}

func (q *UserQuery) Order(keys ...string) *UserQuery {
	q.order = keys
	return q
}

func (q *UserQuery) EqualTo(key string, value interface{}) *UserQuery {
	q.where[key] = wrapCondition("", value, "")
	return q
}

func (q *UserQuery) NotEqualTo(key string, value interface{}) *UserQuery {
	q.where[key] = wrapCondition("$ne", value, "")
	return q
}

func (q *UserQuery) SizeEqualTo(key string, count int) *UserQuery {
	q.where[key] = wrapCondition("$size", count, "")
	return nil
}

func (q *UserQuery) GreaterThan(key string, value interface{}) *UserQuery {
	q.where[key] = wrapCondition("$gt", value, "")
	return q
}

func (q *UserQuery) GreaterThanOrEqualTo(key string, value interface{}) *UserQuery {
	q.where[key] = wrapCondition("$lte", value, "")
	return q
}

func (q *UserQuery) LessThan(key string, value interface{}) *UserQuery {
	q.where[key] = wrapCondition("$lt", value, "")
	return q
}

func (q *UserQuery) LessThanOrEqualTo(key string, value interface{}) *UserQuery {
	q.where[key] = wrapCondition("$lte", value, "")
	return nil
}

func (q *UserQuery) In(key string, data interface{}) *UserQuery {
	q.where[key] = wrapCondition("$in", data, "")
	return q
}

func (q *UserQuery) NotIn(key string, data interface{}) *UserQuery {
	q.where[key] = wrapCondition("$nin", data, "")
	return q
}

func (q *UserQuery) Regexp(key, expr, options string) *UserQuery {
	q.where[key] = wrapCondition("$regex", expr, options)
	return q
}

func (q *UserQuery) Contains(key, substring string) *UserQuery {
	q.Regexp(key, substring, "")
	return q
}

func (q *UserQuery) ContainsAll(key string, objects interface{}) *UserQuery {
	q.where[key] = wrapCondition("$all", objects, "")
	return q
}

func (q *UserQuery) StartsWith(key, prefix string) *UserQuery {
	q.Regexp(key, fmt.Sprint("^", prefix), "")
	return q
}

package leancloud

import "fmt"

type UserQuery struct {
	c       *Client
	class   *Class
	where   map[string]interface{}
	include []string
	order   []string
	limit   int
	skip    int
}

func (q *UserQuery) Find(users interface{}, authOptions ...AuthOption) error {
	_, err := objectQuery(q, users, false, false, authOptions...)
	if err != nil {
		return err
	}

	return nil
}

func (q *UserQuery) First(user interface{}, authOptions ...AuthOption) error {
	_, err := objectQuery(q, user, false, true, authOptions...)
	if err != nil {
		return err
	}

	return nil
}

func (q *UserQuery) Count(authOptions ...AuthOption) (int, error) {
	resp, err := objectQuery(q, nil, true, false, authOptions...)
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
func (q *UserQuery) Or(queries ...*Query) *UserQuery {
	qArray := make([]map[string]interface{}, 1)
	for _, v := range queries {
		qArray = append(qArray, v.where)
	}
	q.where["$or"] = qArray
	return q
}

func (q *UserQuery) And(queries ...*UserQuery) *UserQuery {
	qArray := make([]map[string]interface{}, 1)
	for _, v := range queries {
		qArray = append(qArray, v.where)
	}
	q.where["$and"] = qArray
	return q
}

func (q *UserQuery) Near(key string, point *GeoPoint) *UserQuery {
	return q
}

func (q *UserQuery) WithinGeoBox(key string, point *GeoPoint) *UserQuery {
	return q
}

func (q *UserQuery) WithinKilometers(key string, point *GeoPoint) *UserQuery {
	return q
}

func (q *UserQuery) WithinMiles(key string, point *GeoPoint) *UserQuery {
	return q
}

func (q *UserQuery) WithinRadians(key string, point *GeoPoint) *UserQuery {
	return q
}

func (q *UserQuery) Include(keys ...string) *UserQuery {
	q.include = append(q.include, keys...)
	return q
}

func (q *UserQuery) Select(keys ...string) *UserQuery {
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

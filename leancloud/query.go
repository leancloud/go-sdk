package leancloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/levigross/grequests"
)

// Query contain parameters of queries
type Query struct {
	c          *Client
	class      *Class
	where      map[string]interface{}
	include    []string
	keys       []string
	order      []string
	limit      int
	skip       int
	includeACL bool
}

// Find fetch results of the Query
func (q *Query) Find(objects interface{}, authOptions ...AuthOption) error {
	_, err := objectQuery(q, objects, false, false, authOptions...)
	if err != nil {
		return err
	}

	return nil
}

// First fetch the first result of the Query
func (q *Query) First(object interface{}, authOptions ...AuthOption) error {
	_, err := objectQuery(q, object, false, true, authOptions...)
	if err != nil {
		return err
	}

	return nil
}

// Count returns the count of results of the Query
func (q *Query) Count(authOptions ...AuthOption) (int, error) {
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

func (q *Query) Skip(count int) *Query {
	q.skip = count
	return q
}

func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

func (q *Query) Order(keys ...string) *Query {
	q.order = keys
	return q
}

func (q *Query) Or(queries ...*Query) *Query {
	qArray := []map[string]interface{}{}
	for _, v := range queries {
		qArray = append(qArray, v.where)
	}
	q.where["$or"] = qArray
	return q
}

func (q *Query) And(queries ...*Query) *Query {
	qArray := []map[string]interface{}{}
	for _, v := range queries {
		qArray = append(qArray, v.where)
	}
	q.where["$and"] = qArray
	return q
}

func (q *Query) Near(key string, point *GeoPoint) *Query {
	q.where[key] = wrapCondition("$nearSphere", point, "")
	return q
}

func (q *Query) WithinGeoBox(key string, southwest *GeoPoint, northeast *GeoPoint) *Query {
	q.where[key] = wrapCondition("$withinBox", []GeoPoint{*southwest, *northeast}, "")
	return q
}

func (q *Query) WithinKilometers(key string, point *GeoPoint) *Query {
	q.where[key] = wrapCondition("$maxDistanceInKilometers", point, "")
	return q
}

func (q *Query) WithinMiles(key string, point *GeoPoint) *Query {
	q.where[key] = wrapCondition("$maxDistanceInMiles", point, "")
	return q
}

func (q *Query) WithinRadians(key string, point *GeoPoint) *Query {
	q.where[key] = wrapCondition("$maxDistanceInRadians", point, "")
	return q
}

func (q *Query) Include(keys ...string) *Query {
	q.include = append(q.include, keys...)
	return q
}

func (q *Query) Select(keys ...string) *Query {
	q.keys = append(q.keys, keys...)
	return q
}

func (q *Query) MatchesQuery(key string, query *Query) *Query {
	q.where[key] = wrapCondition("$inQuery", query, "")
	return q
}

func (q *Query) NotMatchesQuery(key string, query *Query) *Query {
	q.where[key] = wrapCondition("$notInQuery", query, "")
	return q
}

func (q *Query) MatchesKeyQuery(key, queryKey string, query *Query) *Query {
	q.where[key] = map[string]interface{}{
		"$select": map[string]interface{}{
			"query": map[string]interface{}{
				"className": query.class,
				"where":     query.where,
			},
			"key": queryKey,
		},
	}
	return q
}

func (q *Query) EqualTo(key string, value interface{}) *Query {
	q.where[key] = wrapCondition("", value, "")
	return q
}

func (q *Query) NotEqualTo(key string, value interface{}) *Query {
	q.where[key] = wrapCondition("$ne", value, "")
	return q
}
func (q *Query) Exists(key string) *Query {
	q.where[key] = wrapCondition("$exists", "", "")
	return q
}
func (q *Query) NotExists(key string) *Query {
	q.where[key] = wrapCondition("$notexists", "", "")
	return q
}

func (q *Query) GreaterThan(key string, value interface{}) *Query {
	q.where[key] = wrapCondition("$gt", value, "")
	return q
}

func (q *Query) GreaterThanOrEqualTo(key string, value interface{}) *Query {
	q.where[key] = wrapCondition("$gte", value, "")
	return q
}

func (q *Query) LessThan(key string, value interface{}) *Query {
	q.where[key] = wrapCondition("$lt", value, "")
	return q
}

func (q *Query) LessThanOrEqualTo(key string, value interface{}) *Query {
	q.where[key] = wrapCondition("$lte", value, "")
	return q
}

func (q *Query) In(key string, data interface{}) *Query {
	q.where[key] = wrapCondition("$in", data, "")
	return q
}

func (q *Query) NotIn(key string, data interface{}) *Query {
	q.where[key] = wrapCondition("$nin", data, "")
	return q
}

func (q *Query) Regexp(key, expr, options string) *Query {
	q.where[key] = wrapCondition("$regex", expr, options)
	return q
}

func (q *Query) Contains(key, substring string) *Query {
	q.Regexp(key, substring, "")
	return q
}

func (q *Query) ContainsAll(key string, objects interface{}) *Query {
	q.where[key] = wrapCondition("$all", objects, "")
	return q
}

func (q *Query) StartsWith(key, prefix string) *Query {
	q.Regexp(key, fmt.Sprint("^", prefix), "")
	return q
}

func (q *Query) IncludeACL() *Query {
	q.includeACL = true
	return q
}

func wrapCondition(verb string, value interface{}, options string) interface{} {
	switch verb {
	case "$ne", "$lt", "$lte", "$gt", "$gte", "$in", "$nin", "$all", "nearShpere":
		return map[string]interface{}{
			verb: encode(value, false),
		}
	case "$withinBox":
		return encode(map[string]interface{}{
			"$box": value,
		}, true)
	case "$regex":
		return map[string]interface{}{
			"$regex":   value,
			"$options": options,
		}
	case "$exists":
		return map[string]interface{}{
			"$exists": true,
		}
	case "$notexists":
		return map[string]interface{}{
			"$exists": false,
		}
	case "$inQuery":
		queryMap, err := formatQuery(value, false, false)
		if err != nil {
			return nil
		}
		queryMap["className"] = value.(*Query).class.Name
		return map[string]interface{}{
			"$inQuery": queryMap,
		}
	case "$notInQuery":
		queryMap, err := formatQuery(value, false, false)
		if err != nil {
			return nil
		}
		queryMap["className"] = value.(*Query).class.Name
		return map[string]interface{}{
			"$notInQuery": queryMap,
		}
	default:
		switch v := value.(type) {
		case time.Time:
			return encodeDate(&v)
		default:
			return encode(value, false)
		}
	}
}

func objectQuery(query interface{}, objects interface{}, count bool, first bool, authOptions ...AuthOption) (interface{}, error) {
	path := fmt.Sprint("/1.1/")
	var client *Client
	var options *grequests.RequestOptions
	params, err := wrapParams(query, count, first)
	if err != nil {
		return nil, err
	}

	switch v := query.(type) {
	case *Query:
		if v.class.Name == "_User" {
			path = fmt.Sprint(path, "users")
		} else if v.class.Name == "_File" {
			path = fmt.Sprint(path, "classes/files")
		} else if v.class.Name == "_Role" {
			path = fmt.Sprint(path, "roles")
		} else {
			path = fmt.Sprint(path, "classes/", v.class.Name)
		}
		options = v.c.getRequestOptions()
		client = v.c
	}

	options.Params = params

	resp, err := client.request(methodGet, path, options, authOptions...)
	if err != nil {
		return nil, err
	}

	respJSON := make(map[string]interface{})
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return nil, fmt.Errorf("unable to parse response %w", err)
	}
	if count {
		return respJSON["count"], nil
	}

	results := respJSON["results"].([]interface{})
	switch query.(type) {
	case *Query:
		decodedObjects, err := decodeArray(results, true)

		if err != nil {
			return nil, err
		}

		rDecodedObjects := reflect.ValueOf(decodedObjects)

		if !first {
			if err := bind(rDecodedObjects, reflect.ValueOf(objects).Elem()); err != nil {
				return nil, err
			}
		} else if rDecodedObjects.Len() == 0 {
			return nil, errors.New("no matched object found")
		} else {
			if err := bind(rDecodedObjects.Index(0), reflect.ValueOf(objects).Elem()); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func wrapParams(query interface{}, count, first bool) (map[string]string, error) {
	var where map[string]interface{}
	var order string
	var include string
	var keys string
	var skip, limit int

	switch v := query.(type) {
	case *Query:
		where = v.where
		order = strings.Join(v.order, ",")
		include = strings.Join(v.include, ",")
		keys = strings.Join(v.keys, ",")
		skip, limit = v.skip, v.limit
	}

	whereString, err := json.Marshal(where)
	if err != nil {
		return nil, fmt.Errorf("unable to wrap params %w", err)
	}

	params := map[string]string{
		"where": string(whereString),
	}

	if skip != 0 {
		params["skip"] = fmt.Sprintf("%d", skip)
	}

	if limit != 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}

	if len(order) != 0 {
		params["order"] = order
	}

	if len(include) != 0 {
		params["include"] = include
	}

	if len(keys) != 0 {
		params["keys"] = keys
	}

	if count {
		params["count"] = "1"
	}

	if first {
		params["limit"] = "1"
	}

	return params, nil
}

func formatQuery(query interface{}, count, first bool) (map[string]interface{}, error) {
	paramsInString, err := wrapParams(query, count, first)
	if err != nil {
		return nil, err
	}
	paramsInInterface := make(map[string]interface{})
	for k, v := range paramsInString {
		paramsInInterface[k] = interface{}(v)
	}
	paramsInInterface["where"] = query.(*Query).where

	return paramsInInterface, nil
}

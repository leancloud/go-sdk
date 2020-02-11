package lean

import (
	"encoding/json"
	"fmt"

	"github.com/levigross/grequests"
)

type Query struct {
	client    *Client
	className string
	where     map[string]string
	order     []string
	limit     int
	skip      int
}

func (client *Client) NewQuery(className string) *Query {
	return &Query{
		client:    client,
		className: className,
	}
}

func (client *Client) CloudQuery(cql string, pvalues interface{}, objects interface{}, authOptions ...AuthOption) error {
	options := client.getRequestOptions()

	params := map[string]string{
		"cql": cql,
	}

	if pvalues != nil {
		jsonBytes, err := json.Marshal(pvalues)

		if err != nil {
			return err
		}

		params["pvalues"] = string(jsonBytes)
	}

	options.Params = params

	resp, err := client.request(ServiceAPI, methodGet, "/1.1/cloudQuery", options, authOptions...)

	if err != nil {
		return err
	}

	result := &cqlResponse{}

	err = resp.JSON(result)

	if err != nil {
		return err
	}

	return decodeObjects(result.ClassName, result.Results, objects)
}

func (query *Query) EqualTo(key string, value string) *Query {
	query.where[key] = value
	return query
}

func (query *Query) Limit(limit int) *Query {
	return query
}

func (query *Query) Skip(limit int) *Query {
	return query
}

func (query *Query) ascend(field string) *Query {
	return query
}

func (query *Query) descend(field string) *Query {
	return query
}

func (query *Query) Get(object Object, authOptions ...AuthOption) error {
	resp, err := query.request(fmt.Sprint("/1.1/classes/", query.className, "/", object.getObjectMeta().ObjectID), authOptions)

	if err != nil {
		return err
	}

	result := make(objectResponse)

	err = resp.JSON(&result)

	if err != nil {
		return err
	}

	return decodeObject(result, object)
}

func (query *Query) Count() (rows int, err error) {
	return 0, nil
}

func (query *Query) First(object Object, authOptions ...AuthOption) error {
	query.Limit(1)

	resp, err := query.request(fmt.Sprint("/1.1/classes/", query.className), authOptions)

	if err != nil {
		return err
	}

	result := &objectsResponse{}

	err = resp.JSON(result)

	if err != nil {
		return err
	}

	return decodeObject(result.Results[0], object)
}

func (query *Query) Find(objects interface{}, authOptions ...AuthOption) error {
	resp, err := query.request(fmt.Sprint("/1.1/classes/", query.className), authOptions)

	if err != nil {
		return err
	}

	result := &objectsResponse{}

	err = resp.JSON(result)

	if err != nil {
		return err
	}

	return decodeObjects(query.className, result.Results, objects)
}

func (query *Query) request(path string, authOptions []AuthOption) (*grequests.Response, error) {
	params, err := encodeQuery(query)

	if err != nil {
		return nil, err
	}

	options := query.client.getRequestOptions()

	options.Params = params

	return query.client.request(ServiceAPI, methodGet, path, options, authOptions...)
}

func encodeQuery(query *Query) (params map[string]string, err error) {
	if len(query.where) > 0 {
		jsonBytes, err := json.Marshal(query.where)

		if err != nil {
			return params, err
		}

		params["where"] = string(jsonBytes)
	}

	if query.limit != 0 {
		params["limit"] = string(query.limit)
	}

	if query.skip != 0 {
		params["skip"] = string(query.skip)
	}

	return params, nil
}

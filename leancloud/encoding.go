package leancloud

import (
	"fmt"
	"reflect"
	"time"
)

func encodeObject(object interface{}) map[string]interface{} {
	mapObject := make(map[string]interface{})
	if reflect.TypeOf(object).Kind() == reflect.Struct {
		v := reflect.ValueOf(object)
		s := reflect.TypeOf(object)
		for i := 0; i < v.NumField(); i++ {
			tag, ok := s.Field(i).Tag.Lookup("json")
			if !ok || tag == "" {
				tag = s.Field(i).Name
			}
			switch v.Field(i).Type() {
			case reflect.TypeOf(time.Time{}):
				date := v.Field(i).Interface().(time.Time)
				mapObject[tag] = encodeDate(date)
				break
			default:
				mapObject[s.Field(i).Tag.Get("json")] = v.Field(i).Interface()
			}
		}
	}
	if reflect.TypeOf(object).Kind() == reflect.Map {
		iter := reflect.ValueOf(object).MapRange()
		for iter.Next() {
			switch iter.Value().Elem().Type() {
			case reflect.TypeOf(time.Time{}):
				date := iter.Value().Interface().(time.Time)
				mapObject[iter.Key().String()] = encodeDate(date)
				break
			default:
				mapObject[iter.Key().String()] = iter.Value().Interface()
			}
		}
	}

	return mapObject
}

func encodeDate(date time.Time) map[string]interface{} {
	return map[string]interface{}{
		"__type": "Date",
		"iso":    fmt.Sprint(date.In(time.FixedZone("UTC", 0)).Format("2006-01-02T15:04:05.000Z")),
	}
}

func decodeObject(fields map[string]interface{}, object interface{}) {
	v := reflect.ValueOf(object)
	s := reflect.TypeOf(object)
	sv := reflect.ValueOf(&object).Elem()
	for i := 0; i < v.NumField(); i++ {
		tag, ok := s.Field(i).Tag.Lookup("json")
		if !ok || tag == "" {
			tag = s.Field(i).Name
		}
		if fields[tag] != nil {
			switch reflect.ValueOf(fields[tag]).Kind() {
			case reflect.Interface:
				break
			default:
				sv.Field(i).Set(reflect.ValueOf(fields[tag]))
			}
		}
	}
}

func decodeDate(data map[string]interface{}) (*time.Time, error) {
	date, err := time.Parse(time.RFC3339, data["iso"].(string))
	if err != nil {
		return nil, err
	}
	return &date, nil
}

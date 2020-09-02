package leancloud

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func encodeObject(object interface{}) map[string]interface{} {
	mapObject := make(map[string]interface{})
	if reflect.TypeOf(object).Kind() == reflect.Struct {
		v := reflect.ValueOf(object)
		t := reflect.TypeOf(object)
		for i := 0; i < v.NumField(); i++ {
			tag, option := parseTag(t.Field(i).Tag.Get("json"))
			if option == "omitempty" && v.Field(i).IsZero() {
				continue
			}
			if tag == "" {
				tag = t.Field(i).Name
			}
			switch v.Field(i).Type() {
			case reflect.TypeOf(time.Time{}):
				date := v.Field(i).Interface().(time.Time)
				mapObject[tag] = encodeDate(date)
				break
			default:
				mapObject[tag] = v.Field(i).Interface()
			}
		}
	} else if reflect.TypeOf(object).Kind() == reflect.Map {
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

func decodeFields(fields map[string]interface{}) map[string]interface{} {
	objectMap := make(map[string]interface{})
	iter := reflect.ValueOf(fields).MapRange()
	for iter.Next() {
		switch iter.Value().Elem().Kind() {
		case reflect.Map:
			intf, _ := iter.Value().Interface().(map[string]interface{})
			if reflect.ValueOf(intf["__type"]).IsValid() {
				switch intf["__type"].(string) {
				case "Date":
					date, _ := decodeDate(intf)
					objectMap[iter.Key().String()] = date
					break
				}
			}
			break
		default:
			objectMap[iter.Key().String()] = iter.Value().Interface()
		}
	}

	return objectMap
}

func decodeObject(fields map[string]interface{}, object interface{}) {
	v := reflect.Indirect(reflect.ValueOf(object))
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		tag, ok := t.Field(i).Tag.Lookup("json")
		if !ok || tag == "" {
			tag = t.Field(i).Name
		}

		fv := reflect.ValueOf(fields[tag])
		if fields[tag] != nil {
			switch fv.Kind() {
			case reflect.Map:
				data, _ := fields[tag].(map[string]interface{})
				mapType, _ := data["__type"].(string)
				switch mapType {
				case "Date":
					date, _ := decodeDate(data)
					v.Field(i).Set(reflect.ValueOf(date))
				}
			case reflect.String:
				if tag == "createdAt" || tag == "updatedAt" {
					timeAt, err := time.Parse(time.RFC3339, fv.String())
					if err != nil {
						panic(err)
					}
					v.Field(i).Set(reflect.ValueOf(timeAt))
				} else {
					v.Field(i).Set(fv)
				}
			default:
				v.Field(i).Set(fv.Convert(t.Field(i).Type))
			}
		}
	}

	switch v := object.(type) {
	case *Object:
		v.fields = fields
	case *User:
		v.fields = fields
	}
}

func decodeDate(data map[string]interface{}) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, data["iso"].(string))
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

func parseTag(tag string) (name string, option string) {
	parts := strings.Split(tag, ",")

	if len(parts) > 1 {
		return parts[0], parts[1]
	}

	return parts[0], ""
}

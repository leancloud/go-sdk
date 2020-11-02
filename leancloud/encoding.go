package leancloud

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func encode(object interface{}) interface{} {
	switch o := object.(type) {
	case Object:
		if o.isPointer {
			return map[string]interface{}{
				"__type":    "Pointer",
				"className": o.fields["className"],
				"objectId":  o.fields["objectId"],
				"createdAt": encode(o.CreatedAt),
				"updatedAt": encode(o.UpdatedAt),
			}
		}
		r := encode(o.fields)
		rmap, _ := r.(map[string]interface{})
		rmap["__type"] = "Object"
		return rmap
	case User:
		r := encode(o.fields)
		rmap, _ := r.(map[string]interface{})
		rmap["__type"] = "Object"
		return rmap
	case GeoPoint:
		return map[string]interface{}{
			"__type":    "GeoPoint",
			"latitude":  o.Latitude,
			"longitude": o.Longitude,
		}
	case time.Time:
		return map[string]interface{}{
			"__type": "Date",
			"iso":    fmt.Sprint(o.In(time.FixedZone("UTC", 0)).Format("2006-01-02T15:04:05.000Z")),
		}
	case []byte:
		return map[string]interface{}{
			"__type": "Byte",
			"base64": base64.StdEncoding.EncodeToString(o),
		}
	case File:
		return map[string]interface{}{
			"__type": "File",
		}
	case Relation:
		return map[string]interface{}{
			"__type": "Relation",
		}
	case ACL:
		return map[string]interface{}{
			"__type": "ACL",
		}
	default:
		v := reflect.ValueOf(object)
		var array []interface{}
		switch v.Kind() {
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			for i := 0; i < v.Len(); i++ {
				r := encode(v.Index(i).Interface())
				array = append(array, r)
			}
			return array
		case reflect.Map:
			mapObject := make(map[string]interface{})
			iter := v.MapRange()
			for iter.Next() {
				mapObject[iter.Key().String()] = encode(iter.Value().Interface())
			}
			return mapObject
		default:
			return object
		}
	}
}

func encodeObject(object interface{}) map[string]interface{} {
	v := reflect.ValueOf(object)
	t := reflect.TypeOf(object)
	mapObject := make(map[string]interface{})
	mapObject["__type"] = "Object"
	for i := 0; i < v.NumField(); i++ {
		tag, option := parseTag(t.Field(i).Tag.Get("json"))
		if option == "omitempty" && v.Field(i).IsZero() {
			continue
		}
		if tag == "" {
			tag = t.Field(i).Name
		}
		mapObject[tag] = encode(v.Field(i).Interface())
	}
	return mapObject
}

func encodeDate(date time.Time) map[string]interface{} {
	return map[string]interface{}{
		"__type": "Date",
		"iso":    fmt.Sprint(date.In(time.FixedZone("UTC", 0)).Format("2006-01-02T15:04:05.000Z")),
	}
}

func transform(fields map[string]interface{}, object interface{}) error {
	return nil
}

func decode(fields interface{}) (interface{}, error) {
	mapFields, ok := fields.(map[string]interface{})
	if !ok {
		switch reflect.ValueOf(fields).Kind() {
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			var array []interface{}
			v := reflect.ValueOf(fields)
			for i := 0; i < v.Len(); i++ {
				r, err := decode(v.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				array = append(array, r)
			}
			return array, nil
		default:
			return fields, nil
		}
	}
	if mapFields["__type"] != nil {
		fieldType, ok := mapFields["__type"].(string)
		if !ok {
			return nil, fmt.Errorf("unexpected error when parse __type: want string but %v", reflect.TypeOf(mapFields["__type"]))
		}

		switch fieldType {
		case "Pointer":
			fallthrough
		case "Object":
			mapObject := make(map[string]interface{})
			iter := reflect.ValueOf(mapFields).MapRange()
			for iter.Next() {
				if iter.Key().String() != "__type" {
					r, err := decode(iter.Value().Interface())
					if err != nil {
						return nil, err
					}
					mapObject[iter.Key().String()] = r
				}
			}
			createdAt, ok := mapObject["createdAt"].(string)
			if !ok {
				return nil, fmt.Errorf("unexpected error when parse createdAt: want type string but %v", reflect.TypeOf(mapObject["createdAt"]))
			}
			realCreatedAt, err := time.Parse(time.RFC3339, createdAt)
			if err != nil {
				return nil, fmt.Errorf("unexpected error when parse createdAt: %v", err)
			}

			updatedAt, ok := mapObject["updatedAt"].(string)
			if !ok {
				return nil, fmt.Errorf("unexpected error when parse updatedAt: want type string but %v", reflect.TypeOf(mapObject["createdAt"]))
			}
			realUpdatedAt, err := time.Parse(time.RFC3339, updatedAt)
			if err != nil {
				return nil, fmt.Errorf("unexpected error when parse updatedAt: %v", err)
			}

			objectID, ok := mapObject["objectId"].(string)
			if !ok {
				return nil, fmt.Errorf("unexpected error when parse objectId: want type string but %v", reflect.TypeOf(mapObject["objectId"]))
			}

			object := Object{
				ID:        objectID,
				CreatedAt: realCreatedAt,
				UpdatedAt: realUpdatedAt,
				fields:    mapObject,
				isPointer: false,
			}

			if mapFields["__type"] == "Pointer" {
				object.isPointer = true
			}

			return object, nil
		case "Date":
			iso, ok := mapFields["iso"].(string)
			if !ok {
				return nil, fmt.Errorf("unexpected error when parse Date: iso expected string but %v", reflect.TypeOf(mapFields["iso"]))
			}
			date, err := time.Parse(time.RFC3339, iso)
			if err != nil {
				return nil, fmt.Errorf("unexpected error when parse Date %v", err)
			}
			return date, nil
		case "Byte":
			base64String, ok := mapFields["base64"].(string)
			if !ok {
				return nil, fmt.Errorf("unexpected error when parse Byte: base64 string expected string but %v", reflect.TypeOf(mapFields["base64"]))
			}
			bytes, err := base64.StdEncoding.DecodeString(base64String)
			if err != nil {
				return nil, fmt.Errorf("unexpected error when parse Byte %v", err)
			}
			return bytes, nil
		case "File":
			return nil, nil
		case "Relation":
			return nil, nil
		case "ACL":
			return nil, nil
		default:
			return fields, nil
		}
	} else {
		mapObject := make(map[string]interface{})
		iter := reflect.ValueOf(mapFields).MapRange()
		for iter.Next() {
			if iter.Key().String() != "__type" {
				r, err := decode(iter.Value().Interface())
				if err != nil {
					return nil, err
				}
				mapObject[iter.Key().String()] = r
			}
		}

		return mapObject, nil
	}
}

func decodeFields(fields map[string]interface{}) (map[string]interface{}, error) {
	objectMap := make(map[string]interface{})
	iter := reflect.ValueOf(fields).MapRange()
	for iter.Next() {
		switch iter.Value().Elem().Kind() {
		case reflect.Map:
			intf, ok := iter.Value().Interface().(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("unable to assert type map from fields")
			}
			if reflect.ValueOf(intf["__type"]).IsValid() {
				switch intf["__type"].(string) {
				case "Date":
					date, err := decodeDate(intf)
					if err != nil {
						return nil, fmt.Errorf("unable to decode Date %w", err)
					}
					objectMap[iter.Key().String()] = date
					break
				}
			}
			break
		default:
			objectMap[iter.Key().String()] = iter.Value().Interface()
		}
	}

	return objectMap, nil
}

func decodeObject(fields map[string]interface{}, object interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(object))
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		tag, ok := t.Field(i).Tag.Lookup("json")
		if !ok || tag == "" {
			tag = t.Field(i).Name
		}

		if fields[tag] != nil {
			fv := reflect.ValueOf(fields[tag])
			switch fv.Kind() {
			case reflect.Map:
				data, ok := fields[tag].(map[string]interface{})
				if !ok {
					return fmt.Errorf("unable to assert type map from fields")
				}
				mapType, ok := data["__type"].(string)
				if !ok {
					return fmt.Errorf("unable to assert type string from fields")
				}
				switch mapType {
				case "Date":
					date, err := decodeDate(data)
					if err != nil {
						object = nil
						return fmt.Errorf("unable to decode Date %w", err)
					}
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

	return nil
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

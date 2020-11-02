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
		return encodeGeoPoint(&o)
	case time.Time:
		return encodeDate(o)
	case []byte:
		return encodeBytes(o)
	case File:
		return encodeFile(&o)
	case Relation:
		return encodeRelation(&o)
	case ACL:
		return encodeACL(&o)
	default:
		switch reflect.ValueOf(object).Kind() {
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			return encodeArray(object)
		case reflect.Map:
			return encodeMap(o)
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

func encodeMap(fields interface{}) map[string]interface{} {
	encodedMap := make(map[string]interface{})
	v := reflect.ValueOf(fields)

	for iter := v.MapRange(); iter.Next(); {
		encodedMap[iter.Key().String()] = encode(iter.Value().Interface())
	}

	return encodedMap
}

func encodeArray(array interface{}) []interface{} {
	var encodedArray []interface{}
	v := reflect.ValueOf(array)
	for i := 0; i < v.Len(); i++ {
		encodedArray = append(encodedArray, encode(v.Index(i).Interface()))
	}

	return encodedArray
}

func encodePointer(pointer *Object) map[string]interface{} {
	return map[string]interface{}{
		"__type":    "Pointer",
		"objectId":  pointer.ID,
		"className": pointer.fields["className"],
	}
}

func encodeDate(date time.Time) map[string]interface{} {
	return map[string]interface{}{
		"__type": "Date",
		"iso":    fmt.Sprint(date.In(time.FixedZone("UTC", 0)).Format("2006-01-02T15:04:05.000Z")),
	}
}

func encodeGeoPoint(point *GeoPoint) map[string]interface{} {
	return map[string]interface{}{
		"__type":    "GeoPoint",
		"latitude":  point.Latitude,
		"longitude": point.Longitude,
	}
}

func encodeBytes(bytes []byte) map[string]interface{} {
	return map[string]interface{}{
		"__type": "Byte",
		"base64": base64.StdEncoding.EncodeToString(bytes),
	}
}

func encodeFile(file *File) map[string]interface{} {
	return nil
}

func encodeACL(acl *ACL) map[string]interface{} {
	return nil
}

func encodeRelation(relation *Relation) map[string]interface{} {
	return nil
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
			return decodeArray(fields)
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
			return decodeObject(fields)
		case "Date":
			iso, ok := mapFields["iso"].(string)
			if !ok {
				return nil, fmt.Errorf("unexpected error when parse Date: iso expected string but %v", reflect.TypeOf(mapFields["iso"]))
			}
			return decodeDate(iso)
		case "Byte":
			base64String, ok := mapFields["base64"].(string)
			if !ok {
				return nil, fmt.Errorf("unexpected error when parse Byte: base64 string expected string but %v", reflect.TypeOf(mapFields["base64"]))
			}
			return decodeBytes(base64String)
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
		return decodeMap(fields)
	}
}

func decodeObject(fields interface{}) (*Object, error) {
	decodedFields, err := decodeMap(fields)
	if err != nil {
		return nil, err
	}

	objectID, ok := decodedFields["objectId"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected error when parse objectId: want type string but %v", reflect.TypeOf(decodedFields["objectId"]))
	}

	createdAt, ok := decodedFields["createdAt"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected error when parse createdAt: want type string but %v", reflect.TypeOf(decodedFields["createdAt"]))
	}
	decodedCreatedAt, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, fmt.Errorf("unexpected error when parse createdAt: %v", err)
	}

	updatedAt, ok := decodedFields["updatedAt"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected error when parse updatedAt: want type string but %v", reflect.TypeOf(decodedFields["createdAt"]))
	}
	decodedUpdatedAt, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("unexpected error when parse updatedAt: %v", err)
	}

	return &Object{
		ID:        objectID,
		CreatedAt: decodedCreatedAt,
		UpdatedAt: decodedUpdatedAt,
		fields:    decodedFields,
		isPointer: false,
	}, nil
}

func decodePointer(pointer interface{}) (*Object, error) {
	decodedFields, err := decodeMap(pointer)
	if err != nil {
		return nil, err
	}

	objectID, ok := decodedFields["objectId"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected error when parse objectId: want type string but %v", reflect.TypeOf(decodedFields["objectId"]))
	}
	return &Object{
		ID:        objectID,
		isPointer: true,
		fields:    decodedFields,
	}, nil
}

func decodeArray(array interface{}) ([]interface{}, error) {
	var decodedArray []interface{}
	v := reflect.ValueOf(array)
	for i := 0; i < v.Len(); i++ {
		r, err := decode(v.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		decodedArray = append(decodedArray, r)
	}
	return decodedArray, nil
}

func decodeMap(fields interface{}) (map[string]interface{}, error) {
	var decodedMap map[string]interface{}
	iter := reflect.ValueOf(fields).MapRange()
	for iter.Next() {
		if iter.Key().String() != "__type" {
			r, err := decode(iter.Value().Interface())
			if err != nil {
				return nil, err
			}
			decodedMap[iter.Key().String()] = r
		}
	}
	return decodedMap, nil
}

func decodeBytes(bytesStr string) ([]byte, error) {
	bytes, err := base64.StdEncoding.DecodeString(bytesStr)
	if err != nil {
		return nil, fmt.Errorf("unexpected error when parse Byte %v", err)
	}
	return bytes, nil
}

func decodeDate(dateStr string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

/*
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
*/

func parseTag(tag string) (name string, option string) {
	parts := strings.Split(tag, ",")

	if len(parts) > 1 {
		return parts[0], parts[1]
	}

	return parts[0], ""
}

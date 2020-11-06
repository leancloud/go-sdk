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
			return encodePointer(&o)
		}
		encodedObject := encodeMap(o.fields)
		encodedObject["__type"] = "Object"
		return encodedObject
	case User:
		return encodeMap(o.fields)
	case GeoPoint:
		return encodeGeoPoint(&o)
	case time.Time:
		return encodeDate(o)
	case []byte:
		return encodeBytes(o)
	case File:
		return encodeFile(&o, true)
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
	encodedObject := make(map[string]interface{})
	encodedObject["__type"] = "Object"
	for i := 0; i < v.NumField(); i++ {
		tag, option := parseTag(t.Field(i).Tag.Get("json"))
		if option == "omitempty" && v.Field(i).IsZero() {
			continue
		}
		if tag == "" {
			tag = t.Field(i).Name
		}
		encodedObject[tag] = encode(v.Field(i).Interface())
	}
	return encodedObject
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

func encodeFile(file *File, embedded bool) map[string]interface{} {
	if embedded {
		return map[string]interface{}{
			"__type": "File",
			"id":     file.ID,
		}
	}

	return map[string]interface{}{
		"__type":    "File",
		"name":      file.Name,
		"mime_type": file.MIME,
		"metaData":  file.Meatadata,
	}
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
		case "GeoPoint":
			return decodeGeoPoint(mapFields)
		case "File":
			return decodeFile(mapFields)
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
		if decodedFields["updatedAt"] == nil {
			updatedAt = ""
		} else {
			return nil, fmt.Errorf("unexpected error when parse updatedAt: want type string but %v", reflect.TypeOf(decodedFields["updatedAt"]))
		}
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

func decodeUser(fields interface{}) (*User, error) {
	object, err := decodeObject(fields)
	if err != nil {
		return nil, err
	}

	sessionToken, ok := object.fields["sessionToken"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected error when parse sessionToken: want type string but %v", reflect.TypeOf(object.fields["sessionToken"]))
	}
	return &User{
		Object:       *object,
		sessionToken: sessionToken,
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
	decodedMap := make(map[string]interface{})
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

func decodeGeoPoint(v map[string]interface{}) (*GeoPoint, error) {
	latitude, ok := v["latitude"].(float64)
	if !ok {
		return nil, fmt.Errorf("latitude want type float64 but %v", reflect.TypeOf(v["latitude"]))
	}
	longitude, ok := v["longitude"].(float64)
	if !ok {
		return nil, fmt.Errorf("longitude want type float64 but %v", reflect.TypeOf(v["longitude"]))
	}
	return &GeoPoint{
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}

func decodeFile(fields map[string]interface{}) (*File, error) {
	file := new(File)

	decodedFields, err := decodeMap(fields)
	if err != nil {
		return nil, err
	}
	file.fields = decodedFields

	objectID, ok := decodedFields["objectId"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected error when parse objectId: want type string but %v", reflect.TypeOf(decodedFields["objectId"]))
	}
	file.ID = objectID

	createdAt, ok := decodedFields["createdAt"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected error when parse createdAt: want type string but %v", reflect.TypeOf(decodedFields["createdAt"]))
	}
	decodedCreatedAt, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, fmt.Errorf("unexpected error when parse createdAt: %v", err)
	}
	file.CreatedAt = decodedCreatedAt

	updatedAt, ok := decodedFields["updatedAt"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected error when parse updatedAt: want type string but %v", reflect.TypeOf(decodedFields["updatedAt"]))
	}
	decodedUpdatedAt, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("unexpected error when parse updatedAt: %v", err)
	}
	file.UpdatedAt = decodedUpdatedAt

	key, ok := decodedFields["key"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected error when parse key from response: want type string but %v", reflect.TypeOf(decodedFields["key"]))
	}
	file.Key = key

	url, ok := decodedFields["url"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected error when parse url from response: want type string but %v", reflect.TypeOf(decodedFields["url"]))
	}
	file.URL = url

	bucket, ok := decodedFields["bucket"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected error when parse bucket from response: want type string but %v", reflect.TypeOf(decodedFields["bucket"]))
	}
	file.Bucket = bucket

	provider, ok := decodedFields["provider"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected error when parse provider from response: want type string but %v", reflect.TypeOf(decodedFields["provider"]))
	}
	file.Provider = provider

	return file, nil
}

func parseTag(tag string) (name string, option string) {
	parts := strings.Split(tag, ",")

	if len(parts) > 1 {
		return parts[0], parts[1]
	}

	return parts[0], ""
}

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
		return encodePointer(&o)
	case User:
		return encodeMap(&o)
	case GeoPoint:
		return encodeGeoPoint(&o)
	case time.Time:
		return encodeDate(&o)
	case File:
		return encodeFile(&o, true)
	case Relation:
		return encodeRelation(&o)
	case ACL:
		return encodeACL(&o)
	case *Object:
		return encodePointer(o)
	case *User:
		return encodeMap(o.fields)
	case *GeoPoint:
		return encodeGeoPoint(o)
	case *time.Time:
		return encodeDate(o)
	case []byte:
		return encodeBytes(o)
	case *File:
		return encodeFile(o, true)
	case *Relation:
		return encodeRelation(o)
	case *ACL:
		return encodeACL(o)
	default:
		switch reflect.ValueOf(object).Kind() {
		case reflect.Slice, reflect.Array:
			return encodeArray(object)
		case reflect.Map:
			return encodeMap(o)
		case reflect.Struct:
			return encodeObject(o)
		case reflect.Interface, reflect.Ptr:
			return encode(reflect.Indirect(reflect.ValueOf(o)).Interface())
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
		if tag == "createdAt" || tag == "updatedAt" {
			date, _ := v.Field(i).Interface().(time.Time)
			encodedObject[tag] = fmt.Sprint(date.In(time.FixedZone("UTC", 0)).Format("2006-01-02T15:04:05.000Z"))
		} else {
			if v.Field(i).Kind() == reflect.Ptr && v.Field(i).IsNil() {
				continue
			}
			encodedObject[tag] = encode(v.Field(i).Interface())
		}
	}
	return encodedObject
}

func encodeMap(fields interface{}) map[string]interface{} {
	encodedMap := make(map[string]interface{})
	v := reflect.ValueOf(fields)
	for iter := v.MapRange(); iter.Next(); {
		if iter.Key().String() == "createdAt" || iter.Key().String() == "updatedAt" {
			date, _ := iter.Value().Interface().(time.Time)
			encodedMap[iter.Key().String()] = fmt.Sprint(date.In(time.FixedZone("UTC", 0)).Format("2006-01-02T15:04:05.000Z"))
		} else {
			encodedMap[iter.Key().String()] = encode(iter.Value().Interface())
		}
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

func encodePointer(pointer interface{}) map[string]interface{} {
	switch v := pointer.(type) {
	case *Object:
		if v.ref == nil {
			return nil
		}
		return map[string]interface{}{
			"__type":    "Pointer",
			"objectId":  v.ID,
			"className": v.ref.class,
		}
	case *User:
		if v.ref == nil {
			return nil
		}
		return map[string]interface{}{
			"__type":    "Pointer",
			"objectId":  v.ID,
			"className": v.ref.class,
		}
	}

	return nil
}

func encodeDate(date *time.Time) map[string]interface{} {
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
	if len(bytes) == 0 {
		return nil
	}

	return map[string]interface{}{
		"__type": "Bytes",
		"base64": base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(string(bytes)))),
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

func encodeACL(acl *ACL) map[string]map[string]bool {
	return acl.content
}

func encodeRelation(relation *Relation) map[string]interface{} {
	return nil
}

func bind(src reflect.Value, dst reflect.Value) error {
	tdst := dst.Type()
	switch dst.Kind() {
	case reflect.Struct:
		if src.Kind() == reflect.Map {
			for i := 0; i < tdst.NumField(); i++ {
				tag, ok := tdst.Field(i).Tag.Lookup("json")
				if !ok || tag == "" {
					tag = tdst.Field(i).Name
				}
				mapIndex := src.MapIndex(reflect.ValueOf(tag))
				if mapIndex.Kind() == reflect.Ptr && mapIndex.IsNil() {
					continue
				}
				if mapIndex.IsValid() {
					if dst.Field(i).Kind() == reflect.Ptr && dst.Field(i).IsNil() {
						pv := reflect.New(dst.Field(i).Type().Elem())
						if err := bind(src.MapIndex(reflect.ValueOf(tag)), pv); err != nil {
							return err
						}
						dst.Field(i).Set(pv)
					} else {
						if err := bind(src.MapIndex(reflect.ValueOf(tag)), dst.Field(i)); err != nil {
							return err
						}
					}
				}
			}
		} else {
			if src.Kind() != reflect.Interface && src.Kind() != reflect.Ptr {
				dst.Set(src)
			} else {
				if err := bind(src.Elem(), dst); err != nil {
					return err
				}
			}
		}
	case reflect.Array, reflect.Slice:
		var isrc reflect.Value
		if src.Kind() != reflect.Slice {
			isrc = src.Elem()
		} else {
			isrc = src
		}
		if isrc.IsValid() {
			slice := reflect.MakeSlice(dst.Type(), isrc.Len(), isrc.Len())
			for i := 0; i < isrc.Len(); i++ {
				var isrcIndex reflect.Value
				if isrc.Index(i).Kind() != reflect.Interface {
					isrcIndex = isrc.Index(i)
				} else {
					isrcIndex = reflect.Indirect(isrc.Index(i))
				}
				if slice.Index(i).Kind() == reflect.Ptr && slice.Index(i).IsNil() {
					pv := reflect.New(slice.Index(i).Type())
					if err := bind(isrcIndex, pv); err != nil {
						return err
					}
					slice.Index(i).Set(reflect.Indirect(pv))
				} else {
					if err := bind(isrcIndex, slice.Index(i)); err != nil {
						return err
					}
				}
			}
			dst.Set(slice)
		}
	case reflect.String:
		dst.Set(reflect.ValueOf(src.Interface()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if src.Kind() != reflect.Interface {
			dst.Set(src.Convert(dst.Type()))
		} else {
			dst.Set(src.Elem().Convert(dst.Type()))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if src.Kind() != reflect.Interface {
			dst.Set(src.Convert(dst.Type()))
		} else {
			dst.Set(src.Elem().Convert(dst.Type()))
		}
	case reflect.Float32, reflect.Float64:
		if src.Kind() != reflect.Interface {
			dst.Set(src.Convert(dst.Type()))
		} else {
			dst.Set(src.Elem().Convert(dst.Type()))
		}
	case reflect.Bool:
		dst.SetBool(src.Elem().Bool())
	case reflect.Ptr:
		if !dst.IsNil() {
			if dst.Elem().Kind() != reflect.Interface && dst.Elem().Kind() != reflect.Ptr {
				if src.Kind() != reflect.Interface && src.Kind() != reflect.Ptr {
					if src.Kind() == reflect.Array || src.Kind() == reflect.Slice {
						if err := bind(src, dst.Elem()); err != nil {
							return err
						}
					} else {
						dst.Elem().Set(src.Convert(dst.Type().Elem()))
					}
				} else {
					if err := bind(src.Elem(), dst); err != nil {
						return err
					}
				}
			} else {
				if err := bind(src, dst.Elem()); err != nil {
					return err
				}
			}
		} else {
			pv := reflect.New(dst.Type().Elem())
			if dst.Elem().Kind() != reflect.Interface && dst.Elem().Kind() != reflect.Ptr {
				if src.Kind() != reflect.Interface && src.Kind() != reflect.Ptr {
					pv.Elem().Set(src.Convert(dst.Type().Elem()))
				} else {
					if err := bind(src.Elem(), pv); err != nil {
						return err
					}
				}
			} else {
				if err := bind(src, dst.Elem()); err != nil {
					return err
				}
			}
			dst.Set(pv)
		}
	default:
		if src.Kind() != reflect.Interface && src.Kind() != reflect.Ptr {
			dst.Set(src)
		} else {
			if err := bind(src.Elem(), dst); err != nil {
				return err
			}
		}
	}

	return nil
}

func decode(fields interface{}) (interface{}, error) {
	mapFields, ok := fields.(map[string]interface{})
	if !ok {
		switch reflect.ValueOf(fields).Kind() {
		case reflect.Array, reflect.Slice:
			return decodeArray(fields)
		case reflect.Interface, reflect.Ptr:
			return decode(reflect.Indirect(reflect.ValueOf(fields)).Interface())
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
			return decodePointer(fields)
		case "Object":
			return decodeObject(fields)
		case "Date":
			iso, ok := mapFields["iso"].(string)
			if !ok {
				return nil, fmt.Errorf("unexpected error when parse Date: iso expected string but %v", reflect.TypeOf(mapFields["iso"]))
			}
			return decodeDate(iso)
		case "Bytes":
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
	decodedFields["createdAt"] = decodedCreatedAt

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
	decodedFields["updatedAt"] = decodedUpdatedAt

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
			switch iter.Key().String() {
			case "ACL":

			default:
				r, err := decode(iter.Value().Interface())
				if err != nil {
					return nil, err
				}
				decodedMap[iter.Key().String()] = r
			}
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

func decodeDate(dateStr string) (*time.Time, error) {
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return &time.Time{}, err
	}
	return &date, nil
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

func decodeACL(fields map[string]interface{}) (*ACL, error) {
	return nil, nil
}

func parseTag(tag string) (name string, option string) {
	parts := strings.Split(tag, ",")

	if len(parts) > 1 {
		return parts[0], parts[1]
	}

	return parts[0], ""
}

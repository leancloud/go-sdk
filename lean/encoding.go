package lean

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const StructFieldTag = "json"

func decodeObjects(className string, source []objectResponse, objects interface{}) error {
	objectType := reflect.TypeOf(objects).Elem().Elem()
	objectsInterface := reflect.ValueOf(objects)
	resultSlice := reflect.Indirect(reflect.ValueOf(objects))

	for _, objectResp := range source {
		object := reflect.New(objectType).Interface().(Object)
		err := decodeObject(objectResp, object)

		if err != nil {
			return err
		}

		resultSlice = reflect.Append(resultSlice, reflect.ValueOf(object).Elem())
	}

	objectsInterface.Elem().Set(resultSlice)

	return nil
}

func decodeObject(source objectResponse, object Object) error {
	objectStruct := reflect.Indirect(reflect.ValueOf(object))
	objectType := objectStruct.Type()

	// fmt.Println(source)

	for i := 0; i < objectType.NumField(); i++ {
		field := objectType.Field(i)

		jsonFieldName, tagSet := field.Tag.Lookup(StructFieldTag)

		if !tagSet {
			jsonFieldName = field.Name
		}

		switch v := reflect.New(field.Type).Interface().(type) {
		case *ObjectMeta:
			meta := objectStruct.Field(i).Addr().Interface().(*ObjectMeta)

			var err error

			meta.ObjectID = source["objectId"].(string)
			meta.CreatedAt, err = time.Parse(time.RFC3339Nano, source["createdAt"].(string))

			if err != nil {
				panic(err)
			}

			meta.UpdatedAt, err = time.Parse(time.RFC3339Nano, source["updatedAt"].(string))

			if err != nil {
				panic(err)
			}

		case *String:
			sourceFieldValue := reflect.ValueOf(source[jsonFieldName])

			if sourceFieldValue.Kind() == reflect.String {
				objectStruct.Field(i).Set(reflect.ValueOf(NewString(source[jsonFieldName].(string))))
			}

		case *Integer:
			sourceFieldValue := reflect.ValueOf(source[jsonFieldName])

			if sourceFieldValue.Kind() == reflect.Int {
				objectStruct.Field(i).Set(reflect.ValueOf(NewInteger(int(source[jsonFieldName].(int)))))
			}

		case *Float:
			sourceFieldValue := reflect.ValueOf(source[jsonFieldName])

			if sourceFieldValue.Kind() == reflect.Float32 {
				objectStruct.Field(i).Set(reflect.ValueOf(NewFloat(float64(source[jsonFieldName].(float64)))))
			}

		case *Boolean:
			sourceFieldValue := reflect.ValueOf(source[jsonFieldName])

			if sourceFieldValue.Kind() == reflect.Bool {
				objectStruct.Field(i).Set(reflect.ValueOf(NewBoolean(source[jsonFieldName].(bool))))
			}

		case *Date:
			sourceFieldValue := reflect.ValueOf(source[jsonFieldName])

			if sourceFieldValue.Kind() == reflect.String {
				dateTime, err := time.Parse(time.RFC3339Nano, source[jsonFieldName].(string))

				if err != nil {
					panic(err)
				}

				objectStruct.Field(i).Set(reflect.ValueOf(NewDate(dateTime)))
			}

		default:
			fmt.Println("Unknown type", v)
		}

		// fmt.Println(field)
	}

	return nil
}

func mergeToObject(source []byte, object Object) error {
	fmt.Println("source", string(source))
	return json.Unmarshal(source, object)

	// objectStruct := reflect.Indirect(reflect.ValueOf(object))
	//
	// var err error
	//
	// meta := object.getObjectMeta()
	//
	// mapping, _ := getFieldMapping(object)
	//
	// for jsonField, value := range source {
	// 	switch jsonField {
	// 	case "objectId":
	// 		meta.ObjectID = value.(string)
	// 	case "createdAt":
	// 		meta.CreatedAt, err = time.Parse(time.RFC3339Nano, value.(string))
	//
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	default:
	// 		fieldName, ok := mapping[jsonField]
	//
	// 		if ok {
	// 			objectStruct.FieldByName(fieldName).Set
	// 		}
	//
	// 	}
	// }
	//
	// return nil
}

func encodeObject(object Object, dest map[string]interface{}) error {
	objectStruct := reflect.Indirect(reflect.ValueOf(object))
	objectType := objectStruct.Type()

	_, mapping := getFieldMapping(object)

	for i := 0; i < objectType.NumField(); i++ {
		field := objectType.Field(i)

		jsonFieldName, ok := mapping[field.Name]

		if ok {
			dest[jsonFieldName] = objectStruct.Field(i).Interface()
		}
	}

	return nil
}

func parseTag(tag string) (name string, option string) {
	parts := strings.Split(tag, ",")

	if len(parts) > 1 {
		return parts[0], parts[1]
	}

	return parts[0], ""
}

func getFieldMapping(object Object) (jsonToField map[string]string, fieldToJSON map[string]string) {
	jsonToField = make(map[string]string)
	fieldToJSON = make(map[string]string)

	objectStruct := reflect.Indirect(reflect.ValueOf(object))
	objectType := objectStruct.Type()

	for i := 0; i < objectType.NumField(); i++ {
		field := objectType.Field(i)

		if field.Anonymous {
			continue
		}

		jsonFieldName := field.Name

		fieldTag, tagSet := field.Tag.Lookup(StructFieldTag)

		var omitEmpty bool

		if tagSet {
			fieldName, option := parseTag(fieldTag)

			if fieldName == "" || fieldName == "-" {
				continue
			}

			omitEmpty = option == "omitempty"

			jsonFieldName = fieldName
		}

		jsonToField[jsonFieldName] = field.Name

		if !omitEmpty || !objectStruct.Field(i).Addr().MethodByName("IsNull").Call([]reflect.Value{})[0].Bool() {
			fieldToJSON[field.Name] = jsonFieldName
		}
	}

	return jsonToField, fieldToJSON
}

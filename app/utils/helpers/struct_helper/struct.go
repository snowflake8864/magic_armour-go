package struct_helper

import (
	"reflect"
	"strings"
)

// Copy 复制结构体
// dst 目标结构体，src 源结构体
// 必须传入指针，且不能为nil
// 它会把src与dst的相同字段名的值，复制到dst中
func Copy(dst, src interface{}) {
	dstValue := reflect.ValueOf(dst).Elem()
	srcValue := reflect.ValueOf(src).Elem()

	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		srcName := srcValue.Type().Field(i).Name
		dstFieldByName := dstValue.FieldByName(srcName)

		if dstFieldByName.IsValid() {
			switch dstFieldByName.Kind() {
			case reflect.Ptr:
				switch srcField.Kind() {
				case reflect.Ptr:
					if srcField.IsNil() {
						dstFieldByName.Set(reflect.New(dstFieldByName.Type().Elem()))
					} else {
						dstFieldByName.Set(srcField)
					}
				default:
					dstFieldByName.Set(srcField.Addr())
				}
			default:
				switch srcField.Kind() {
				case reflect.Ptr:
					if srcField.IsNil() {
						dstFieldByName.Set(reflect.Zero(dstFieldByName.Type()))
					} else {
						dstFieldByName.Set(srcField.Elem())
					}
				default:
					dstFieldByName.Set(srcField)
				}
			}
		}
	}
}

// SimpleStructToMapWithJsonTags  converts a simple struct to a map[string]interface{} using JSON tags as keys.
// 忽略字段是: Struct\Slice\Array
func SimpleStructToMapWithJsonTags(obj interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil // or handle non-struct types differently if needed
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		jsonTag := field.Tag.Get("json")

		// If no JSON tag is specified, use the field name as the key
		if jsonTag == "" {
			jsonTag = field.Name
		} else if jsonTag == "-" {
			// Skip fields with "-" JSON tag
			continue
		} else if strings.Contains(jsonTag, ",") {
			// Split the tag into name and options, and use the name part
			jsonTag = strings.SplitN(jsonTag, ",", 2)[0]
		}

		// Handle nested structs and slices recursively
		switch fieldVal.Kind() {
		case reflect.Struct:
			break
		case reflect.Slice, reflect.Array:
			break
		default:
			data[jsonTag] = fieldVal.Interface()
		}
	}

	return data
}

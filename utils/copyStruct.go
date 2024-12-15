package utils

import (
	"reflect"
)

// CopyStruct copies all fields from src to dst
func CopyStruct(src, dst interface{}) {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		field := srcVal.Type().Field(i)
		dstField := dstVal.FieldByName(field.Name)
		if dstField.IsValid() && dstField.CanSet() {
			dstField.Set(srcVal.Field(i))
		}
	}
}

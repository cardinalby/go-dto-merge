package dtomerge

import (
	"fmt"
	"reflect"
)

func mergeStructs(src reflect.Value, patch reflect.Value, options Options) (reflect.Value, error) {
	result := reflect.New(src.Type()).Elem()
	result.Set(src)
	structType := src.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		srcFieldValue := src.FieldByIndex(field.Index)
		resultFieldValue := result.FieldByIndex(field.Index)
		if field.IsExported() {
			mergedFieldValue, err := mergeAny(srcFieldValue, patch.FieldByIndex(field.Index), options)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("'%s' field: %w", field.Name, err)
			}
			resultFieldValue.Set(mergedFieldValue)
		}
	}
	return result, nil
}

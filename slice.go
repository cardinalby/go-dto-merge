package dtomerge

import (
	"fmt"
	"reflect"
)

func mergeComparableSlicesUnique(src, patch reflect.Value) (reflect.Value, error) {
	result := reflect.MakeSlice(src.Type(), 0, src.Len())
	srcElemsMap := reflect.MakeMap(reflect.MapOf(src.Type().Elem(), reflect.TypeOf(true)))
	// index all src elements
	trueValue := reflect.ValueOf(true)
	for i := 0; i < src.Len(); i++ {
		elem := src.Index(i)
		srcElemsMap.SetMapIndex(elem, trueValue)
		result = reflect.Append(result, elem)
	}

	for i := 0; i < patch.Len(); i++ {
		if foundSrcElam := srcElemsMap.MapIndex(patch.Index(i)); !foundSrcElam.IsValid() {
			result = reflect.Append(result, patch.Index(i))
		}
	}
	return result, nil
}

func mergeSlicesByIndex(src, patch reflect.Value, mCtx mergeContext) (reflect.Value, error) {
	result := reflect.MakeSlice(src.Type(), 0, src.Len())
	for i := 0; i < src.Len(); i++ {
		result = reflect.Append(result, src.Index(i))
	}
	for i := 0; i < patch.Len(); i++ {
		if i < result.Len() {
			mergedElement, err := mergeAny(result.Index(i), patch.Index(i), mCtx)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("merging patch slice[%d]: %w", i, err)
			}
			result.Index(i).Set(mergedElement)
		} else {
			result = reflect.Append(result, patch.Index(i))
		}
	}
	return result, nil
}
